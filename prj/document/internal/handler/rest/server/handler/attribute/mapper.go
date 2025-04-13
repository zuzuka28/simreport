package attribute

import (
	"errors"

	openapi "github.com/zuzuka28/simreport/prj/document/internal/handler/rest/gen"
	"github.com/zuzuka28/simreport/prj/document/internal/model"
)

var errNoBody = errors.New("no body")

func mapAttributeRequestToQuery(
	in openapi.PostAttributeRequestObject,
) (model.AttributeQuery, error) {
	if in.Body == nil {
		return model.AttributeQuery{}, errNoBody
	}

	return model.AttributeQuery{
		ID: in.Body.Attribute,
	}, nil
}

func mapDocumentsToSearchResponse(
	in []model.Attribute,
) openapi.PostAttributeResponseObject {
	items := make([]openapi.Attribute, 0, len(in))

	for _, v := range in {
		v := v

		items = append(items, openapi.Attribute{
			Label: v.Label,
			Value: v.Value,
		})
	}

	return openapi.PostAttribute200JSONResponse{
		AttributeResultJSONResponse: openapi.AttributeResultJSONResponse{
			Items: &items,
		},
	}
}
