package ui

import (
	"github.com/rivo/tview"
)

type Header struct {
	*tview.TextView
}

func (h *Header) SetTitle(s string) *Header {
	h.SetText(s)
	return h
}
func (h *Header) SetBorder(b bool) *Header {
	h.TextView.SetBorder(b)
	return h
}

func NewHeaderView() *Header {
	textView := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	return &Header{textView}
}
