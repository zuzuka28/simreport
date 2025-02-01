package similarity

import (
	"fmt"
	"time"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
)

func mapSearchSimilarRequestToModel(
	in *pb.DocumentId,
) (model.SimilarityQuery, error) {
	if in.GetDocumentId() == "" {
		return model.SimilarityQuery{}, fmt.Errorf("%w: empty id", model.ErrInvalid)
	}

	return model.SimilarityQuery{
		ID:   in.GetDocumentId(),
		Item: model.Document{}, //nolint:exhaustruct
	}, nil
}

func mapDocumentSimilarMatchToPb(
	in *model.SimilarityMatch,
) *pb.SimilarityDocumentMatch {
	if in == nil {
		return nil
	}

	return &pb.SimilarityDocumentMatch{
		Id:            in.ID,
		Rate:          float32(in.Rate),
		Highlights:    in.Highlights,
		SimilarImages: in.SimilarImages,
	}
}

func mapSearchSimilarResponseToPb(
	in []*model.SimilarityMatch,
) *pb.SearchSimilarResponse {
	docs := make([]*pb.SimilarityDocumentMatch, 0, len(in))

	for _, v := range in {
		docs = append(docs, mapDocumentSimilarMatchToPb(v))
	}

	return &pb.SearchSimilarResponse{
		Error:     nil,
		Documents: docs,
	}
}

func mapSearchSimilarityHistoryRequestToModel(
	in *pb.SearchSimilarityHistoryRequest,
) (model.SimilarityHistoryQuery, error) {
	tFrom, err := time.Parse(time.RFC3339, in.GetDateFrom())
	if err != nil {
		return model.SimilarityHistoryQuery{}, fmt.Errorf("invalid DateFrom format: %w", err)
	}

	tTo, err := time.Parse(time.RFC3339, in.GetDateTo())
	if err != nil {
		return model.SimilarityHistoryQuery{}, fmt.Errorf("invalid DateTo format: %w", err)
	}

	return model.SimilarityHistoryQuery{
		DocumentID: in.GetDocumentId(),
		Limit:      int(in.GetLimit()),
		Offset:     int(in.GetOffset()),
		DateFrom:   tFrom,
		DateTo:     tTo,
	}, nil
}

func mapSearchSimilarityHistoryResponseToPb(
	in *model.SimilarityHistoryList,
) *pb.SearchSimilarityHistoryResponse {
	items := make([]*pb.SimilaritySearchHistory, 0, len(in.Items))

	for _, v := range in.Items {
		items = append(items, mapSimilaritySearchHistoryToPb(v))
	}

	return &pb.SearchSimilarityHistoryResponse{
		Error:     nil,
		Documents: items,
		Count:     int64(in.Count),
	}
}

func mapSimilaritySearchHistoryToPb(in *model.SimilarityHistory) *pb.SimilaritySearchHistory {
	docs := make([]*pb.SimilarityDocumentMatch, 0, len(in.Matches))

	for _, v := range in.Matches {
		docs = append(docs, mapDocumentSimilarMatchToPb(v))
	}

	return &pb.SimilaritySearchHistory{
		Date:       in.Date.Format(time.RFC3339),
		DocumentId: in.DocumentID,
		Id:         in.ID,
		Matches:    docs,
	}
}
