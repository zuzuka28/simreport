package server

import (
	"context"

	pb "github.com/zuzuka28/simreport/prj/similarityindex/pkg/pb/v1"
)

type (
	Handler interface {
		SearchSimilar(
			ctx context.Context,
			params *pb.SearchSimilarRequest,
		) (*pb.SearchSimilarResponse, error)
	}
)
