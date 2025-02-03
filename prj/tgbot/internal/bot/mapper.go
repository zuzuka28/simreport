package bot

import "github.com/zuzuka28/simreport/prj/tgbot/internal/model"

func mapSimilarityMatchesToResponse(
	in []*model.SimilarityMatch,
) []*similarityMatch {
	res := make([]*similarityMatch, 0, len(in))

	for _, v := range in {
		res = append(res, &similarityMatch{
			ID:         v.ID,
			Rate:       v.Rate,
			Highlights: v.Highlights,
		})
	}

	return res
}
