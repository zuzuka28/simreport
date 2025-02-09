package similarity

import (
	"regexp"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/zuzuka28/simreport/prj/similarity/internal/model"
)

var (
	reNormalize1 = regexp.MustCompile(`[^a-zA-Zа-яА-ЯёЁ0-9_ ]`)
	reNormalize2 = regexp.MustCompile(`\s+`)
)

func normalize(text string) string {
	text = reNormalize1.ReplaceAllString(text, "")
	text = reNormalize2.ReplaceAllString(text, " ")
	text = strings.ToLower(text)

	return text
}

func shingle(text string, shingleSize int) map[string]struct{} {
	words := strings.Fields(text)

	var shingles []string

	for i := 0; i <= len(words)-shingleSize; i++ {
		shingle := strings.Join(words[i:i+shingleSize], " ")
		shingles = append(shingles, shingle)
	}

	shingleSet := make(map[string]struct{})
	for _, shingle := range shingles {
		shingleSet[shingle] = struct{}{}
	}

	return shingleSet
}

func jaccardSimilarity(set1, set2 map[string]struct{}) float64 {
	intersection := 0

	for item := range set1 {
		if _, exists := set2[item]; exists {
			intersection++
		}
	}

	union := len(set1) + len(set2) - intersection

	if union == 0 {
		return 0
	}

	return float64(intersection) / float64(union)
}

func normalizeMatches(
	items map[string]*match,
) map[string]*match {
	for k, v := range items {
		v.text = normalize(v.text)

		items[k] = v
	}

	return items
}

func shingleMatches(
	items map[string]*match,
) map[string]*match {
	for k, v := range items {
		v.shingles = shingle(v.text, shingleSize)
		items[k] = v
	}

	return items
}

func rerank(
	q string,
	items map[string]*match,
) map[string]*match {
	sourcesnormalized := normalizeMatches(items)
	sourcesshingled := shingleMatches(sourcesnormalized)

	text := string(q)
	text = normalize(text)
	textshingles := shingle(text, shingleSize)

	for k, v := range sourcesshingled {
		v.Rate = jaccardSimilarity(textshingles, v.shingles)
		sourcesshingled[k] = v
	}

	return items
}

func highlight(
	q string,
	items map[string]*match,
) map[string]*match {
	dmp := diffmatchpatch.New()

	for _, v := range items {
		var highlights []string

		diffs := dmp.DiffMain(q, v.text, false)

		for _, v := range diffs {
			if v.Type == diffmatchpatch.DiffEqual {
				highlights = append(highlights, v.Text)
			}
		}

		v.Highlights = highlights
	}

	return items
}

func deduplicateByID(
	items []*model.SimilarityMatch,
) []*model.SimilarityMatch {
	uniq := make(map[string]struct{})

	var result []*model.SimilarityMatch

	for _, v := range items {

		_, ok := uniq[v.ID]
		if ok {
			continue
		}

		result = append(result, v)
		uniq[v.ID] = struct{}{}
	}

	return result
}
