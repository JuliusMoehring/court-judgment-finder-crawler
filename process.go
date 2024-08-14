package main

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/JuliusMoehring/court-judgment-finder-crawler/bgh"
	"github.com/JuliusMoehring/court-judgment-finder-crawler/download"
	"github.com/JuliusMoehring/court-judgment-finder-crawler/embedder"
	filestorage "github.com/JuliusMoehring/court-judgment-finder-crawler/file-storage"
	"github.com/JuliusMoehring/court-judgment-finder-crawler/logger"
	"github.com/JuliusMoehring/court-judgment-finder-crawler/pdf"
	vectorstore "github.com/JuliusMoehring/court-judgment-finder-crawler/vector-store"
)

type Processor struct {
	logger      logger.Logger
	downloader  download.Downloader
	fileStorage filestorage.FileStorage
	pdfReader   pdf.Reader
	embedder    embedder.Embedder
	vectorStore vectorstore.VectorStore
}

func NewProcessor(logger logger.Logger, downloader download.Downloader, fileStorage filestorage.FileStorage, pdfReader pdf.Reader, embedder embedder.Embedder, vectorStore vectorstore.VectorStore) *Processor {
	return &Processor{
		logger:      logger,
		downloader:  downloader,
		fileStorage: fileStorage,
		pdfReader:   pdfReader,
		embedder:    embedder,
		vectorStore: vectorStore,
	}
}

func (p *Processor) shouldProcessLink(ctx context.Context, path string) (bool, error) {
	pdfUploaded, err := p.fileStorage.Exists(ctx, path)
	if err != nil {
		p.logger.Errorf("processor", "failed checking if document is uploaded: %s", err)
		return false, err
	}

	documentID, err := p.vectorStore.GetDocumentIDByPath(ctx, path)
	// If the error is not that the document is not found, we return the error
	if err != nil && !errors.Is(err, vectorstore.ErrDocumentNotFound) {
		p.logger.Errorf("processor", "failed checking if document exists in vector store: %s", err)
		return false, err
	}

	// If the PDF is already uploaded and the document is already in the vector store, we can skip this PDF
	if pdfUploaded && documentID != "" {
		p.logger.Debugf("processor", "skipping already uploaded document: '%s', id: '%s'", path, documentID)
		return false, nil
	}

	return true, nil
}

func (p *Processor) pdfToText(ctx context.Context, data []byte) (string, error) {
	if err := os.MkdirAll("temp", os.ModePerm); err != nil {
		p.logger.Errorf("processor", "failed creating temp directory: %s", err)
		return "", err
	}

	file, err := os.CreateTemp("temp", "temp-*.pdf")
	if err != nil {
		p.logger.Errorf("processor", "failed creating temp file: %s", err)
		return "", err
	}

	defer func() {
		if err := file.Close(); err != nil {
			p.logger.Errorf("processor", "failed closing temp file: %s", err)
		}

		if err := os.Remove(file.Name()); err != nil {
			p.logger.Errorf("processor", "failed removing temp file: %s", err)
		}
	}()

	if _, err = file.Write(data); err != nil {
		p.logger.Errorf("processor", "failed writing data to temp file: %s", err)
		return "", err
	}

	bytes, err := p.pdfReader.Read(ctx, file.Name())
	if err != nil {
		p.logger.Errorf("processor", "failed reading temp file: %s", err)
		return "", err
	}

	return string(bytes), nil
}

func (p *Processor) processLink(ctx context.Context, link string) error {
	path, err := bgh.PathFromURL(link)
	if err != nil {
		p.logger.Errorf("failed to create path from url '%s': %s", link, err)
		return nil
	}

	shouldProcess, err := p.shouldProcessLink(ctx, path)
	if err != nil {
		return err
	}

	if !shouldProcess {
		return nil
	}

	data, err := p.downloader.Download(ctx, link)
	if err != nil {
		p.logger.Errorf("processor", "failed downloading document: %s", err)
		return err
	}

	if err := p.fileStorage.Save(ctx, data, path); err != nil {
		p.logger.Errorf("processor", "failed saving document to file storage: %s", err)
		return err
	}

	text, err := p.pdfToText(ctx, data)
	if err != nil {
		p.logger.Errorf("processor", "failed converting pdf to text: %s", err)
		return err
	}

	var judgementPages []vectorstore.CreateDocumentParamsPage

	pages := strings.Split(text, "\f")

	for i, page := range pages {
		if len(page) == 0 {
			continue
		}

		embedding, err := p.embedder.Embed(ctx, page)
		if err != nil {
			p.logger.Errorf("processor", "failed to create embddings for page %d of link %s: %s", i+1, link, err)
			return err
		}

		judgementPages = append(judgementPages, vectorstore.CreateDocumentParamsPage{
			Text:      page,
			Embedding: embedding,
		})
	}

	return p.vectorStore.CreateDocument(ctx, vectorstore.CreateDocumentParams{
		FilePath: path,
		Pages:    judgementPages,
	})
}

func (p *Processor) Process(ctx context.Context, downloadLinks <-chan string) error {
	for link := range downloadLinks {
		p.logger.Debugf("processor", "processing link: '%s'", link)

		if err := p.processLink(ctx, link); err != nil {
			p.logger.Errorf("processor", "failed processing link: '%s'", err)
			return err
		}

		p.logger.Debugf("processor", "processed link: '%s'. %d more links to process.", link, len(downloadLinks))
	}

	return nil
}
