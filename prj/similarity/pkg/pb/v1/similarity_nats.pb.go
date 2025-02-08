// Code generated by protoc-gen-natsmicro. DO NOT EDIT.

package similarity

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

// SimilarityServiceServerInterface is the server API for SimilarityService service.
type SimilarityServiceServerInterface interface {
	SearchSimilar(context.Context, *DocumentId) (*SearchSimilarResponse, error)
	SearchSimilarityHistory(context.Context, *SearchSimilarityHistoryRequest) (*SearchSimilarityHistoryResponse, error)
	mustEmbedUnimplementedGreeterServer()
}

type UnimplementedSimilarityServiceServer struct{}

func (UnimplementedSimilarityServiceServer) SearchSimilar(context.Context, *DocumentId) (*SearchSimilarResponse, error) {
	return nil, errors.New("method SearchSimilar not implemented")
}
func (UnimplementedSimilarityServiceServer) SearchSimilarityHistory(context.Context, *SearchSimilarityHistoryRequest) (*SearchSimilarityHistoryResponse, error) {
	return nil, errors.New("method SearchSimilarityHistory not implemented")
}

func (UnimplementedSimilarityServiceServer) mustEmbedUnimplementedSimilarityServiceServer() {}

func (UnimplementedSimilarityServiceServer) testEmbeddedByValue() {}

type UnsafeSimilarityServiceServer interface {
	mustEmbedUnimplementedGreeterServer()
}

type SimilarityServiceServerConfig struct {
	micro.Config
	RequestTimeout       time.Duration
	Middleware           Middleware
	RequestErrorHandler  ErrorHandler
	ResponseErrorHandler ErrorHandler
}

type SimilarityServiceServer struct {
	nc   *nats.Conn
	impl SimilarityServiceServerInterface
	cfg  SimilarityServiceServerConfig
	done chan struct{}

	requestErrorHandlerFunc  ErrorHandler
	responseErrorHandlerFunc ErrorHandler
}

// NewSimilarityServiceServer  creates a new NATS microservice server.
func NewSimilarityServiceServer(
	cfg SimilarityServiceServerConfig,
	nc *nats.Conn,
	impl SimilarityServiceServerInterface,
) *SimilarityServiceServer {
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

	return &SimilarityServiceServer{
		nc:                       nc,
		impl:                     impl,
		cfg:                      cfg,
		done:                     make(chan struct{}),
		requestErrorHandlerFunc:  cfg.RequestErrorHandler,
		responseErrorHandlerFunc: cfg.ResponseErrorHandler,
	}
}

// Start starts the microservice and blocking until application context done.
func (s *SimilarityServiceServer) Start(ctx context.Context) error {
	srv, err := micro.AddService(s.nc, s.cfg.Config)
	if err != nil {
		return fmt.Errorf("failed to start microservice: %w", err)
	}

	defer srv.Stop()

	group := srv.AddGroup(s.cfg.Config.Name)

	// Register handlers
	if err := group.AddEndpoint("search_similar", s.toMicroHandler(s.handleSearchSimilar)); err != nil {
		return fmt.Errorf("failed to add endpoint SearchSimilar: %w", err)
	}
	if err := group.AddEndpoint("search_similarity_history", s.toMicroHandler(s.handleSearchSimilarityHistory)); err != nil {
		return fmt.Errorf("failed to add endpoint SearchSimilarityHistory: %w", err)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.done:
		return nil
	}
}

func (s *SimilarityServiceServer) Stop() error {
	s.done <- struct{}{}
	return nil
}

func (s *SimilarityServiceServer) toMicroHandler(h Handler) micro.HandlerFunc {
	return func(req micro.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.RequestTimeout)
		defer cancel()

		if s.cfg.Middleware != nil {
			h = s.cfg.Middleware(h)
		}

		h(ctx, req)
	}
}

func (s *SimilarityServiceServer) handleSearchSimilar(
	ctx context.Context,
	req micro.Request,
) {
	msg := new(DocumentId)

	if err := proto.Unmarshal(req.Data(), msg); err != nil {
		s.requestErrorHandlerFunc(ctx, req, err)
		return
	}

	res, err := s.impl.SearchSimilar(ctx, msg)
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

func (s *SimilarityServiceServer) handleSearchSimilarityHistory(
	ctx context.Context,
	req micro.Request,
) {
	msg := new(SearchSimilarityHistoryRequest)

	if err := proto.Unmarshal(req.Data(), msg); err != nil {
		s.requestErrorHandlerFunc(ctx, req, err)
		return
	}

	res, err := s.impl.SearchSimilarityHistory(ctx, msg)
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

type SimilarityServiceClientConfig struct {
	MicroSubject string
	Middleware   InvokerMiddleware
}

// SimilarityServiceClient is the client API for SimilarityService service.
type SimilarityServiceClient struct {
	nc  *nats.Conn
	cfg SimilarityServiceClientConfig
}

// NewSimilarityServiceClient creates a new NATS microservice client.
func NewSimilarityServiceClient(
	cfg SimilarityServiceClientConfig,
	nc *nats.Conn,
) *SimilarityServiceClient {
	return &SimilarityServiceClient{nc: nc, cfg: cfg}
}

func (c *SimilarityServiceClient) SearchSimilar(
	ctx context.Context,
	req *DocumentId,
) (*SearchSimilarResponse, error) {
	resp := new(SearchSimilarResponse)

	data, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	reqmsg := nats.NewMsg(c.cfg.MicroSubject + ".search_similar")
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

func (c *SimilarityServiceClient) SearchSimilarityHistory(
	ctx context.Context,
	req *SearchSimilarityHistoryRequest,
) (*SearchSimilarityHistoryResponse, error) {
	resp := new(SearchSimilarityHistoryResponse)

	data, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	reqmsg := nats.NewMsg(c.cfg.MicroSubject + ".search_similarity_history")
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
