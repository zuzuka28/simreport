// Code generated by protoc-gen-natsmicro. DO NOT EDIT.

package document

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"google.golang.org/protobuf/proto"
)

type Handler func(ctx context.Context, req micro.Request)

type ErrorHandler func(ctx context.Context, msg micro.Request, err error)

type Middleware func(Handler) Handler

// DocumentServiceServerInterface is the server API for DocumentService service.
type DocumentServiceServerInterface interface {
	FetchDocument(context.Context, *FetchDocumentRequest) (*FetchDocumentResponse, error)
	UploadDocument(context.Context, *UploadDocumentRequest) (*UploadDocumentResponse, error)
	SearchAttribute(context.Context, *SearchAttributeRequest) (*SearchAttributeResponse, error)
	SearchDocument(context.Context, *SearchRequest) (*SearchDocumentResponse, error)
	mustEmbedUnimplementedGreeterServer()
}

type UnimplementedDocumentServiceServer struct{}

func (UnimplementedDocumentServiceServer) FetchDocument(context.Context, *FetchDocumentRequest) (*FetchDocumentResponse, error) {
	return nil, errors.New("method FetchDocument not implemented")
}
func (UnimplementedDocumentServiceServer) UploadDocument(context.Context, *UploadDocumentRequest) (*UploadDocumentResponse, error) {
	return nil, errors.New("method UploadDocument not implemented")
}
func (UnimplementedDocumentServiceServer) SearchAttribute(context.Context, *SearchAttributeRequest) (*SearchAttributeResponse, error) {
	return nil, errors.New("method SearchAttribute not implemented")
}
func (UnimplementedDocumentServiceServer) SearchDocument(context.Context, *SearchRequest) (*SearchDocumentResponse, error) {
	return nil, errors.New("method SearchDocument not implemented")
}

func (UnimplementedDocumentServiceServer) mustEmbedUnimplementedDocumentServiceServer() {}

func (UnimplementedDocumentServiceServer) testEmbeddedByValue() {}

type UnsafeDocumentServiceServer interface {
	mustEmbedUnimplementedGreeterServer()
}

type DocumentServiceServerConfig struct {
	micro.Config
	RequestTimeout       time.Duration
	Middleware           Middleware
	RequestErrorHandler  ErrorHandler
	ResponseErrorHandler ErrorHandler
}

type DocumentServiceServer struct {
	nc   *nats.Conn
	impl DocumentServiceServerInterface
	cfg  DocumentServiceServerConfig
	done chan struct{}

	requestErrorHandlerFunc  ErrorHandler
	responseErrorHandlerFunc ErrorHandler
}

// NewDocumentServiceServer  creates a new NATS microservice server.
func NewDocumentServiceServer(
	cfg DocumentServiceServerConfig,
	nc *nats.Conn,
	impl DocumentServiceServerInterface,
) *DocumentServiceServer {
	if t, ok := impl.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}

	if cfg.RequestTimeout == 0 {
		cfg.RequestTimeout = time.Second * 60
	}

	if cfg.RequestErrorHandler == nil {
		cfg.RequestErrorHandler = func(_ context.Context, req micro.Request, _ error) {
			req.Error("500", "unproccessable request", nil)
		}
	}

	if cfg.ResponseErrorHandler == nil {
		cfg.ResponseErrorHandler = func(_ context.Context, req micro.Request, _ error) {
			req.Error("500", "internal server error", nil)
		}
	}

	return &DocumentServiceServer{
		nc:                       nc,
		impl:                     impl,
		cfg:                      cfg,
		done:                     make(chan struct{}),
		requestErrorHandlerFunc:  cfg.RequestErrorHandler,
		responseErrorHandlerFunc: cfg.ResponseErrorHandler,
	}
}

