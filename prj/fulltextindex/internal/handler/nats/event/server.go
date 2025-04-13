package server

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const indexIngestionSubject = "fulltextindex_ingestion"

type Server struct {
	conn *nats.Conn
	dh   IndexerHandler
}

func NewServer(
	conn *nats.Conn,
	dh IndexerHandler,
) *Server {
	return &Server{
		conn: conn,
		dh:   dh,
	}
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.initWorkQueue(ctx); err != nil {
		return fmt.Errorf("init work queue: %w", err)
	}

	js, err := s.conn.JetStream()
	if err != nil {
		return fmt.Errorf("create jetstream: %w", err)
	}

	sub, err := js.QueueSubscribe(indexIngestionSubject, "workers", s.dh.Save)
	if err != nil {
		return fmt.Errorf("create indexer sub: %w", err)
	}

	defer func() { _ = sub.Drain() }()

	<-ctx.Done()

	return nil
}

func (s *Server) initWorkQueue(ctx context.Context) error {
	js, err := jetstream.New(s.conn)
	if err != nil {
		return fmt.Errorf("create jetstream: %w", err)
	}

	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{ //nolint:exhaustruct
		Name: indexIngestionSubject,
		Sources: []*jetstream.StreamSource{
			{
				Name:              "documentstatus",
				OptStartSeq:       0,
				OptStartTime:      nil,
				FilterSubject:     "documentstatus.document_parsed",
				SubjectTransforms: nil,
				External:          nil,
				Domain:            "",
			},
		},
		SubjectTransform: &jetstream.SubjectTransformConfig{
			Source:      "documentstatus.document_parsed",
			Destination: indexIngestionSubject,
		},
		Retention: jetstream.WorkQueuePolicy,
	})
	if err != nil {
		return fmt.Errorf("create mirroring stream: %w", err)
	}

	return nil
}
