package shingleindex

import "shingleindex/internal/model"

func mapCandidatesToMatches(in []string) []*model.MinhashSimilarMatch {
	res := make([]*model.MinhashSimilarMatch, 0, len(in))

	for _, v := range in {
		res = append(res, &model.MinhashSimilarMatch{
			DocumentID: v,
		})
	}

	return res
}