// Start starts the microservice and blocking until application context done.
func (s *DocumentServiceServer) Start(ctx context.Context) error {
	srv, err := micro.AddService(s.nc, s.cfg.Config)
	if err != nil {
		return fmt.Errorf("failed to start microservice: %w", err)
	}

	defer srv.Stop()

	group := srv.AddGroup(s.cfg.Config.Name)

	// Register handlers
	if err := group.AddEndpoint("fetch_document", s.toMicroHandler(s.handleFetchDocument)); err != nil {
		return fmt.Errorf("failed to add endpoint FetchDocument: %w", err)
	}
	if err := group.AddEndpoint("upload_document", s.toMicroHandler(s.handleUploadDocument)); err != nil {
		return fmt.Errorf("failed to add endpoint UploadDocument: %w", err)
	}
	if err := group.AddEndpoint("search_attribute", s.toMicroHandler(s.handleSearchAttribute)); err != nil {
		return fmt.Errorf("failed to add endpoint SearchAttribute: %w", err)
	}
	if err := group.AddEndpoint("search_document", s.toMicroHandler(s.handleSearchDocument)); err != nil {
		return fmt.Errorf("failed to add endpoint SearchDocument: %w", err)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.done:
		return nil
	}
}

func (s *DocumentServiceServer) Stop() error {
	s.done <- struct{}{}
	return nil
}

func (s *DocumentServiceServer) toMicroHandler(h Handler) micro.HandlerFunc {
	return func(req micro.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.RequestTimeout)
		defer cancel()

		if s.cfg.Middleware != nil {
			h = s.cfg.Middleware(h)
		}

		h(ctx, req)
	}
}

func (s *DocumentServiceServer) handleFetchDocument(
	ctx context.Context,
	req micro.Request,
) {
	msg := new(FetchDocumentRequest)

	if err := proto.Unmarshal(req.Data(), msg); err != nil {
		s.requestErrorHandlerFunc(ctx, req, err)
		return
	}

	res, err := s.impl.FetchDocument(ctx, msg)
	if err != nil {
		s.responseErrorHandlerFunc(ctx, req, err)
		return
	}

	resp, err := proto.Marshal(res)
	if err != nil {
		s.responseErrorHandlerFunc(ctx, req, err)
		return
	}

	req.Respond(resp)
}

func (s *DocumentServiceServer) handleUploadDocument(
	ctx context.Context,
	req micro.Request,
) {
	msg := new(UploadDocumentRequest)

	if err := proto.Unmarshal(req.Data(), msg); err != nil {
		s.requestErrorHandlerFunc(ctx, req, err)
		return
	}

	res, err := s.impl.UploadDocument(ctx, msg)
	if err != nil {
		s.responseErrorHandlerFunc(ctx, req, err)
		return
	}

	resp, err := proto.Marshal(res)
	if err != nil {
		s.responseErrorHandlerFunc(ctx, req, err)
		return
	}

	req.Respond(resp)
}

func (s *DocumentServiceServer) handleSearchAttribute(
	ctx context.Context,
	req micro.Request,
) {
	msg := new(SearchAttributeRequest)

	if err := proto.Unmarshal(req.Data(), msg); err != nil {
		s.requestErrorHandlerFunc(ctx, req, err)
		return
	}

	res, err := s.impl.SearchAttribute(ctx, msg)
	if err != nil {
		s.responseErrorHandlerFunc(ctx, req, err)
		return
	}

	resp, err := proto.Marshal(res)
	if err != nil {
		s.responseErrorHandlerFunc(ctx, req, err)
		return
	}

	req.Respond(resp)
}

func (s *DocumentServiceServer) handleSearchDocument(
	ctx context.Context,
	req micro.Request,
) {
	msg := new(SearchRequest)

	if err := proto.Unmarshal(req.Data(), msg); err != nil {
		s.requestErrorHandlerFunc(ctx, req, err)
		return
	}

	res, err := s.impl.SearchDocument(ctx, msg)
	if err != nil {
		s.responseErrorHandlerFunc(ctx, req, err)
		return
	}

	resp, err := proto.Marshal(res)
	if err != nil {
		s.responseErrorHandlerFunc(ctx, req, err)
		return
	}

	req.Respond(resp)
}

type ClientError struct {
	Status      string
	Description string
}

func (ce *ClientError) Error() string {
	return fmt.Sprintf("[%s] %s", ce.Status, ce.Description)
}

type Invoker func(ctx context.Context, msg *nats.Msg) (*nats.Msg, error)

type InvokerMiddleware func(Invoker) Invoker

