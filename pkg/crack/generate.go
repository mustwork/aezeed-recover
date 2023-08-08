package crack

import (
	"github.com/lightningnetwork/lnd/aezeed"
	"golang.org/x/exp/slices"
	"time"
)

const (
	NumberOfWorkers = 32
)

type Result struct {
	Trials     Trials
	Password   string
	CipherSeed *aezeed.CipherSeed
	Highlight  []int
}

func (r Result) Exhausted() bool {
	if r.Trials.Tried < r.Trials.Total {
		return false
	} else if r.Trials.Tried == r.Trials.Total {
		return true
	} else {
		panic("tried more than total")
	}
}

func FindMissingWords(s MnemonicSettings, callback func(progress Result)) {
	assertConfigValid(s, false)
	start := time.Now()
	missingWords := s.MnemonicLength - len(s.Mnemonic)
	pool := NewWorkerPool(NumberOfWorkers, s.Password)
	go func() {
		tried := 0
		// actually less than 8 + 23*2048, because of version bits,
		// and even less because of timestamp
		total := 1
		for i := 0; i < missingWords; i++ {
			total *= (s.MnemonicLength - i) * len(s.Wordlist)
		}
		for seed := range pool.Out {
			tried++
			if seed != nil {
				pool.Close()
				mnemonic, err := seed.ToMnemonic([]byte(s.Password))
				if err != nil {
					panic(err)
				}
				xs := complementedIndices(s.Mnemonic, mnemonic)
				callback(Result{Trials{start, total, tried}, s.Password, seed, xs})
				// if we don't return here there may be trailing events from pool.Out from other workers
				return
			}
			callback(Result{Trials{start, total, tried}, s.Password, seed, nil})
		}
	}()
	go pool.Consume(ComplementedMnemonics(s))
}

func assertConfigValid(s MnemonicSettings, expectComplete bool) {
	if expectComplete && len(s.Mnemonic) != s.MnemonicLength {
		panic("mnemonic incomplete")
	} else if !expectComplete && len(s.Mnemonic) == s.MnemonicLength {
		panic("mnemonic too long")
	}
	for _, word := range s.Mnemonic {
		if !slices.Contains(s.Wordlist, word) {
			panic("mnemonic contains word not in wordlist")
		}
	}
	if len(s.Mnemonic) == s.MnemonicLength {
		m := aezeed.Mnemonic(s.Mnemonic)
		_, err := m.ToCipherSeed([]byte(s.Password))
		if err == nil {
			panic("mnemonic is already valid")
		}
	}
}

// complementedIndices returns the indices of the words that have been added.
func complementedIndices(xs []string, ys aezeed.Mnemonic) (result []int) {
	x := 0
	for i, y := range ys {
		if x >= len(xs) || xs[x] != y {
			result = append(result, i)
		} else {
			x++
		}
	}
	return
}

// ComplementedMnemonics iterates over all possible mnemonics that complete the given partial mnemonic -- unvalidated.
func ComplementedMnemonics(s MnemonicSettings) <-chan *aezeed.Mnemonic {
	missingWords := s.MnemonicLength - len(s.Mnemonic)
	out := make(chan *aezeed.Mnemonic)
	go func() {
		// TODO - version and timestamp can be optimized explicitly
		// TODO - start from behind, as it seems more likely that errors are made at the end rather than at the beginning
		defer close(out)
		for i := NewInsertIndexer(s, missingWords); !i.Done(); i.Increment() {
			out <- i.Apply(s.Mnemonic)
		}
	}()
	return out
}

func FindSwappedWords(s MnemonicSettings, callback func(progress Result)) {
	start := time.Now()
	assertConfigValid(s, true)
	go func() {
		total := 23
		for i := 0; i < 23; i++ {
			m := make([]string, s.MnemonicLength)
			copy(m, s.Mnemonic)
			m[i], m[i+1] = m[i+1], m[i]
			mnemonic := aezeed.Mnemonic(m)
			cipherSeed, err := mnemonic.ToCipherSeed([]byte(s.Password))
			if err != nil {
				callback(Result{Trials{start, total, i + 1}, s.Password, nil, nil})
				continue
			}
			xs := findSwappedIndices(s.Mnemonic, mnemonic)
			callback(Result{Trials{start, total, i}, s.Password, cipherSeed, xs})
			break
		}
	}()
}

func findSwappedIndices(xs []string, ys aezeed.Mnemonic) (result []int) {
	return findWrongIndices(xs, ys)
}

// FindWrongWord finds a single wrong word in the mnemonic.
func FindWrongWord(s MnemonicSettings, callback func(progress Result)) {
	start := time.Now()
	assertConfigValid(s, true)
	go func() {
		total := s.MnemonicLength * len(s.Wordlist)
		tries := 1
		for i := 0; i < s.MnemonicLength; i++ {
			for _, word := range s.Wordlist {
				m := make([]string, s.MnemonicLength)
				copy(m, s.Mnemonic)
				m[i] = word
				mnemonic := aezeed.Mnemonic(m)
				cipherSeed, err := mnemonic.ToCipherSeed([]byte(s.Password))
				if err != nil {
					callback(Result{Trials{start, total, tries}, s.Password, nil, nil})
					tries++
					continue
				}
				xs := findWrongIndices(s.Mnemonic, mnemonic)
				callback(Result{Trials{start, total, tries}, s.Password, cipherSeed, xs})
				return
			}
		}
	}()
}

func findWrongIndices(xs []string, ys aezeed.Mnemonic) (result []int) {
	x := 0
	for i, y := range ys {
		if xs[x] != y {
			result = append(result, i)
		}
		x++
	}
	return
}
