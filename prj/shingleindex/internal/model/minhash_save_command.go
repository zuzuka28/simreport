package model

import "github.com/zuzuka28/simreport/lib/minhash"

type MinhashSaveCommand struct {
	DocumentID string
	Minhash    *minhash.MinHash
}
