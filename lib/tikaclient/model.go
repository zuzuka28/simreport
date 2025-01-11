package tikaclient

import "fmt"

//nolint:gochecknoglobals
var OCRLanguage = []string{"rus", "eng"}

type Result struct {
	Content []byte
	Sha256  string
	Name    string
}

type ClientError struct {
	StatusCode int
}

func (e ClientError) Error() string {
	return fmt.Sprintf("response code %d", e.StatusCode)
}
