package docxparser

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"regexp"
	"simrep/internal/model"
	"strings"
	"time"
)

var errEmptyFile = errors.New("file content is empty")

var (
	mainRegex   = regexp.MustCompile(`word/document.xml`)
	headerRegex = regexp.MustCompile(`word/header[0-9]*\.xml`)
	footerRegex = regexp.MustCompile(`word/footer[0-9]*\.xml`)
	imageRegex  = regexp.MustCompile(`word/media/.*\.(jpg|jpeg|png|bmp)`)
)

func Parse(item model.File) (model.Document, error) {
	if len(item.Content) == 0 {
		return model.Document{}, errEmptyFile
	}

	text, images, err := processDocx(item.Content)
	if err != nil {
		return model.Document{}, err
	}

	imageIDs := make([]string, len(images))
	for i, img := range images {
		imageIDs[i] = img.Sha256
	}

	return model.Document{
		ID:          item.Sha256,
		Name:        item.Name,
		ImageIDs:    imageIDs,
		TextID:      text,
		LastUpdated: time.Now(),
		Source:      item,
		Images:      images,
	}, nil
}

func processDocx(docx []byte) (string, []model.File, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(docx), int64(len(docx)))
	if err != nil {
		return "", nil, fmt.Errorf("new zip reader: %w", err)
	}

	text, err := extractText(zipReader)
	if err != nil {
		return "", nil, err
	}

	images, err := extractImages(zipReader)
	if err != nil {
		return "", nil, err
	}

	return text, images, nil
}

func extractText(zipReader *zip.Reader) (string, error) {
	var builder strings.Builder

	mainFiles := filesByRegex(zipReader, mainRegex)

	for _, f := range mainFiles {
		data, err := readZipFile(f)
		if err != nil {
			return "", err
		}

		builder.WriteString(xml2text(data))
		builder.WriteString("\n\n")
	}

	return builder.String(), nil
}

func extractImages(zipReader *zip.Reader) ([]model.File, error) {
	imgFiles := filesByRegex(zipReader, imageRegex)
	imgs := make([]model.File, 0, len(imgFiles))

	for _, file := range imgFiles {
		content, err := readZipFile(file)
		if err != nil {
			return nil, fmt.Errorf("read zip file: %w", err)
		}

		imgs = append(imgs, model.File{
			Name:        "",
			Content:     content,
			Sha256:      sha256String(content),
			LastUpdated: time.Time{},
		})
	}

	return imgs, nil
}

func filesByRegex(
	zipReader *zip.Reader,
	re *regexp.Regexp,
) []*zip.File {
	var files []*zip.File //nolint:prealloc

	for _, f := range zipReader.File {
		if !re.MatchString(f.Name) {
			continue
		}

		files = append(files, f)
	}

	return files
}

func xml2text(xmlStr []byte) string {
	decoder := xml.NewDecoder(bytes.NewReader(xmlStr))

	var result strings.Builder

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		if startElem, ok := token.(xml.StartElement); ok {
			processStartElement(decoder, startElem, &result)
		}
	}

	return result.String()
}

func processStartElement(
	decoder *xml.Decoder,
	elem xml.StartElement,
	result *strings.Builder,
) {
	switch elem.Name.Local {
	case "t":
		var text string

		if err := decoder.DecodeElement(&text, &elem); err == nil {
			result.WriteString(text)
		}

	case "tab":
		result.WriteString("\t")

	case "br":
		result.WriteString("\n")

	case "p":
		result.WriteString("\n\n")
	}
}

func readZipFile(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, fmt.Errorf("open zip file: %w", err)
	}

	defer rc.Close()

	res, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("read content: %w", err)
	}

	return res, nil
}

func sha256String(in []byte) string {
	hash := sha256.New()
	_, _ = hash.Write(in)

	return hex.EncodeToString(hash.Sum(nil))
}
