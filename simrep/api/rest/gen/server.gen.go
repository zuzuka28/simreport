// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package openapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
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

// SimilaritySearchHistory defines model for SimilaritySearchHistory.
type SimilaritySearchHistory struct {
	Date       *time.Time               `json:"date,omitempty"`
	DocumentID *string                  `json:"documentID,omitempty"`
	Id         *string                  `json:"id,omitempty"`
	Matches    *[]AnalyzedDocumentMatch `json:"matches,omitempty"`
}

// SimilaritySearchHistoryRequest defines model for SimilaritySearchHistoryRequest.
type SimilaritySearchHistoryRequest struct {
	DateFrom   *time.Time `json:"dateFrom,omitempty"`
	DateTo     *time.Time `json:"dateTo,omitempty"`
	DocumentID *string    `json:"documentID,omitempty"`
	Limit      *int       `json:"limit,omitempty"`
	Offset     *int       `json:"offset,omitempty"`
}

// UploadRequest defines model for UploadRequest.
type UploadRequest struct {
	// FileID Уникальный идентификатор файла
	FileID string `json:"fileID"`

	// GroupID Уникальный идентификатор группы
	GroupID *string `json:"groupID,omitempty"`

	// ParentID Уникальный идентификатор родительского документа
	ParentID *string `json:"parentID,omitempty"`

	// Version Версия документа
	Version *int `json:"version,omitempty"`
}

// DocumentId defines model for DocumentId.
type DocumentId = string

