package generator

const serviceTmpl = `
// Code generated by protoc-gen-natsmicro. DO NOT EDIT.

package {{ .Package }}

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"google.golang.org/protobuf/proto"
)

type Handler func(ctx context.Context, req micro.Request)

type Middleware func(Handler) Handler

{{ range .Services }}
// {{ .Name }}Server is the server API for {{ .Name }} service.
type {{ .Name }}Server interface {
    {{- range .Methods }}
    {{ .Name }}(context.Context, *{{ .InputType }}) (*{{ .OutputType }}, error)
    {{- end }}
}

type {{ .Name }}ServerConfig struct {
	micro.Config
	RequestTimeout  time.Duration
	Middleware      Middleware
	OnError         func(ctx context.Context, err error)
}

type {{ lower .Name }}Server struct {
    srv    micro.Service
    impl   {{ .Name }}Server
    cfg    {{ .Name }}ServerConfig
}

// New{{ .Name }}Server creates a new NATS microservice server.
func New{{ .Name }}Server(
    cfg {{ .Name }}ServerConfig,
    nc *nats.Conn,
    impl {{ .Name }}Server,
) (*{{ lower .Name }}Server, error) {
    srv, err := micro.AddService(nc, cfg.Config)
    if err != nil {
        return nil, fmt.Errorf("failed to create microservice: %w", err)
    }

    if cfg.RequestTimeout == 0 {
        cfg.RequestTimeout = time.Second * 60
    }

    s := &{{ lower .Name }}Server{
        srv:    srv,
        impl:   impl,
        cfg: cfg,
    }

    // Register handlers
    {{- range .Methods }}
    if err := srv.AddEndpoint(cfg.Name+".{{ lower .Name }}", s.toMicroHandler(s.handle{{ .Name }})); err != nil {
        return nil, fmt.Errorf("failed to add endpoint {{ .Name }}: %w", err)
    }
    {{- end }}

    return s, nil
}

// Stop stops the microservice.
func (s *{{ lower .Name }}Server) Stop() error {
    return s.srv.Stop()
}

func (s *{{ lower .Name }}Server) toMicroHandler(h Handler) micro.HandlerFunc {
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
func (s *{{ lower .Reciever }}Server) handle{{ .Name }}(
    ctx context.Context,
    req micro.Request,
) {
	msg := new({{ .InputType }})

	if err := proto.Unmarshal(req.Data(), msg); err != nil {
		req.Error("500", "unproccessable request", nil, nil)
		return
	}

	res, err := s.impl.{{ .Name }}(ctx, msg)
	if err != nil {
		req.Error("500", "server error", nil, nil)
		return
	}

	resp, err := proto.Marshal(res)
	if err != nil {
	    if s.cfg.OnError != nil {
	        s.cfg.OnError(ctx, err)
	    }

		req.Error("500", "server error", nil, nil)
		return
	}

	req.Respond(resp)
}
{{ end }}

// {{ .Name }}Client is the client API for {{ .Name }} service.
type {{ .Name }}Client struct {
    nc     *nats.Conn
}

// New{{ .Name }}Client creates a new NATS microservice client.
func New{{ .Name }}Client(nc *nats.Conn) *{{ .Name }}Client {
    return &{{ .Name }}Client{nc: nc}
}

{{ range .Methods }}
func (c *{{ .Reciever }}Client) {{ .Name }}(
    ctx context.Context,
    req *{{ .InputType }},
) (*{{ .OutputType }}, error) {
    resp := new({{ .OutputType }})

    data, err := proto.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }

    msg, err := c.nc.RequestWithContext(ctx, "{{ lower .Reciever }}.{{ lower .Name }}", data)
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
