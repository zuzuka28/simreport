package model

import (
	"strconv"
)

type Document struct {
	ParentID string
	Name     string

	Version int
	GroupID []string

	Source File
}

func (d *Document) ID() string {
	return d.ParentID + "_" + strconv.Itoa(d.Version)
}
