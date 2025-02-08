package similarityindexclient

import (
	"context"
	"fmt"
	"time"

	"github.com/zuzuka28/simreport/prj/similarity/internal/model"

	pb "github.com/zuzuka28/simreport/prj/similarityindex/pkg/pb/v1"
)

func (s *Repository) SearchSimilar(
	ctx context.Context,
	query model.SimilarityQuery,
) ([]*model.SimilarityMatch, error) {
	const op = "searchSimilar"

	t := time.Now()

	resp, err := s.cli.SearchSimilar(ctx, &pb.SearchSimilarRequest{
		Id: query.ID,
	})
	if err != nil {
		s.m.IncSimilarityIndexRequests(s.index, op, mapErrorToStatus(err), time.Since(t).Seconds())
		return nil, fmt.Errorf("do request: %w", mapErrorToModel(err))
	}

	s.m.IncSimilarityIndexRequests(s.index, op, "200", time.Since(t).Seconds())

	return parseSearchSimilarResponse(resp), nil
}
