package shingleindex

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/shingleindex/internal/model"

	"github.com/zuzuka28/simreport/lib/minhash"
)

type saveService struct {
	r Repository
}

func newSaveService(
	r Repository,
) *saveService {
	return &saveService{
		r: r,
	}
}

func (s *saveService) Save(
	ctx context.Context,
	cmd model.DocumentSaveCommand,
) error {
	ncmd := mapDocumentToMinhashSaveCommand(cmd)

	if err := s.r.Save(ctx, ncmd); err != nil {
		return fmt.Errorf("save analyzed document: %w", err)
	}

	return nil
}

func mapDocumentToMinhashSaveCommand(
	in model.DocumentSaveCommand,
) model.MinhashSaveCommand {
	text := string(in.Item.Text)
	text = normalize(text)

	shingles := shingle(text, shingleSize)

	mh := minhash.New(
		permutations,
		hasher,
		seed,
	)

	for shingle := range shingles {
		mh.Push([]byte(shingle))
	}

	return model.MinhashSaveCommand{
		DocumentID: in.Item.ID,
		Minhash:    mh,
	}
}
