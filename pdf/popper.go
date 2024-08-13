package pdf

import (
	"context"
	"os/exec"
)

type PopperPDFReader struct {
}

func NewPopperPDFReader() Reader {
	_, err := exec.LookPath("pdftotext")
	if err != nil {
		panic("pdftotext not found in PATH")
	}

	return &PopperPDFReader{}
}

func (p *PopperPDFReader) Read(ctx context.Context, path string) ([]byte, error) {
	bytes, err := exec.CommandContext(ctx, "pdftotext", path, "-").Output()
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
