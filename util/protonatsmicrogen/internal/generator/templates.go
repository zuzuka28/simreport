package generator

const serviceTmpl = `
// Code generated by protoc-gen-nats. DO NOT EDIT.

package {{ .Package }}

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

{{ range .Services }}
// {{ .Name }}Server is the server API for {{ .Name }} service.
type {{ .Name }}Server interface {
	{{- range .Methods }}
	{{ .Name }}(context.Context, *{{ .InputType }}) (*{{ .OutputType }}, error)
	{{- end }}
}

type {{ lower .Name }}Server struct {
	srv  micro.Service
	impl {{ .Name }}Server
}

// New{{ .Name }}Server creates a new NATS microservice server
func New{{ .Name }}Server(nc *nats.Conn, impl {{ .Name }}Server) (*{{ lower .Name }}Server, error) {
	srv, err := micro.AddService(nc, micro.Config{
		Name: "{{ .Name }}",
	})
	if err != nil {
		return nil, err
	}

	s := &{{ lower .Name }}Server{
		srv:  srv,
		impl: impl,
	}

	// Register handlers
	{{- range .Methods }}
	if err := srv.AddEndpoint("{{ .Name }}", micro.HandlerFunc(s.handle{{ .Name }})); err != nil {
		return nil, err
	}
	{{- end }}

	return s, nil
}

{{ range .Methods }}
func (s *{{ lower .Name }}Server) handle{{ .Name }}(ctx context.Context, req any) (any, error) {
	msg, ok := req.(*{{ .InputType }})
	if !ok {
		return nil, fmt.Errorf("invalid request type: %T", req)
	}

	return s.impl.{{ .Name }}(ctx, msg)
}
{{ end }}

// {{ .Name }}Client is the client API for {{ .Name }} service.
type {{ .Name }}Client struct {
	nc *nats.Conn
}

// New{{ .Name }}Client creates a new NATS microservice client
func New{{ .Name }}Client(nc *nats.Conn) *{{ .Name }}Client {
	return &{{ .Name }}Client{nc: nc}
}

{{ range .Methods }}
func (c *{{ $.Name }}Client) {{ .Name }}(ctx context.Context, req *{{ .InputType }}) (*{{ .OutputType }}, error) {
	resp := new({{ .OutputType }})

	msg, err := c.nc.RequestWithContext(ctx, "{{ $.Name }}.{{ .Name }}", req)
	if err != nil {
		return nil, err
	}

	if err := resp.UnmarshalVT(msg.Data); err != nil {
		return nil, err
	}

	return resp, nil
}
{{ end }}
{{ end }}
`
