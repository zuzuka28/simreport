package generator

const serviceTmpl = `
// Code generated by protoc-gen-natsmicro. DO NOT EDIT.

package {{ .GoPackageName }}

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

{{ range .Services }}
// {{ .GoName }}ServerInterface is the server API for {{ .GoName }} service.
type {{ .GoName }}ServerInterface interface {
    {{- range .Methods }}
    {{ .GoName }}(context.Context, *{{ .Input.GoIdent.GoName }}) (*{{ .Output.GoIdent.GoName }}, error)
    {{- end }}
    mustEmbedUnimplementedGreeterServer()
}

type Unimplemented{{ .GoName }}Server struct{}

{{- range .Methods }}
func (Unimplemented{{ .Parent.GoName }}Server) {{ .GoName }}(context.Context, *{{ .Input.GoIdent.GoName }}) (*{{ .Output.GoIdent.GoName }}, error) {
	return nil, errors.New("method {{ .GoName }} not implemented")
}
{{- end }}

func (Unimplemented{{ .GoName }}Server) mustEmbedUnimplemented{{ .GoName }}Server() {}

func (Unimplemented{{ .GoName }}Server) testEmbeddedByValue() {}

type Unsafe{{ .GoName }}Server interface {
	mustEmbedUnimplementedGreeterServer()
}

type {{ .GoName }}ServerConfig struct {
	micro.Config
	RequestTimeout              time.Duration
	Middleware                  Middleware
	RequestErrorHandler         ErrorHandler
	ResponseErrorHandler        ErrorHandler
}

type {{ .GoName }}Server struct {
    nc     *nats.Conn
    impl   {{ .GoName }}ServerInterface
    cfg    {{ .GoName }}ServerConfig
    done   chan struct{}

    requestErrorHandlerFunc     ErrorHandler
    responseErrorHandlerFunc    ErrorHandler
}

// New{{ .GoName }}Server  creates a new NATS microservice server.
func New{{ .GoName }}Server(
    cfg {{ .GoName }}ServerConfig,
    nc *nats.Conn,
    impl {{ .GoName }}ServerInterface,
) (*{{ .GoName }}Server) {
	if t, ok := impl.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}

    if cfg.RequestTimeout == 0 {
        cfg.RequestTimeout = time.Second * 60
    }

    if cfg.RequestErrorHandler == nil {
        cfg.RequestErrorHandler = func(_ context.Context, req micro.Request, _ error) {
            req.Error("500", "unproccessable request", nil, nil)
        }
    }

    if cfg.ResponseErrorHandler == nil {
        cfg.ResponseErrorHandler = func(_ context.Context, req micro.Request, _ error) {
            req.Error("500", "internal server error", nil, nil)
        }
    }

    return &{{ .GoName }}Server{
        nc: nc,
        impl:   impl,
        cfg:    cfg,
        done:   make(chan struct{}),
        requestErrorHandlerFunc: cfg.RequestErrorHandler,
        responseErrorHandlerFunc: cfg.ResponseErrorHandler,
    }
}

// Start starts the microservice and blocking until application context done.
func (s *{{ .GoName }}Server) Start(ctx context.Context) error {
    srv, err := micro.AddService(s.nc, s.cfg.Config)
    if err != nil {
        return fmt.Errorf("failed to start microservice: %w", err)
    }

    defer srv.Stop()

    group := srv.AddGroup(s.cfg.Config.Name)

    // Register handlers
    {{- range .Methods }}
    if err := group.AddEndpoint("{{ snakecase .GoName }}", s.toMicroHandler(s.handle{{ .GoName }})); err != nil {
        return fmt.Errorf("failed to add endpoint {{ .GoName }}: %w", err)
    }
    {{- end }}

    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-s.done:
        return nil
    }
}

func (s *{{ .GoName }}Server) Stop() error {
    s.done<-struct{}{}
    return nil
}

func (s *{{ .GoName }}Server) toMicroHandler(h Handler) micro.HandlerFunc {
	return func(req micro.Request) {
	    ctx, cancel := context.WithTimeout(context.Background(), s.cfg.RequestTimeout)
	    defer cancel()

	    if s.cfg.Middleware != nil {
            h = s.cfg.Middleware(h)
	    }

	    h(ctx, req)
	}
}

{{ range .Methods }}
func (s *{{ .Parent.GoName }}Server) handle{{ .GoName }}(
    ctx context.Context,
    req micro.Request,
) {
	msg := new({{ .Input.GoIdent.GoName }})

	if err := proto.Unmarshal(req.Data(), msg); err != nil {
		s.requestErrorHandlerFunc(ctx, req, err)
		return
	}

	res, err := s.impl.{{ .GoName }}(ctx, msg)
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
{{ end }}

type {{ .GoName }}ClientConfig struct {
	MicroSubject string
}

// {{ .GoName }}Client is the client API for {{ .GoName }} service.
type {{ .GoName }}Client struct {
    nc     *nats.Conn
    cfg    {{ .GoName }}ClientConfig
}

// New{{ .GoName }}Client creates a new NATS microservice client.
func New{{ .GoName }}Client(
    cfg {{ .GoName }}ClientConfig,
    nc *nats.Conn,
) *{{ .GoName }}Client {
    return &{{ .GoName }}Client{nc: nc, cfg: cfg}
}

{{ range .Methods }}
func (c *{{ .Parent.GoName }}Client) {{ .GoName }}(
    ctx context.Context,
    req *{{ .Input.GoIdent.GoName }},
) (*{{ .Output.GoIdent.GoName }}, error) {
    resp := new({{ .Output.GoIdent.GoName }})

    data, err := proto.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }

    msg, err := c.nc.RequestWithContext(ctx, c.cfg.MicroSubject+".{{ snakecase .GoName }}", data)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }

    if err := proto.Unmarshal(msg.Data, resp); err != nil {
        return nil, fmt.Errorf("failed to unmarshal response: %w", err)
    }

    return resp, nil
}
{{ end }}
{{ end }}
`
