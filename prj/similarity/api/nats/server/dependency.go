package server

import (
	"context"

	pb "github.com/zuzuka28/simreport/prj/similarity/pkg/pb/v1"
)

type (
	Metrics interface {
		IncNatsMicroRequest(op string, status string, size int, dur float64)
	}

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
