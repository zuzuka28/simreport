package shingleindex

import (
	"github.com/zuzuka28/simreport/prj/shingleindex/internal/model"
	pb "github.com/zuzuka28/simreport/prj/shingleindex/pkg/pb/v1"
)

func mapDocumentToResponse(
	in []*model.DocumentSimilarMatch,
) *pb.SearchSimilarResponse {
	items := make([]*pb.SimilarityDocumentMatch, 0, len(in))

	for _, v := range in {
		items = append(items, &pb.SimilarityDocumentMatch{
			Id:            v.ID,
			Rate:          v.Rate,
			Highlights:    v.Highlights,
			SimilarImages: v.SimilarImages,
		})
	}

	return &pb.SearchSimilarResponse{
		Error:   nil,
		Matches: items,
	}
}
