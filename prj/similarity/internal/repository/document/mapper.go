package document

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zuzuka28/simreport/prj/similarity/internal/model"

	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
)

var errInternal = errors.New("internal error")

const idPartsCount = 2

func parseFetchDocumentResponse(in *pb.FetchDocumentResponse) (model.Document, error) {
	raw := in.GetDocument()

	if in == nil {
		return model.Document{}, nil //nolint:exhaustruct
	}

	id := strings.Split(raw.GetId(), "_")

	if len(id) != idPartsCount {
		return model.Document{}, errInternal
	}

	parentID, ver := id[0], id[1]

	version, err := strconv.Atoi(ver)
	if err != nil {
		return model.Document{}, fmt.Errorf("parse document version: %w", err)
	}

	source := mapFileToModel(raw.GetSource())
	text := mapFileToModel(raw.GetText())

	imgs := make([]model.File, 0, len(raw.GetImages()))
	for _, v := range raw.GetImages() {
		imgs = append(imgs, mapFileToModel(v))
	}

	imgIDs := make([]string, 0, len(imgs))
	for _, v := range imgs {
		imgIDs = append(imgIDs, v.Sha256)
	}

	return model.Document{
		ParentID:    parentID,
		Name:        raw.GetName(),
		LastUpdated: time.Time{},
		Version:     version,
		GroupID:     raw.GetGroupIds(),
		SourceID:    source.Sha256,
		TextID:      text.Sha256,
		ImageIDs:    imgIDs,
		WithContent: raw.GetSource() != nil,
		Source:      source,
		Text:        text,
		Images:      imgs,
	}, nil
}

func mapFileToModel(
	in *pb.File,
) model.File {
	return model.File{
		Name:        in.GetFilename(),
		Content:     in.GetContent(),
		Sha256:      in.GetId(),
		LastUpdated: time.Time{},
	}
}

func isErr(in *pb.Error) error {
	if in == nil {
		return nil
	}

	status := in.GetStatus()
	switch status {
	case 404:
		return fmt.Errorf("%w: %s", model.ErrNotFound, in.GetDescription())

	case 400:
		return fmt.Errorf("%w: %s", model.ErrInvalid, in.GetDescription())

	default:
		return fmt.Errorf("%w: %s", errInternal, in.GetDescription())
	}
}

func mapDocumentQueryToPb(
	query model.DocumentQuery,
) *pb.FetchDocumentRequest {
	includeToModel := map[model.DocumentQueryInclude]pb.DocumentQueryInclude{
		model.DocumentQueryIncludeSource: pb.DocumentQueryInclude_DOCUMENT_QUERY_INCLUDE_SOURCE,
		model.DocumentQueryIncludeText:   pb.DocumentQueryInclude_DOCUMENT_QUERY_INCLUDE_TEXT,
		model.DocumentQueryIncludeImages: pb.DocumentQueryInclude_DOCUMENT_QUERY_INCLUDE_IMAGES,
	}

	include := make([]pb.DocumentQueryInclude, 0, len(query.Include))

	for _, v := range query.Include {
		include = append(include, includeToModel[v])
	}

	return &pb.FetchDocumentRequest{
		Id:          query.ID,
		WithContent: query.WithContent,
		Include:     include,
	}
}
