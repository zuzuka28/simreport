package model

type AnalyzedImage struct {
	ID        string
	Vector    Vector
	HashImage HashImage
}

type AnalyzedDocument struct {
	ID         string
	Text       string
	TextVector Vector
	Images     []AnalyzedImage
}

func (d *AnalyzedDocument) ImagesIDs() []string {
	res := make([]string, 0, len(d.Images))
	for _, v := range d.Images {
		res = append(res, v.ID)
	}

	return res
}
