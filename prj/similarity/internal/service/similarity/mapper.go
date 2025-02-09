package similarity

import "github.com/zuzuka28/simreport/prj/similarity/internal/model"

func extractSourceIDs(in []*model.SimilarityMatch) []string {
	ids := make([]string, 0, len(in))
	for _, v := range in {
		ids = append(ids, v.ID)
	}

	return ids
}
