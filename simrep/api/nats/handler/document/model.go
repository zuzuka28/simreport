package document

type Document struct {
	ID   string `json:"id"`
	Text []byte `json:"text"`
}
