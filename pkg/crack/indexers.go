package crack

import (
	"github.com/lightningnetwork/lnd/aezeed"
)

// Indexer is an interface for iterating over the search space of an incomplete or wrong mnemonic.
// (Named Indexer for lack of a better word).
type Indexer interface {
	Done() bool
	Increment() Indexer
	Apply([]string) *aezeed.Mnemonic
	// TODO - the Indexer should be capable of returning the current position, the total size of the search space, and the remaining search space.
	//Size() int
	//Position() int
	//Remaining() int
}

// insertSubIndexer keeps track of a single substitution.
type insertSubIndexer struct {
	position int
	word     int
}

// InsertIndexer is used to keep track of multiple insertion- and word indices when complementing a mnemonic.
type InsertIndexer struct {
	settings MnemonicSettings
	xs       []insertSubIndexer
}

// Done returns true if the indexer has reached the end of the search space.
func (i *InsertIndexer) Done() bool {
	for k := 0; k < len(i.xs); k++ {
		x := i.xs[len(i.xs)-1-k]
		if (x.word < len(i.settings.Wordlist)-1) || (x.position-k < i.settings.MnemonicLength-1) {
			return false
		}
	}
	return true
}

// Increment forwards the current position by the following rules:
// If the last sub-indexer's current word index is less than 2047, increment it.
// If the last sub-indexer's current word index is 2047, and the current position is less than 23, increment the current position and reset the word index.
// If the last sub-indexer's current word index is 2047, and the current position is 23, increment the second to last sub-indexer according to the same rules, and reset the last sub-indexer.
func (i *InsertIndexer) Increment() {
	for k := 0; k < len(i.xs); k++ {
		idx := len(i.xs) - 1 - k
		x := i.xs[idx] // watch out, this is a copy
		if x.word < len(i.settings.Wordlist)-1 {
			i.xs[idx].word++
			return
		} else if x.position-k < i.settings.MnemonicLength-1 {
			i.xs[idx].position++
			i.xs[idx].word = 0
			return
		} else {
			if idx == 0 {
				return // we're done
			}
			i.xs[idx].position = i.xs[idx-1].position + 1
			i.xs[idx].word = 0
			continue
		}
	}
	panic("Increment called on i that is already done")
}

// Apply returns a new mnemonic with the missing words inserted at the positions according to the indexer.
func (i *InsertIndexer) Apply(xs []string) *aezeed.Mnemonic {
	start := 0
	mnemonic := aezeed.Mnemonic{}
	for k, x := range i.xs {
		copy(mnemonic[start:x.position], xs[start-k:x.position-k])
		mnemonic[x.position] = aezeed.DefaultWordList[x.word]
		start = x.position + 1
	}
	copy(mnemonic[start:], xs[start-len(i.xs):])
	return &mnemonic
}

// NewInsertIndexer creates an Indexer for n amount of missing words.
func NewInsertIndexer(s MnemonicSettings, n int) *InsertIndexer {
	xs := make([]insertSubIndexer, n)
	for i := 0; i < n; i++ {
		xs[i] = insertSubIndexer{position: i, word: 0}
	}
	return &InsertIndexer{s, xs}
}
