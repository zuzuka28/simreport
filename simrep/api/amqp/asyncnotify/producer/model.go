package producer

type notification struct {
	DocumentID string `json:"documentID"`
	Action     string `json:"action"`
	UserData   any    `json:"userData"`
}
