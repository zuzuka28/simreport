package document

type parsedDocument struct {
	ID          string   `json:"id"`
	Sha256      string   `json:"sha256"`
	ImageIDs    []string `json:"imageIDs"`
	TextContent string   `json:"textContent"`
}
