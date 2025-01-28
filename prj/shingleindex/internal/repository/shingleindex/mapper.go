package shingleindex

import "github.com/zuzuka28/simreport/prj/shingleindex/internal/model"

func mapCandidatesToMatches(in []string) []*model.MinhashSimilarMatch {
	res := make([]*model.MinhashSimilarMatch, 0, len(in))

	for _, v := range in {
		res = append(res, &model.MinhashSimilarMatch{
			DocumentID: v,
		})
	}

	return res
}
