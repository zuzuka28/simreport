package attribute

import (
	"fmt"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
)

func mapSearchAttributeRequestToModel(
	in *pb.SearchAttributeRequest,
) (model.AttributeQuery, error) {
	if in.GetAttribute() == "" {
		return model.AttributeQuery{}, fmt.Errorf("%w: empty id", model.ErrInvalid)
	}

	return model.AttributeQuery{
		ID: in.GetAttribute(),
	}, nil
}

func mapSearchAttributeToPb(
	in []model.Attribute,
) *pb.SearchAttributeResponse {
	attrs := make([]*pb.Attribute, 0, len(in))

	for _, v := range in {
		attrs = append(attrs, &pb.Attribute{
			Label: v.Label,
			Value: v.Value,
		})
	}

	return &pb.SearchAttributeResponse{
		Items: attrs,
	}
}
