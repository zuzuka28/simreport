package fulltextindexclient

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/similarity/internal/model"

	pb "github.com/zuzuka28/simreport/prj/similarityindex/pkg/pb/v1"
)

func (s *Repository) SearchSimilar(
	ctx context.Context,
	query model.SimilarityQuery,
) ([]*model.SimilarityMatch, error) {
	resp, err := s.cli.SearchSimilar(ctx, &pb.SearchSimilarRequest{
		Id: query.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	if err := isErr(resp.GetError()); err != nil {
		return nil, err
	}

	return parseSearchSimilarResponse(resp), nil
}
