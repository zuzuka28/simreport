package server

import (
	"context"

	pb "github.com/zuzuka28/simreport/prj/document/pkg/pb/v1"
)

type (
	Metrics interface {
		IncNatsMicroRequest(op string, status string, size int, dur float64)
	}

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
			params *pb.UploadDocumentRequest,
		) (*pb.UploadDocumentResponse, error)
	}

	AttributeHandler interface {
		SearchAttribute(
			ctx context.Context,
			params *pb.SearchAttributeRequest,
		) (*pb.SearchAttributeResponse, error)
	}
)
