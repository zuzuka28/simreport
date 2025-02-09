package shingleindex

import (
	"github.com/zuzuka28/simreport/lib/minhash"
	"github.com/zuzuka28/simreport/prj/shingleindex/internal/model"
)

func mapMinhashMatchToDocumentMatch(
	in *model.MinhashSimilarMatch,
) *model.DocumentSimilarMatch {
	return &model.DocumentSimilarMatch{
		ID:            in.DocumentID,
		Rate:          0,
		Highlights:    nil,
		SimilarImages: nil,
	}
}

func mapMinhashMatchesToDocumentMatches(
	in []*model.MinhashSimilarMatch,
) []*model.DocumentSimilarMatch {
	res := make([]*model.DocumentSimilarMatch, 0, len(in))

	for _, v := range in {
		res = append(res, mapMinhashMatchToDocumentMatch(v))
	}

	return res
}

func mapDocumentToMinhashSimilarQuery(
	in model.DocumentSimilarQuery,
) model.MinhashSimilarQuery {
	text := string(in.Item.Text)
	text = normalize(text)

	shingles := shingle(text, shingleSize)

	mh := minhash.New(
		permutations,
		hasher,
		seed,
	)

	for shingle := range shingles {
		mh.Push([]byte(shingle))
	}

	return model.MinhashSimilarQuery{
		Minhash: mh,
	}
}
