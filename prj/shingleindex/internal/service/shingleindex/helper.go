package shingleindex

import (
	"regexp"
	"strings"
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
