package crack

import (
	"math"
	"time"
)

type MnemonicSettings struct {
	Wordlist       []string
	Password       string
	Mnemonic       []string
	MnemonicLength int
}

type Trials struct {
	Start time.Time
	Total int
	Tried int
}

func (t *Trials) Percentage() float64 {
	return float64(t.Tried) / float64(t.Total) * 100
}
func (t *Trials) Elapsed() time.Duration {
	return time.Since(t.Start)
}
func (t *Trials) Remaining() time.Duration {
	if t.Tried == 0 {
		return time.Duration(math.MaxInt64)
	}
	perTrial := t.Elapsed() / time.Duration(t.Tried)
	return perTrial * time.Duration(t.Total-t.Tried)
}
