package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/JuliusMoehring/court-judgment-finder-crawler/bgh"
	"github.com/JuliusMoehring/court-judgment-finder-crawler/download"
	"github.com/JuliusMoehring/court-judgment-finder-crawler/embedder"
	filestorage "github.com/JuliusMoehring/court-judgment-finder-crawler/file-storage"
	"github.com/JuliusMoehring/court-judgment-finder-crawler/logger"
	"github.com/JuliusMoehring/court-judgment-finder-crawler/pdf"
	vectorstore "github.com/JuliusMoehring/court-judgment-finder-crawler/vector-store"
	"github.com/joho/godotenv"
)

const WORKERS = 10

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	ctx := context.Background()

	// Initialize services
	logger := logger.NewStdOutLogger()
	downloader := download.NewSimpleDownloader(logger)
	fileStorage := filestorage.NewS3FileStorage(ctx, logger, "court-judgement-finder")
	pdfReader := pdf.NewPopperPDFReader()
	embedder := embedder.NewOpenAIEmbedder()
	vectorStore := vectorstore.NewSurrealDBVectorStore()
	defer vectorStore.Close()

	// Initialize crawler
	crawler := bgh.NewCrawler(logger)

	links, err := crawler.Crawl(ctx)
	if err != nil {
		panic(fmt.Sprintf("could not crawl BGH: %s", err))
	}

	processor := NewProcessor(logger, downloader, fileStorage, pdfReader, embedder, vectorStore)

	downloadLinks := make(chan string, len(links))

	var wg sync.WaitGroup

	for i := 0; i < WORKERS; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := processor.Process(ctx, downloadLinks); err != nil {
				logger.Errorf("error processing PDF links: %s", err)
			}
		}()
	}

	for _, link := range links {
		downloadLinks <- link
	}
	close(downloadLinks)

	wg.Wait()

}
