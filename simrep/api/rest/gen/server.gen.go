// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package openapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// AnalyzedDocumentMatch defines model for AnalyzedDocumentMatch.
type AnalyzedDocumentMatch struct {
	// Highlights Список фрагментов текста, которые совпадают в документах.
	Highlights *[]string `json:"highlights,omitempty"`

	// Id Уникальный идентификатор документа.
	Id *string `json:"id,omitempty"`

	// Rate Коэффициент похожести (например, от 0 до 1).
	Rate *float32 `json:"rate,omitempty"`

	// SimilarImages Список идентификаторов похожих изображений.
	SimilarImages *[]string `json:"similarImages,omitempty"`
}

// DocumentSummary defines model for DocumentSummary.
type DocumentSummary struct {
	// Id Идентификатор документа
	Id *string `json:"id,omitempty"`

	// LastUpdated Дата обновления документа
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`

	// Name Имя документа
	Name *string `json:"name,omitempty"`
}

// SearchRequest defines model for SearchRequest.
type SearchRequest struct {
	// Name Имя документа для поиска
	Name *string `json:"name,omitempty"`
}

// UploadRequest defines model for UploadRequest.
type UploadRequest struct {
	// Document Документ для загрузки
	Document *openapi_types.File `json:"document,omitempty"`
}

// DocumentId defines model for DocumentId.
type DocumentId = string

// BadRequest defines model for BadRequest.
type BadRequest struct {
	Error *string `json:"error,omitempty"`
}

// DocumentNotFound defines model for DocumentNotFound.
type DocumentNotFound struct {
	Error *string `json:"error,omitempty"`
}

// SearchResult defines model for SearchResult.
type SearchResult struct {
	Documents *[]DocumentSummary `json:"documents,omitempty"`
}

// ServerError defines model for ServerError.
type ServerError struct {
	Error *string `json:"error,omitempty"`
}

// SimilaritySearchResult defines model for SimilaritySearchResult.
type SimilaritySearchResult struct {
	Documents *[]AnalyzedDocumentMatch `json:"documents,omitempty"`
}

// UploadSuccess defines model for UploadSuccess.
type UploadSuccess struct {
	// DocumentID Уникальный идентификатор загруженного документа
	DocumentID *string `json:"documentID,omitempty"`
}

// PostAnalyzeSimilarMultipartRequestBody defines body for PostAnalyzeSimilar for multipart/form-data ContentType.
type PostAnalyzeSimilarMultipartRequestBody = UploadRequest

// PostDocumentUploadMultipartRequestBody defines body for PostDocumentUpload for multipart/form-data ContentType.
type PostDocumentUploadMultipartRequestBody = UploadRequest

// PostDocumentsSearchJSONRequestBody defines body for PostDocumentsSearch for application/json ContentType.
type PostDocumentsSearchJSONRequestBody = SearchRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Поиск подожих документов
	// (POST /analyze/similar)
	PostAnalyzeSimilar(w http.ResponseWriter, r *http.Request)
	// Загрузка документа
	// (POST /document/upload)
	PostDocumentUpload(w http.ResponseWriter, r *http.Request)
	// Скачать документ
	// (GET /document/{document_id}/download)
	GetDocumentDocumentIdDownload(w http.ResponseWriter, r *http.Request, documentId DocumentId)
	// Поиск документов по имени
	// (POST /documents/search)
	PostDocumentsSearch(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// PostAnalyzeSimilar operation middleware
func (siw *ServerInterfaceWrapper) PostAnalyzeSimilar(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostAnalyzeSimilar(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostDocumentUpload operation middleware
func (siw *ServerInterfaceWrapper) PostDocumentUpload(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostDocumentUpload(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetDocumentDocumentIdDownload operation middleware
func (siw *ServerInterfaceWrapper) GetDocumentDocumentIdDownload(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "document_id" -------------
	var documentId DocumentId

	err = runtime.BindStyledParameterWithOptions("simple", "document_id", mux.Vars(r)["document_id"], &documentId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "document_id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetDocumentDocumentIdDownload(w, r, documentId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostDocumentsSearch operation middleware
func (siw *ServerInterfaceWrapper) PostDocumentsSearch(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostDocumentsSearch(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{})
}

type GorillaServerOptions struct {
	BaseURL          string
	BaseRouter       *mux.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r *mux.Router) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r *mux.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options GorillaServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = mux.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.HandleFunc(options.BaseURL+"/analyze/similar", wrapper.PostAnalyzeSimilar).Methods("POST")

	r.HandleFunc(options.BaseURL+"/document/upload", wrapper.PostDocumentUpload).Methods("POST")

	r.HandleFunc(options.BaseURL+"/document/{document_id}/download", wrapper.GetDocumentDocumentIdDownload).Methods("GET")

	r.HandleFunc(options.BaseURL+"/documents/search", wrapper.PostDocumentsSearch).Methods("POST")

	return r
}

type BadRequestJSONResponse struct {
	Error *string `json:"error,omitempty"`
}

type DocumentNotFoundJSONResponse struct {
	Error *string `json:"error,omitempty"`
}

type SearchResultJSONResponse struct {
	Documents *[]DocumentSummary `json:"documents,omitempty"`
}

type ServerErrorJSONResponse struct {
	Error *string `json:"error,omitempty"`
}

type SimilaritySearchResultJSONResponse struct {
	Documents *[]AnalyzedDocumentMatch `json:"documents,omitempty"`
}

type UploadSuccessJSONResponse struct {
	// DocumentID Уникальный идентификатор загруженного документа
	DocumentID *string `json:"documentID,omitempty"`
}

type PostAnalyzeSimilarRequestObject struct {
	Body *multipart.Reader
}

type PostAnalyzeSimilarResponseObject interface {
	VisitPostAnalyzeSimilarResponse(w http.ResponseWriter) error
}

type PostAnalyzeSimilar200JSONResponse struct {
	SimilaritySearchResultJSONResponse
}

func (response PostAnalyzeSimilar200JSONResponse) VisitPostAnalyzeSimilarResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type PostAnalyzeSimilar400JSONResponse struct{ BadRequestJSONResponse }

func (response PostAnalyzeSimilar400JSONResponse) VisitPostAnalyzeSimilarResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type PostAnalyzeSimilar500JSONResponse struct{ ServerErrorJSONResponse }

func (response PostAnalyzeSimilar500JSONResponse) VisitPostAnalyzeSimilarResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type PostDocumentUploadRequestObject struct {
	Body *multipart.Reader
}

type PostDocumentUploadResponseObject interface {
	VisitPostDocumentUploadResponse(w http.ResponseWriter) error
}

type PostDocumentUpload200JSONResponse struct{ UploadSuccessJSONResponse }

func (response PostDocumentUpload200JSONResponse) VisitPostDocumentUploadResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type PostDocumentUpload400JSONResponse struct{ BadRequestJSONResponse }

func (response PostDocumentUpload400JSONResponse) VisitPostDocumentUploadResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type GetDocumentDocumentIdDownloadRequestObject struct {
	DocumentId DocumentId `json:"document_id"`
}

type GetDocumentDocumentIdDownloadResponseObject interface {
	VisitGetDocumentDocumentIdDownloadResponse(w http.ResponseWriter) error
}

type GetDocumentDocumentIdDownload200ResponseHeaders struct {
	ContentDisposition string
}

type GetDocumentDocumentIdDownload200ApplicationoctetStreamResponse struct {
	Body          io.Reader
	Headers       GetDocumentDocumentIdDownload200ResponseHeaders
	ContentLength int64
}

func (response GetDocumentDocumentIdDownload200ApplicationoctetStreamResponse) VisitGetDocumentDocumentIdDownloadResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/octet-stream")
	if response.ContentLength != 0 {
		w.Header().Set("Content-Length", fmt.Sprint(response.ContentLength))
	}
	w.Header().Set("Content-Disposition", fmt.Sprint(response.Headers.ContentDisposition))
	w.WriteHeader(200)

	if closer, ok := response.Body.(io.ReadCloser); ok {
		defer closer.Close()
	}
	_, err := io.Copy(w, response.Body)
	return err
}

type GetDocumentDocumentIdDownload404JSONResponse struct{ DocumentNotFoundJSONResponse }

func (response GetDocumentDocumentIdDownload404JSONResponse) VisitGetDocumentDocumentIdDownloadResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type PostDocumentsSearchRequestObject struct {
	Body *PostDocumentsSearchJSONRequestBody
}

type PostDocumentsSearchResponseObject interface {
	VisitPostDocumentsSearchResponse(w http.ResponseWriter) error
}

type PostDocumentsSearch200JSONResponse struct{ SearchResultJSONResponse }

func (response PostDocumentsSearch200JSONResponse) VisitPostDocumentsSearchResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type PostDocumentsSearch400JSONResponse struct{ BadRequestJSONResponse }

func (response PostDocumentsSearch400JSONResponse) VisitPostDocumentsSearchResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type PostDocumentsSearch500JSONResponse struct{ ServerErrorJSONResponse }

func (response PostDocumentsSearch500JSONResponse) VisitPostDocumentsSearchResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Поиск подожих документов
	// (POST /analyze/similar)
	PostAnalyzeSimilar(ctx context.Context, request PostAnalyzeSimilarRequestObject) (PostAnalyzeSimilarResponseObject, error)
	// Загрузка документа
	// (POST /document/upload)
	PostDocumentUpload(ctx context.Context, request PostDocumentUploadRequestObject) (PostDocumentUploadResponseObject, error)
	// Скачать документ
	// (GET /document/{document_id}/download)
	GetDocumentDocumentIdDownload(ctx context.Context, request GetDocumentDocumentIdDownloadRequestObject) (GetDocumentDocumentIdDownloadResponseObject, error)
	// Поиск документов по имени
	// (POST /documents/search)
	PostDocumentsSearch(ctx context.Context, request PostDocumentsSearchRequestObject) (PostDocumentsSearchResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// PostAnalyzeSimilar operation middleware
func (sh *strictHandler) PostAnalyzeSimilar(w http.ResponseWriter, r *http.Request) {
	var request PostAnalyzeSimilarRequestObject

	if reader, err := r.MultipartReader(); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode multipart body: %w", err))
		return
	} else {
		request.Body = reader
	}

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PostAnalyzeSimilar(ctx, request.(PostAnalyzeSimilarRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostAnalyzeSimilar")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PostAnalyzeSimilarResponseObject); ok {
		if err := validResponse.VisitPostAnalyzeSimilarResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// PostDocumentUpload operation middleware
func (sh *strictHandler) PostDocumentUpload(w http.ResponseWriter, r *http.Request) {
	var request PostDocumentUploadRequestObject

	if reader, err := r.MultipartReader(); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode multipart body: %w", err))
		return
	} else {
		request.Body = reader
	}

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PostDocumentUpload(ctx, request.(PostDocumentUploadRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostDocumentUpload")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PostDocumentUploadResponseObject); ok {
		if err := validResponse.VisitPostDocumentUploadResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetDocumentDocumentIdDownload operation middleware
func (sh *strictHandler) GetDocumentDocumentIdDownload(w http.ResponseWriter, r *http.Request, documentId DocumentId) {
	var request GetDocumentDocumentIdDownloadRequestObject

	request.DocumentId = documentId

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetDocumentDocumentIdDownload(ctx, request.(GetDocumentDocumentIdDownloadRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetDocumentDocumentIdDownload")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetDocumentDocumentIdDownloadResponseObject); ok {
		if err := validResponse.VisitGetDocumentDocumentIdDownloadResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// PostDocumentsSearch operation middleware
func (sh *strictHandler) PostDocumentsSearch(w http.ResponseWriter, r *http.Request) {
	var request PostDocumentsSearchRequestObject

	var body PostDocumentsSearchJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PostDocumentsSearch(ctx, request.(PostDocumentsSearchRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostDocumentsSearch")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PostDocumentsSearchResponseObject); ok {
		if err := validResponse.VisitPostDocumentsSearchResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}
