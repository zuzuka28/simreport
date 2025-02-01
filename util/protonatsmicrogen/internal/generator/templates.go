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

type ErrorHandler func(ctx context.Context, msg micro.Request, err error)

type Middleware func(Handler) Handler

{{ range .Services }}
// {{ .Name }}Server is the server API for {{ .Name }} service.
type {{ .Name }}Server interface {
    {{- range .Methods }}
    {{ .Name }}(context.Context, *{{ .InputType }}) (*{{ .OutputType }}, error)
    {{- end }}
}

type {{ .Name }}NatsServerConfig struct {
	micro.Config
	RequestTimeout  time.Duration
	Middleware      Middleware
	OnError         ErrorHandler
}

type {{ .Name }}NatsServer struct {
    srv    micro.Service
    impl   {{ .Name }}Server
    cfg    {{ .Name }}NatsServerConfig
}

// New{{ .Name }}NatsServer  creates a new NATS microservice server.
func New{{ .Name }}NatsServer(
    cfg {{ .Name }}NatsServerConfig,
    nc *nats.Conn,
    impl {{ .Name }}Server,
) (*{{ .Name }}NatsServer, error) {
    srv, err := micro.AddService(nc, cfg.Config)
    if err != nil {
        return nil, fmt.Errorf("failed to create microservice: %w", err)
    }

    if cfg.RequestTimeout == 0 {
        cfg.RequestTimeout = time.Second * 60
    }

    s := &{{ .Name }}NatsServer{
        srv:    srv,
        impl:   impl,
        cfg:    cfg,
    }

    group := srv.AddGroup(cfg.Name)

    // Register handlers
    {{- range .Methods }}
    if err := group.AddEndpoint("{{ snakecase .Name }}", s.toMicroHandler(s.handle{{ .Name }})); err != nil {
        return nil, fmt.Errorf("failed to add endpoint {{ .Name }}: %w", err)
    }
    {{- end }}

    return s, nil
}

// Stop stops the microservice.
func (s *{{ .Name }}NatsServer) Stop() error {
    return s.srv.Stop()
}

func (s *{{ .Name }}NatsServer) toMicroHandler(h Handler) micro.HandlerFunc {
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
func (s *{{ .Reciever }}NatsServer) handle{{ .Name }}(
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
	    if s.cfg.OnError != nil {
	        s.cfg.OnError(ctx, req, err)
	    } else {
		    req.Error("500", "server error", nil, nil)
	    }

		return
	}

	resp, err := proto.Marshal(res)
	if err != nil {
		req.Error("500", "server error", nil, nil)
		return
	}

	req.Respond(resp)
}
{{ end }}

type {{ .Name }}NatsClientConfig struct {
	ServerName string
}

// {{ .Name }}Client is the client API for {{ .Name }} service.
type {{ .Name }}NatsClient struct {
    nc     *nats.Conn
    cfg    {{ .Name }}NatsClientConfig
}

// New{{ .Name }}Client creates a new NATS microservice client.
func New{{ .Name }}Client(
    cfg {{ .Name }}NatsClientConfig,
    nc *nats.Conn,
) *{{ .Name }}NatsClient {
    return &{{ .Name }}NatsClient{nc: nc, cfg: cfg}
}

{{ range .Methods }}
func (c *{{ .Reciever }}NatsClient) {{ .Name }}(
    ctx context.Context,
    req *{{ .InputType }},
) (*{{ .OutputType }}, error) {
    resp := new({{ .OutputType }})

    data, err := proto.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }

    msg, err := c.nc.RequestWithContext(ctx, c.cfg.ServerName+".{{ snakecase .Name }}", data)
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
