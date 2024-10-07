package main

import (
	"context"
	"fmt"
	"io"
	"os"
	docprepsrvrepo "simrep/internal/repository/docprepsrv"
	"simrep/internal/service/docprepsrv"
)

func readFile(fname string) ([]byte, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	raw, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	return raw, nil
}

func main() {
	ctx := context.Background()

	r, err := docprepsrvrepo.NewRepository(docprepsrvrepo.Opts{
		APIURL:     "http://localhost:8000",
		HTTPClient: nil,
	})
	if err != nil {
		panic(err)
	}

	s := docprepsrv.NewService(r)

	doc, err := readFile("./data/example.docx")
	if err != nil {
		panic(err)
	}

	res, err := s.PreprocessRawDocument(ctx, doc)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", res)
}
