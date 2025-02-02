package document

import (
	"time"

	"github.com/zuzuka28/simreport/prj/document/internal/model"
	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapSearchRequestToQuery(
	in *pb.SearchRequest,
) model.DocumentSearchQuery {
	return model.DocumentSearchQuery{
		GroupID:  in.GetGroupIds(),
		Name:     in.GetName(),
		ParentID: in.GetParentId(),
		Version:  in.GetVersion(),
	}
}

func mapDocumentsToSearchResponse(
	in []model.Document,
) *pb.SearchDocumentResponse {
	docs := make([]*pb.DocumentSummary, 0, len(in))

	for _, v := range in {
		v := v

		docID := v.ID()

		docs = append(docs, &pb.DocumentSummary{
			GroupIds:    v.GroupID,
			Id:          docID,
			LastUpdated: v.LastUpdated.Format(time.RFC3339),
			Name:        v.Name,
			ParentId:    v.ParentID,
			Version:     int64(v.Version),
		})
	}

	return &pb.SearchDocumentResponse{
		Error:     nil,
		Documents: docs,
	}
}

func mapUploadRequestToCommand(
	in *pb.UploadDocumentRequest,
) model.DocumentSaveCommand {
	return model.DocumentSaveCommand{
		Item: model.Document{
			ParentID:    in.GetParentId(),
			Name:        in.GetFile().GetFilename(),
			LastUpdated: time.Time{},
			Version:     int(in.GetVersion()),
			GroupID:     in.GetGroupIds(),
			SourceID:    "",
			TextID:      "",
			ImageIDs:    nil,
			WithContent: false,
			Source: model.File{
				Name:        in.GetFile().GetFilename(),
				Content:     in.GetFile().GetContent(),
				Sha256:      in.GetFile().GetId(),
				LastUpdated: time.Time{},
			},
			Text:   model.File{}, //nolint:exhaustruct
			Images: nil,
		},
	}
}

func mapUploadCommandToResponse(
	doc *model.Document,
) *pb.UploadDocumentResponse {
	return &pb.UploadDocumentResponse{
		Error:      nil,
		DocumentId: doc.ID(),
	}
}

func mapFetchDocumentRequestToQuery(
	in *pb.FetchDocumentRequest,
) model.DocumentQuery {
	includes := make([]model.DocumentQueryInclude, 0, len(in.GetInclude()))

	includemap := map[pb.DocumentQueryInclude]model.DocumentQueryInclude{
		pb.DocumentQueryInclude_DOCUMENT_QUERY_INCLUDE_SOURCE: model.DocumentQueryIncludeSource,
		pb.DocumentQueryInclude_DOCUMENT_QUERY_INCLUDE_TEXT:   model.DocumentQueryIncludeText,
		pb.DocumentQueryInclude_DOCUMENT_QUERY_INCLUDE_IMAGES: model.DocumentQueryIncludeImages,
	}

	for _, v := range in.GetInclude() {
		includes = append(includes, includemap[v])
	}

	return model.DocumentQuery{
		ID:          in.GetId(),
		WithContent: in.GetWithContent(),
		Include:     includes,
	}
}

func mapFileToPb(
	in model.File,
) *pb.File {
	return &pb.File{
		Content:  in.Content,
		Id:       in.Sha256,
		Filename: in.Name,
	}
}

func mapDocumentToPb(
	in model.Document,
) *pb.Document {
	imgs := make([]*pb.File, 0, len(in.Images))

	for _, v := range in.Images {
		imgs = append(imgs, mapFileToPb(v))
	}

	return &pb.Document{
		Id:          in.ID(),
		Name:        in.Name,
		LastUpdated: timestamppb.New(in.LastUpdated),
		Version:     int64(in.Version),
		GroupIds:    in.GroupID,
		Source:      mapFileToPb(in.Source),
		Text:        mapFileToPb(in.Text),
		Images:      imgs,
	}
}

func mapFetchDocumentResponseToPb(
	in model.Document,
) *pb.FetchDocumentResponse {
	return &pb.FetchDocumentResponse{
		Error:    nil,
		Document: mapDocumentToPb(in),
	}
}
