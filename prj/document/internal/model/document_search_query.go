package model

type DocumentSearchQuery struct {
	GroupID  []string
	Name     string
	ParentID string
	SourceID []string
	Version  string
}
