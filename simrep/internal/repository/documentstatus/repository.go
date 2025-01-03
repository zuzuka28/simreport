package documentstatus

import (
	"context"
	"errors"
	"fmt"
	"simrep/internal/model"
	"sync"

	"github.com/nats-io/nats.go/jetstream"
	"golang.org/x/sync/errgroup"
)

const subject = "documentstatus"

type Repository struct {
	kv jetstream.KeyValue
	p  jetstream.Publisher
}

func NewRepository(
	kv jetstream.KeyValue,
	p jetstream.Publisher,
) *Repository {
	return &Repository{
		kv: kv,
		p:  p,
	}
}

func (r *Repository) Update(
	ctx context.Context,
	cmd model.DocumentStatusUpdateCommand,
) error {
	key := cmd.ID

	fullsubject := subject + "." + string(cmd.Status)

	_, err := r.p.Publish(ctx, fullsubject, []byte(cmd.ID))
	if err != nil {
		return fmt.Errorf("publush update: %w", err)
	}

	_, err = r.kv.Put(ctx, key, []byte(cmd.Status))
	if err != nil {
		return fmt.Errorf("put status: %w", err)
	}

	return nil
}

func (r *Repository) Fetch(
	ctx context.Context,
	query model.DocumentStatusQuery,
) ([]*model.DocumentStatus, error) {
	if len(query.IDs) > 0 {
		return r.fetchByIDs(ctx, query.IDs)
	}

	ids, err := r.topKeys(ctx, query.Limit)
	if err != nil {
		return nil, fmt.Errorf("fetch top keys: %w", err)
	}

	return r.fetchByIDs(ctx, ids)
}

func (r *Repository) topKeys(
	ctx context.Context,
	limit int,
) ([]string, error) {
	w, err := r.kv.ListKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("list keys: %w", err)
	}

	defer func() { _ = w.Stop() }()

	var keys []string

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context canceled: %w", ctx.Err())

		case key, ok := <-w.Keys():
			if !ok {
				return keys, nil
			}

			keys = append(keys, key)
		}

		if limit != 0 && len(keys) >= limit {
			break
		}
	}

	return keys, nil
}

func (r *Repository) fetchByIDs(
	ctx context.Context,
	ids []string,
) ([]*model.DocumentStatus, error) {
	eg, egCtx := errgroup.WithContext(ctx)

	var (
		res   []*model.DocumentStatus
		resMu sync.Mutex
	)

	for _, id := range ids {
		eg.Go(func() error {
			el, err := r.kv.Get(egCtx, id)
			if err != nil && !errors.Is(err, jetstream.ErrNoKeysFound) {
				return fmt.Errorf("get key: %w", err)
			}

			resMu.Lock()
			defer resMu.Unlock()

			if el == nil {
				res = append(res, &model.DocumentStatus{
					ID:     id,
					Status: model.DocumentProcessingStatusNotFound,
				})

				return nil
			}

			res = append(res, &model.DocumentStatus{
				ID:     id,
				Status: model.DocumentProcessingStatus(el.Value()),
			})

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("fetch keys: %w", err)
	}

	return res, nil
}