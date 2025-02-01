package analyze

import (
	"fmt"
	"time"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
)

func mapSearchSimilarDocumentsRequestToModel(
	in *pb.DocumentId,
) (model.DocumentSimilarQuery, error) {
	if in.GetDocumentId() == "" {
		return model.DocumentSimilarQuery{}, fmt.Errorf("%w: empty id", model.ErrInvalid)
	}

	return model.DocumentSimilarQuery{
		ID:   in.GetDocumentId(),
		Item: model.Document{}, //nolint:exhaustruct
	}, nil
}

func mapDocumentSimilarMatchToPb(
	in *model.DocumentSimilarMatch,
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

func mapSearchSimilarDocumentsResponseToPb(
	in []*model.DocumentSimilarMatch,
) *pb.SearchSimilarDocumentsResponse {
	docs := make([]*pb.SimilarityDocumentMatch, 0, len(in))

	for _, v := range in {
		docs = append(docs, mapDocumentSimilarMatchToPb(v))
	}

	return &pb.SearchSimilarDocumentsResponse{
		Error:     nil,
		Documents: docs,
	}
}

func mapSearchSimilaritySearchHistoryRequestToModel(
	in *pb.SearchSimilaritySearchHistoryRequest,
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

func mapSearchSimilaritySearchHistoryResponseToPb(
	in *model.SimilarityHistoryList,
) *pb.SearchSimilaritySearchHistoryResponse {
	items := make([]*pb.SimilaritySearchHistory, 0, len(in.Items))

	for _, v := range in.Items {
		items = append(items, mapSimilaritySearchHistoryToPb(v))
	}

	return &pb.SearchSimilaritySearchHistoryResponse{
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
