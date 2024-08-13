package bgh

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JuliusMoehring/court-judgment-finder-crawler/logger"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

const (
	I_ZIVIL_SENATE = "I. Zivilsenat"
	X_ZIVIL_SENATE = "X. Zivilsenat"
	BASE_URL       = "https://juris.bundesgerichtshof.de/cgi-bin/rechtsprechung"
)

type links struct {
	mu sync.Mutex

	links map[string]bool
}

func newLinks() *links {
	return &links{
		links: map[string]bool{},
	}
}

func (l *links) addLink(link string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.links[link] = true
}

func (l *links) getLinks() []string {
	linkSlice := make([]string, 0, len(l.links))

	for link := range l.links {
		linkSlice = append(linkSlice, link)
	}

	return linkSlice
}

type Crawler struct {
	logger logger.Logger
}

func NewCrawler(logger logger.Logger) *Crawler {
	return &Crawler{
		logger: logger,
	}
}

func (c *Crawler) getOverviewURL(year int, page int) (string, error) {
	baseURL, err := url.Parse("https://juris.bundesgerichtshof.de/cgi-bin/rechtsprechung/list.py?Gericht=bgh&Art=en")
	if err != nil {
		return "", err
	}

	query := baseURL.Query()

	query.Add("Datum", strconv.Itoa(year))

	if page > 1 {
		query.Add("Seite", strconv.Itoa(page))
	}

	baseURL.RawQuery = query.Encode()

	return baseURL.String(), nil
}

// Gets all years where court judgments are available
func (c *Crawler) getAvailableYears(collector *colly.Collector, years *[]int) {
	collector.OnHTML("#kaljahr", func(e *colly.HTMLElement) {
		start := time.Now()

		e.DOM.Children().Each(func(index int, item *goquery.Selection) {
			year, err := strconv.Atoi(item.Text())
			if err != nil {
				fmt.Println("Error parsing year: ", err)
				return
			}

			*years = append(*years, year)
		})

		c.logger.Debugf("crawler", "Got %d available years, took %s", len(*years), time.Since(start))
	})
}

// Gets the maximum pages that are available
func (c *Crawler) getAvailablePages(collector *colly.Collector, availablePages *int) {
	collector.OnHTML("body > table.rechts > tbody > tr:nth-child(1) > td:nth-child(4) > table > tbody > tr:nth-child(2) > td > form > table > thead > tr > td.ETitelKopf > table > tbody > tr > td:nth-child(5) > a:nth-child(2)", func(e *colly.HTMLElement) {
		href := e.Attr("href")

		hrefSegments := strings.Split(href, "?")

		if len(hrefSegments) != 2 {
			fmt.Println("Invalid href found")
			return
		}

		queryParams, err := url.ParseQuery(hrefSegments[1])
		if err != nil {
			fmt.Println("Error parsing query params: ", err)
			return
		}

		*availablePages, err = strconv.Atoi(queryParams.Get("Seite"))
		if err != nil {
			fmt.Println("Error parsing available pages: ", err)
			return
		}
	})
}

func (c *Crawler) crawlYear(collector *colly.Collector, year int, links *links) error {
	collector = collector.Clone()

	var availablePages int

	c.getAvailablePages(collector, &availablePages)

	initialURL, err := c.getOverviewURL(year, 1)
	if err != nil {
		return err
	}

	if err = collector.Visit(initialURL); err != nil {
		return err
	}

	for i := 1; i <= availablePages; i++ {
		url, err := c.getOverviewURL(year, i)
		if err != nil {
			return err
		}

		c.crawlTable(collector, links)

		if err = collector.Visit(url); err != nil {
			return err
		}
	}

	return nil
}

// Crawls the table with the court judgments
func (c *Crawler) crawlTable(collector *colly.Collector, links *links) {
	// Get the table with the court judgments
	collector.OnHTML("body > table.rechts > tbody > tr:nth-child(1) > td:nth-child(4) > table > tbody > tr:nth-child(2) > td > form > table > tbody", func(e *colly.HTMLElement) {
		for _, child := range e.DOM.Children().Nodes {
			row := goquery.NewDocumentFromNode(child)

			// Get the text of the 1st column
			senate := strings.TrimSpace(row.Children().First().Text())

			// Only add links for the predefined senates
			if senate != I_ZIVIL_SENATE && senate != X_ZIVIL_SENATE {
				continue
			}

			// Get the 2nd link in the 3rd column
			link := row.Children().Find("a[type=\"application/pdf\"]")

			href, exists := link.Attr("href")
			if !exists {
				fmt.Println("No href found for link")
				continue
			}

			links.addLink(BASE_URL + "/" + href)
		}
	})
}

// Crawls the Bundesgerichtshof website for court judgments of predefined senates
func (c *Crawler) Crawl(ctx context.Context) ([]string, error) {
	collector := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.AllowedDomains("juris.bundesgerichtshof.de"),
		colly.CacheDir("./bgh/cache/"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	years := []int{}
	links := newLinks()

	c.getAvailableYears(collector, &years)

	initialURL, err := c.getOverviewURL(2024, 1)
	if err != nil {
		return nil, err
	}

	if err = collector.Visit(initialURL); err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	for _, year := range years {
		wg.Add(1)
		go func(collector *colly.Collector, year int) {
			defer wg.Done()

			err := c.crawlYear(collector, year, links)
			if err != nil {
				fmt.Println("Error crawling year: ", err)
			}
		}(collector, year)
	}

	wg.Wait()

	c.logger.Debugf("crawler", "Finished crawling Bundesgerichtshof website, found %d unique pdf links", len(links.getLinks()))

	return links.getLinks(), nil
}