// BadRequest defines model for BadRequest.
type BadRequest struct {
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

// SimilaritySearchHistoryResult defines model for SimilaritySearchHistoryResult.
type SimilaritySearchHistoryResult struct {
	Count     *int                       `json:"count,omitempty"`
	Documents *[]SimilaritySearchHistory `json:"documents,omitempty"`
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

// PostAnalyzeHistoryJSONRequestBody defines body for PostAnalyzeHistory for application/json ContentType.
type PostAnalyzeHistoryJSONRequestBody = SimilaritySearchHistoryRequest

// PostDocumentSearchJSONRequestBody defines body for PostDocumentSearch for application/json ContentType.
type PostDocumentSearchJSONRequestBody = SearchRequest

// PostDocumentUploadJSONRequestBody defines body for PostDocumentUpload for application/json ContentType.
type PostDocumentUploadJSONRequestBody = UploadRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	//  История поиска подожих документов
	// (POST /analyze/history)
	PostAnalyzeHistory(w http.ResponseWriter, r *http.Request)
	// Поиск подожих документов
	// (GET /analyze/{document_id}/similar)
	GetAnalyzeDocumentIdSimilar(w http.ResponseWriter, r *http.Request, documentId DocumentId)
	// Поиск документов по имени
	// (POST /document/search)
	PostDocumentSearch(w http.ResponseWriter, r *http.Request)
	// Загрузка документа
	// (POST /document/upload)
	PostDocumentUpload(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// PostAnalyzeHistory operation middleware
func (siw *ServerInterfaceWrapper) PostAnalyzeHistory(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostAnalyzeHistory(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetAnalyzeDocumentIdSimilar operation middleware
func (siw *ServerInterfaceWrapper) GetAnalyzeDocumentIdSimilar(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "document_id" -------------
	var documentId DocumentId

	err = runtime.BindStyledParameterWithOptions("simple", "document_id", mux.Vars(r)["document_id"], &documentId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "document_id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetAnalyzeDocumentIdSimilar(w, r, documentId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostDocumentSearch operation middleware
func (siw *ServerInterfaceWrapper) PostDocumentSearch(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostDocumentSearch(w, r)
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

	r.HandleFunc(options.BaseURL+"/analyze/history", wrapper.PostAnalyzeHistory).Methods("POST")

	r.HandleFunc(options.BaseURL+"/analyze/{document_id}/similar", wrapper.GetAnalyzeDocumentIdSimilar).Methods("GET")

	r.HandleFunc(options.BaseURL+"/document/search", wrapper.PostDocumentSearch).Methods("POST")

	r.HandleFunc(options.BaseURL+"/document/upload", wrapper.PostDocumentUpload).Methods("POST")

	return r
}

type BadRequestJSONResponse struct {
	Error *string `json:"error,omitempty"`
}

type SearchResultJSONResponse struct {
	Documents *[]DocumentSummary `json:"documents,omitempty"`
}

type ServerErrorJSONResponse struct {
	Error *string `json:"error,omitempty"`
}

type SimilaritySearchHistoryResultJSONResponse struct {
	Count     *int                       `json:"count,omitempty"`
	Documents *[]SimilaritySearchHistory `json:"documents,omitempty"`
}

type SimilaritySearchResultJSONResponse struct {
	Documents *[]AnalyzedDocumentMatch `json:"documents,omitempty"`
}

type UploadSuccessJSONResponse struct {
	// DocumentID Уникальный идентификатор загруженного документа
	DocumentID *string `json:"documentID,omitempty"`
}

type PostAnalyzeHistoryRequestObject struct {
	Body *PostAnalyzeHistoryJSONRequestBody
}

type PostAnalyzeHistoryResponseObject interface {
	VisitPostAnalyzeHistoryResponse(w http.ResponseWriter) error
}

type PostAnalyzeHistory200JSONResponse struct {
	SimilaritySearchHistoryResultJSONResponse
}

func (response PostAnalyzeHistory200JSONResponse) VisitPostAnalyzeHistoryResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type PostAnalyzeHistory400JSONResponse struct{ BadRequestJSONResponse }

func (response PostAnalyzeHistory400JSONResponse) VisitPostAnalyzeHistoryResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type PostAnalyzeHistory500JSONResponse struct{ ServerErrorJSONResponse }

func (response PostAnalyzeHistory500JSONResponse) VisitPostAnalyzeHistoryResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type GetAnalyzeDocumentIdSimilarRequestObject struct {
	DocumentId DocumentId `json:"document_id"`
}

type GetAnalyzeDocumentIdSimilarResponseObject interface {
	VisitGetAnalyzeDocumentIdSimilarResponse(w http.ResponseWriter) error
}

type GetAnalyzeDocumentIdSimilar200JSONResponse struct {
	SimilaritySearchResultJSONResponse
}

func (response GetAnalyzeDocumentIdSimilar200JSONResponse) VisitGetAnalyzeDocumentIdSimilarResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetAnalyzeDocumentIdSimilar400JSONResponse struct{ BadRequestJSONResponse }

func (response GetAnalyzeDocumentIdSimilar400JSONResponse) VisitGetAnalyzeDocumentIdSimilarResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type GetAnalyzeDocumentIdSimilar500JSONResponse struct{ ServerErrorJSONResponse }

func (response GetAnalyzeDocumentIdSimilar500JSONResponse) VisitGetAnalyzeDocumentIdSimilarResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type PostDocumentSearchRequestObject struct {
	Body *PostDocumentSearchJSONRequestBody
}

type PostDocumentSearchResponseObject interface {
	VisitPostDocumentSearchResponse(w http.ResponseWriter) error
}

type PostDocumentSearch200JSONResponse struct{ SearchResultJSONResponse }

func (response PostDocumentSearch200JSONResponse) VisitPostDocumentSearchResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type PostDocumentSearch400JSONResponse struct{ BadRequestJSONResponse }

func (response PostDocumentSearch400JSONResponse) VisitPostDocumentSearchResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type PostDocumentSearch500JSONResponse struct{ ServerErrorJSONResponse }

func (response PostDocumentSearch500JSONResponse) VisitPostDocumentSearchResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type PostDocumentUploadRequestObject struct {
	Body *PostDocumentUploadJSONRequestBody
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

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	//  История поиска подожих документов
	// (POST /analyze/history)
	PostAnalyzeHistory(ctx context.Context, request PostAnalyzeHistoryRequestObject) (PostAnalyzeHistoryResponseObject, error)
	// Поиск подожих документов
	// (GET /analyze/{document_id}/similar)
	GetAnalyzeDocumentIdSimilar(ctx context.Context, request GetAnalyzeDocumentIdSimilarRequestObject) (GetAnalyzeDocumentIdSimilarResponseObject, error)
	// Поиск документов по имени
	// (POST /document/search)
	PostDocumentSearch(ctx context.Context, request PostDocumentSearchRequestObject) (PostDocumentSearchResponseObject, error)
	// Загрузка документа
	// (POST /document/upload)
	PostDocumentUpload(ctx context.Context, request PostDocumentUploadRequestObject) (PostDocumentUploadResponseObject, error)
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

// PostAnalyzeHistory operation middleware
func (sh *strictHandler) PostAnalyzeHistory(w http.ResponseWriter, r *http.Request) {
	var request PostAnalyzeHistoryRequestObject

	var body PostAnalyzeHistoryJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PostAnalyzeHistory(ctx, request.(PostAnalyzeHistoryRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostAnalyzeHistory")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PostAnalyzeHistoryResponseObject); ok {
		if err := validResponse.VisitPostAnalyzeHistoryResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetAnalyzeDocumentIdSimilar operation middleware
func (sh *strictHandler) GetAnalyzeDocumentIdSimilar(w http.ResponseWriter, r *http.Request, documentId DocumentId) {
	var request GetAnalyzeDocumentIdSimilarRequestObject

	request.DocumentId = documentId

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetAnalyzeDocumentIdSimilar(ctx, request.(GetAnalyzeDocumentIdSimilarRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetAnalyzeDocumentIdSimilar")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetAnalyzeDocumentIdSimilarResponseObject); ok {
		if err := validResponse.VisitGetAnalyzeDocumentIdSimilarResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// PostDocumentSearch operation middleware
func (sh *strictHandler) PostDocumentSearch(w http.ResponseWriter, r *http.Request) {
	var request PostDocumentSearchRequestObject

	var body PostDocumentSearchJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PostDocumentSearch(ctx, request.(PostDocumentSearchRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostDocumentSearch")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PostDocumentSearchResponseObject); ok {
		if err := validResponse.VisitPostDocumentSearchResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// PostDocumentUpload operation middleware
func (sh *strictHandler) PostDocumentUpload(w http.ResponseWriter, r *http.Request) {
	var request PostDocumentUploadRequestObject

	var body PostDocumentUploadJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

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
