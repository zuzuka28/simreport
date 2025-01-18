package minhashlsh_test

import (
	"context"
	"hash/fnv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zuzuka28/simreport/lib/minhash"
	"github.com/zuzuka28/simreport/lib/minhashlsh"
	"github.com/zuzuka28/simreport/lib/minhashlsh/inmemorystorage"
)

func TestMinhashLSH(t *testing.T) {
	var (
		permutationsNum = 512
		hasher          = func(b []byte) uint64 {
			h := fnv.New64a()
			h.Write(b)
			return h.Sum64()
		}
		seed = 1
	)

	ctx := context.Background()

	lsh := minhashlsh.New(inmemorystorage.New(), "", permutationsNum, 3, hasher)

	mh1 := minhash.New(permutationsNum, hasher, seed)
	mh2 := minhash.New(permutationsNum, hasher, seed)

	mh1.Push([]byte("value1"))
	mh1.Push([]byte("value2"))
	mh2.Push([]byte("value3"))
	mh2.Push([]byte("value4"))

	sim, _ := mh1.Similarity(mh2)
	t.Log("similarity mh1 on mh2", sim)

	err := lsh.Insert(ctx, "key1", mh1)
	require.NoError(t, err)

	err = lsh.Insert(ctx, "key2", mh2)
	require.NoError(t, err)

	results, err := lsh.Query(ctx, mh1)
	require.NoError(t, err)
	require.Contains(t, results, "key1")
	require.NotContains(t, results, "key2")

	mh3 := minhash.New(permutationsNum, hasher, seed)
	mh3.Push([]byte("value5"))

	results, err = lsh.Query(ctx, mh3)
	require.NoError(t, err)
	require.Empty(t, results)
}
