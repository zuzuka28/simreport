package server

import (
	"context"

	pb "github.com/zuzuka28/simreport/prj/similarity/pkg/pb/v1"
)

type (
	SimilarityHandler interface {
		SearchSimilar(
			ctx context.Context,
			params *pb.DocumentId,
		) (*pb.SearchSimilarResponse, error)
		SearchSimilarityHistory(
			ctx context.Context,
			params *pb.SearchSimilarityHistoryRequest,
		) (*pb.SearchSimilarityHistoryResponse, error)
	}
)
