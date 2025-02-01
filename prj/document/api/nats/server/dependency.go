package server

import (
	"context"

	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
)

type (
	DocumentHandler interface {
		FetchDocument(
			ctx context.Context,
			params *pb.FetchDocumentRequest,
		) (*pb.FetchDocumentResponse, error)
		SearchDocument(
			ctx context.Context,
			params *pb.SearchRequest,
		) (*pb.SearchDocumentResponse, error)
		UploadDocument(
			ctx context.Context,
			params *pb.UploadRequest,
		) (*pb.UploadDocumentResponse, error)
	}

	AttributeHandler interface {
		SearchAttribute(
			ctx context.Context,
			params *pb.SearchAttributeRequest,
		) (*pb.SearchAttributeResponse, error)
	}

	AnalyzeHandler interface {
		SearchSimilarDocuments(
			ctx context.Context,
			params *pb.DocumentId,
		) (*pb.SearchSimilarDocumentsResponse, error)
		SearchSimilaritySearchHistory(
			ctx context.Context,
			params *pb.SearchSimilaritySearchHistoryRequest,
		) (*pb.SearchSimilaritySearchHistoryResponse, error)
	}
)
