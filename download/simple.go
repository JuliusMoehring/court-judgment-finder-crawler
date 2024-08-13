package download

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/JuliusMoehring/court-judgment-finder-crawler/logger"
)

type SimpleDownloader struct {
	logger logger.Logger
}

func NewSimpleDownloader(logger logger.Logger) Downloader {
	return &SimpleDownloader{
		logger: logger,
	}
}

func (d *SimpleDownloader) Download(ctx context.Context, url string) ([]byte, error) {
	client := &http.Client{}

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", response.Status)
	}

	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
