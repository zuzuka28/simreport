//nolint:revive,gomnd,mnd
package main

import (
	"sort"
)

type Match struct {
	A    int
	B    int
	Size int
}

// SequenceMatcher compares sequence of strings. The basic
// algorithm predates, and is a little fancier than, an algorithm
// published in the late 1980's by Ratcliff and Obershelp under the
// hyperbolic name "gestalt pattern matching".  The basic idea is to find
// the longest contiguous matching subsequence that contains no "junk"
// elements (R-O doesn't address junk).  The same idea is then applied
// recursively to the pieces of the sequences to the left and to the right
// of the matching subsequence.  This does not yield minimal edit
// sequences, but does tend to yield matches that "look right" to people.
//
// SequenceMatcher tries to compute a "human-friendly diff" between two
// sequences.  Unlike e.g. UNIX(tm) diff, the fundamental notion is the
// longest *contiguous* & junk-free matching subsequence.  That's what
// catches peoples' eyes.  The Windows(tm) windiff has another interesting
// notion, pairing up elements that appear uniquely in each sequence.
// That, and the method here, appear to yield more intuitive difference
// reports than does diff.  This method appears to be the least vulnerable
// to synching up on blocks of "junk lines", though (like blank lines in
// ordinary text files, or maybe "<P>" lines in HTML files).  That may be
// because this is the only method of the 3 that has a *concept* of
// "junk" <wink>.
//
// Timing:  Basic R-O is cubic time worst case and quadratic time expected
// case.  SequenceMatcher is quadratic time for the worst case and has
// expected-case behavior dependent in a complicated way on how many
// elements the sequences have in common; best case time is linear.
type SequenceMatcher[T comparable] struct {
	a              []T
	b              []T
	b2j            map[T][]int
	IsJunk         func(T) bool
	autoJunk       bool
	bJunk          map[T]struct{}
	matchingBlocks []Match
}

type opts[T comparable] struct {
	isJunk   func(T) bool
	autoJunk bool
}

type SequenceMatcherOpt[T comparable] func(*opts[T])

func WithJunkFunc[T comparable](val func(T) bool) func(*opts[T]) {
	return func(o *opts[T]) {
		o.isJunk = val
	}
}

func WithAutoJunk[T comparable](val bool) func(*opts[T]) {
	return func(o *opts[T]) {
		o.autoJunk = val
	}
}

func NewMatcher[T comparable](configures ...SequenceMatcherOpt[T]) *SequenceMatcher[T] {
	o := &opts[T]{
		isJunk:   nil,
		autoJunk: true,
	}
	for _, configure := range configures {
		configure(o)
	}

	return &SequenceMatcher[T]{
		a:              nil,
		b:              nil,
		b2j:            nil,
		IsJunk:         o.isJunk,
		autoJunk:       o.autoJunk,
		bJunk:          nil,
		matchingBlocks: nil,
	}
}

// Set two sequences to be compared.
func (m *SequenceMatcher[T]) SetSeqs(a, b []T) {
	m.SetSeq1(a)
	m.SetSeq2(b)
}

// Set the first sequence to be compared. The second sequence to be compared is
// not changed.
//
// SequenceMatcher computes and caches detailed information about the second
// sequence, so if you want to compare one sequence S against many sequences,
// use .SetSeq2(s) once and call .SetSeq1(x) repeatedly for each of the other
// sequences.
func (m *SequenceMatcher[T]) SetSeq1(a []T) {
	if &a == &m.a {
		return
	}

	m.a = a
	m.matchingBlocks = nil
}

// Set the second sequence to be compared. The first sequence to be compared is
// not changed.
func (m *SequenceMatcher[T]) SetSeq2(b []T) {
	if &b == &m.b {
		return
	}

	m.b = b
	m.matchingBlocks = nil
	m.chainB()
}

