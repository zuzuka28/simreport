package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	openapi "simrep/api/rest/gen"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	openapimw "github.com/oapi-codegen/nethttp-middleware"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	addr   string
	router *mux.Router
}

type Opts struct {
	Addr string

	Spec              []byte
	FileHandler       FileHandler
	SimilarityHandler SimilarityHandler
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

	compose := struct {
		FileHandler
		SimilarityHandler
	}{
		FileHandler:       opts.FileHandler,
		SimilarityHandler: opts.SimilarityHandler,
	}

	stricthandler := openapi.NewStrictHandler(compose, nil)
	router.Handle("", openapi.Handler(stricthandler))

	return &Server{
		addr:   opts.Addr,
		router: router,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	slog.Info("starting server", "addr", s.addr)

	httpServer := http.Server{ //nolint:exhaustruct
		Addr:              s.addr,
		Handler:           s.router,
		ReadTimeout:       10 * time.Second, //nolint:gomnd,mnd
		WriteTimeout:      10 * time.Second, //nolint:gomnd,mnd
		IdleTimeout:       30 * time.Second, //nolint:gomnd,mnd
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