type DocumentServiceClientConfig struct {
	MicroSubject string
	Middleware   InvokerMiddleware
}

// DocumentServiceClient is the client API for DocumentService service.
type DocumentServiceClient struct {
	nc  *nats.Conn
	cfg DocumentServiceClientConfig
}

// NewDocumentServiceClient creates a new NATS microservice client.
func NewDocumentServiceClient(
	cfg DocumentServiceClientConfig,
	nc *nats.Conn,
) *DocumentServiceClient {
	return &DocumentServiceClient{nc: nc, cfg: cfg}
}

func (c *DocumentServiceClient) FetchDocument(
	ctx context.Context,
	req *FetchDocumentRequest,
) (*FetchDocumentResponse, error) {
	resp := new(FetchDocumentResponse)

	data, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	reqmsg := nats.NewMsg(c.cfg.MicroSubject + ".fetch_document")
	reqmsg.Data = data

	doRequest := c.nc.RequestMsgWithContext

	if c.cfg.Middleware != nil {
		doRequest = c.cfg.Middleware(doRequest)
	}

	respmsg, err := doRequest(ctx, reqmsg)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if respmsg.Header.Get(micro.ErrorHeader) != "" {
		return nil, &ClientError{
			Status:      respmsg.Header.Get(micro.ErrorCodeHeader),
			Description: respmsg.Header.Get(micro.ErrorHeader),
		}
	}

	if err := proto.Unmarshal(respmsg.Data, resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp, nil
}

func (c *DocumentServiceClient) UploadDocument(
	ctx context.Context,
	req *UploadDocumentRequest,
) (*UploadDocumentResponse, error) {
	resp := new(UploadDocumentResponse)

	data, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	reqmsg := nats.NewMsg(c.cfg.MicroSubject + ".upload_document")
	reqmsg.Data = data

	doRequest := c.nc.RequestMsgWithContext

	if c.cfg.Middleware != nil {
		doRequest = c.cfg.Middleware(doRequest)
	}

	respmsg, err := doRequest(ctx, reqmsg)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if respmsg.Header.Get(micro.ErrorHeader) != "" {
		return nil, &ClientError{
			Status:      respmsg.Header.Get(micro.ErrorCodeHeader),
			Description: respmsg.Header.Get(micro.ErrorHeader),
		}
	}

	if err := proto.Unmarshal(respmsg.Data, resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp, nil
}

func (c *DocumentServiceClient) SearchAttribute(
	ctx context.Context,
	req *SearchAttributeRequest,
) (*SearchAttributeResponse, error) {
	resp := new(SearchAttributeResponse)

	data, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	reqmsg := nats.NewMsg(c.cfg.MicroSubject + ".search_attribute")
	reqmsg.Data = data

	doRequest := c.nc.RequestMsgWithContext

	if c.cfg.Middleware != nil {
		doRequest = c.cfg.Middleware(doRequest)
	}

	respmsg, err := doRequest(ctx, reqmsg)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if respmsg.Header.Get(micro.ErrorHeader) != "" {
		return nil, &ClientError{
			Status:      respmsg.Header.Get(micro.ErrorCodeHeader),
			Description: respmsg.Header.Get(micro.ErrorHeader),
		}
	}

	if err := proto.Unmarshal(respmsg.Data, resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp, nil
}

func (c *DocumentServiceClient) SearchDocument(
	ctx context.Context,
	req *SearchRequest,
) (*SearchDocumentResponse, error) {
	resp := new(SearchDocumentResponse)

	data, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	reqmsg := nats.NewMsg(c.cfg.MicroSubject + ".search_document")
	reqmsg.Data = data

	doRequest := c.nc.RequestMsgWithContext

	if c.cfg.Middleware != nil {
		doRequest = c.cfg.Middleware(doRequest)
	}

	respmsg, err := doRequest(ctx, reqmsg)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if respmsg.Header.Get(micro.ErrorHeader) != "" {
		return nil, &ClientError{
			Status:      respmsg.Header.Get(micro.ErrorCodeHeader),
			Description: respmsg.Header.Get(micro.ErrorHeader),
		}
	}

	if err := proto.Unmarshal(respmsg.Data, resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp, nil
}