func (m *SequenceMatcher[T]) chainB() {
	b2j := map[T][]int{}
	for i, s := range m.b {
		indices := b2j[s]
		indices = append(indices, i)
		b2j[s] = indices
	}

	// Purge junk elements
	m.bJunk = map[T]struct{}{}
	if m.IsJunk != nil {
		junk := m.bJunk

		for s := range b2j {
			if m.IsJunk(s) {
				junk[s] = struct{}{}
			}
		}

		for s := range junk {
			delete(b2j, s)
		}
	}

	// Purge remaining popular elements
	popular := map[T]struct{}{}
	n := len(m.b)

	if m.autoJunk && n >= 200 {
		ntest := n/100 + 1
		for s, indices := range b2j {
			if len(indices) > ntest {
				popular[s] = struct{}{}
			}
		}

		for s := range popular {
			delete(b2j, s)
		}
	}

	m.b2j = b2j
}

func (m *SequenceMatcher[T]) isBJunk(s T) bool {
	_, ok := m.bJunk[s]
	return ok
}

// Find longest matching block in a[alo:ahi] and b[blo:bhi].
//
// If IsJunk is not defined:
//
// Return (i,j,k) such that a[i:i+k] is equal to b[j:j+k], where
//
//	alo <= i <= i+k <= ahi
//	blo <= j <= j+k <= bhi
//
// and for all (i',j',k') meeting those conditions,
//
//	k >= k'
//	i <= i'
//	and if i == i', j <= j'
//
// In other words, of all maximal matching blocks, return one that
// starts earliest in a, and of all those maximal matching blocks that
// start earliest in a, return the one that starts earliest in b.
//
// If IsJunk is defined, first the longest matching block is
// determined as above, but with the additional restriction that no
// junk element appears in the block.  Then that block is extended as
// far as possible by matching (only) junk elements on both sides.  So
// the resulting block never matches on junk except as identical junk
// happens to be adjacent to an "interesting" match.
//
// If no blocks match, return (alo, blo, 0).
func (m *SequenceMatcher[T]) findLongestMatch(alo, ahi, blo, bhi int) Match {
	// CAUTION:  stripping common prefix or suffix would be incorrect.
	// E.g.,
	//    ab
	//    acab
	// Longest matching block is "ab", but if common prefix is
	// stripped, it's "a" (tied with "b").  UNIX(tm) diff does so
	// strip, so ends up claiming that ab is changed to acab by
	// inserting "ca" in the middle.  That's minimal but unintuitive:
	// "it's obvious" that someone inserted "ac" at the front.
	// Windiff ends up at the same place as diff, but by pairing up
	// the unique 'b's and then matching the first two 'a's.
	besti, bestj, bestsize := alo, blo, 0

	// find longest junk-free match
	// during an iteration of the loop, j2len[j] = length of longest
	// junk-free match ending with a[i-1] and b[j]
	j2len := map[int]int{}

	for i := alo; i != ahi; i++ {
		// look at all instances of a[i] in b; note that because
		// b2j has no junk keys, the loop is skipped if a[i] is junk
		newj2len := map[int]int{}

		for _, j := range m.b2j[m.a[i]] {
			// a[i] matches b[j]
			if j < blo {
				continue
			}

			if j >= bhi {
				break
			}

			k := j2len[j-1] + 1
			newj2len[j] = k

			if k > bestsize {
				besti, bestj, bestsize = i-k+1, j-k+1, k
			}
		}

		j2len = newj2len
	}

	// Extend the best by non-junk elements on each end.  In particular,
	// "popular" non-junk elements aren't in b2j, which greatly speeds
	// the inner loop above, but also means "the best" match so far
	// doesn't contain any junk *or* popular non-junk elements.
	for besti > alo && bestj > blo && !m.isBJunk(m.b[bestj-1]) &&
		m.a[besti-1] == m.b[bestj-1] {
		besti, bestj, bestsize = besti-1, bestj-1, bestsize+1
	}

	for besti+bestsize < ahi && bestj+bestsize < bhi &&
		!m.isBJunk(m.b[bestj+bestsize]) &&
		m.a[besti+bestsize] == m.b[bestj+bestsize] {
		bestsize++
	}

	// Now that we have a wholly interesting match (albeit possibly
	// empty!), we may as well suck up the matching junk on each
	// side of it too.  Can't think of a good reason not to, and it
	// saves post-processing the (possibly considerable) expense of
	// figuring out what to do with it.  In the case of an empty
	// interesting match, this is clearly the right thing to do,
	// because no other kind of match is possible in the regions.
	for besti > alo && bestj > blo && m.isBJunk(m.b[bestj-1]) &&
		m.a[besti-1] == m.b[bestj-1] {
		besti, bestj, bestsize = besti-1, bestj-1, bestsize+1
	}

	for besti+bestsize < ahi && bestj+bestsize < bhi &&
		m.isBJunk(m.b[bestj+bestsize]) &&
		m.a[besti+bestsize] == m.b[bestj+bestsize] {
		bestsize++
	}

	return Match{A: besti, B: bestj, Size: bestsize}
}

