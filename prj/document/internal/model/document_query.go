package model

type DocumentQueryInclude string

const (
	DocumentQueryIncludeSource = "source"
	DocumentQueryIncludeText   = "text"
	DocumentQueryIncludeImages = "images"
)

type DocumentQuery struct {
	ID          string
	WithContent bool
	Include     []DocumentQueryInclude
}
