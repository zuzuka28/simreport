package server

import (
	"context"
	openapi "simrep/api/rest/gen"
)

type (
	FileHandler interface {
		PostFilesUpload(
			ctx context.Context,
			request openapi.PostFilesUploadRequestObject,
		) (openapi.PostFilesUploadResponseObject, error)
		GetFilesFileId(
			ctx context.Context,
			request openapi.GetFilesFileIdRequestObject,
		) (openapi.GetFilesFileIdResponseObject, error)
	}

	SimilarityHandler interface {
		PostFilesCompare(
			ctx context.Context,
			request openapi.PostFilesCompareRequestObject,
		) (openapi.PostFilesCompareResponseObject, error)
		GetFilesFileIdCompareAll(
			ctx context.Context,
			request openapi.GetFilesFileIdCompareAllRequestObject,
		) (openapi.GetFilesFileIdCompareAllResponseObject, error)
	}
)