// Return list of triples describing matching subsequences.
//
// Each triple is of the form (i, j, n), and means that
// a[i:i+n] == b[j:j+n].  The triples are monotonically increasing in
// i and in j. It's also guaranteed that if (i, j, n) and (i', j', n') are
// adjacent triples in the list, and the second is not the last triple in the
// list, then i+n != i' or j+n != j'. IOW, adjacent triples never describe
// adjacent equal blocks.
//
// The last triple is a dummy, (len(a), len(b), 0), and is the only
// triple with n==0.
func (m *SequenceMatcher[T]) GetMatchingBlocks() []Match {
	if m.matchingBlocks != nil {
		return m.matchingBlocks
	}

	la, lb := len(m.a), len(m.b)
	queue := []struct {
		alo, ahi, blo, bhi int
	}{{0, la, 0, lb}}

	var matched []Match

	for len(queue) > 0 {
		block := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		alo, ahi, blo, bhi := block.alo, block.ahi, block.blo, block.bhi
		match := m.findLongestMatch(alo, ahi, blo, bhi)
		i, j, k := match.A, match.B, match.Size

		if k > 0 {
			matched = append(matched, match)

			if alo < i && blo < j {
				queue = append(queue, struct {
					alo, ahi, blo, bhi int
				}{alo, i, blo, j})
			}

			if i+k < ahi && j+k < bhi {
				queue = append(queue, struct {
					alo, ahi, blo, bhi int
				}{i + k, ahi, j + k, bhi})
			}
		}
	}

	sort.Slice(matched, func(i, j int) bool {
		if matched[i].A != matched[j].A {
			return matched[i].A < matched[j].A
		}

		return matched[i].B < matched[j].B
	})

	// It's possible that we have adjacent equal blocks in the
	// matching_blocks list now.
	nonAdjacent := []Match{}
	i1, j1, k1 := 0, 0, 0

	for _, b := range matched {
		// Is this block adjacent to i1, j1, k1?
		i2, j2, k2 := b.A, b.B, b.Size

		if i1+k1 == i2 && j1+k1 == j2 {
			// Yes, so collapse them -- this just increases the length of
			// the first block by the length of the second, and the first
			// block so lengthened remains the block to compare against.
			k1 += k2
		} else {
			// Not adjacent.  Remember the first block (k1==0 means it's
			// the dummy we started with), and make the second block the
			// new block to compare against.
			if k1 > 0 {
				nonAdjacent = append(nonAdjacent, Match{i1, j1, k1})
			}

			i1, j1, k1 = i2, j2, k2
		}
	}

	if k1 > 0 {
		nonAdjacent = append(nonAdjacent, Match{i1, j1, k1})
	}

	nonAdjacent = append(nonAdjacent, Match{len(m.a), len(m.b), 0})
	m.matchingBlocks = nonAdjacent

	return m.matchingBlocks
}
