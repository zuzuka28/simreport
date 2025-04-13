package server

import (
	"context"

	pb "github.com/zuzuka28/simreport/prj/similarityindex/pkg/pb/v1"
)

type (
	Metrics interface {
		IncNatsMicroRequest(op string, status string, size int, dur float64)
	}

	Handler interface {
		SearchSimilar(
			ctx context.Context,
			params *pb.SearchSimilarRequest,
		) (*pb.SearchSimilarResponse, error)
	}
)
