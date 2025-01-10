package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	openapi "simrep/api/rest/gen"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gorilla/mux"
	openapimw "github.com/oapi-codegen/nethttp-middleware"
	"github.com/rs/cors"
	"golang.org/x/sync/errgroup"
)

const (
	docMime  = "application/msword"
	docxMime = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	pdfMime  = "application/pdf"
)

type Server struct {
	addr   string
	router http.Handler
}

type Opts struct {
	Port int

	Spec            []byte
	DocumentHandler DocumentHandler
	AnalyzeHandler  AnalyzeHandler
}

func New(
	opts Opts,
) (*Server, error) {
	spec, err := openapi3.NewLoader().LoadFromData(opts.Spec)
	if err != nil {
		return nil, fmt.Errorf("load spec: %w", err)
	}

	router := mux.NewRouter()

	router.Use(openapimw.OapiRequestValidator(spec))
	openapi3filter.RegisterBodyDecoder(docMime, openapi3filter.FileBodyDecoder)
	openapi3filter.RegisterBodyDecoder(docxMime, openapi3filter.FileBodyDecoder)
	openapi3filter.RegisterBodyDecoder(pdfMime, openapi3filter.FileBodyDecoder)

	compose := struct {
		DocumentHandler
		AnalyzeHandler
	}{
		DocumentHandler: opts.DocumentHandler,
		AnalyzeHandler:  opts.AnalyzeHandler,
	}

	stricthandler := openapi.NewStrictHandlerWithOptions(
		compose,
		nil,
		openapi.StrictHTTPServerOptions{
			RequestErrorHandlerFunc:  nil,
			ResponseErrorHandlerFunc: responseErrorHandler,
		},
	)

	router.PathPrefix("").Handler(openapi.Handler(stricthandler))

	handler := cors.New(cors.Options{ //nolint:exhaustruct
		AllowedOrigins: []string{"*"},
	}).Handler(router)

	return &Server{
		addr:   fmt.Sprintf(":%d", opts.Port),
		router: handler,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	slog.Info("starting server", "addr", s.addr)

	httpServer := http.Server{ //nolint:exhaustruct
		Addr:              s.addr,
		Handler:           s.router,
		ReadTimeout:       60 * time.Second, //nolint:gomnd,mnd
		WriteTimeout:      60 * time.Second, //nolint:gomnd,mnd
		IdleTimeout:       60 * time.Second, //nolint:gomnd,mnd
		ReadHeaderTimeout: 2 * time.Second,  //nolint:gomnd,mnd
	}

	errg, ctx := errgroup.WithContext(ctx)

	errg.Go(httpServer.ListenAndServe)

	errg.Go(func() error {
		<-ctx.Done()
		return httpServer.Shutdown(ctx)
	})

	if err := errg.Wait(); err != nil {
		return fmt.Errorf("exit: %w", err)
	}

	return nil
}
