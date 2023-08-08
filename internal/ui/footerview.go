package ui

import (
	"fmt"
	"github.com/mustwork/aezeed-recover/pkg/tview/responsive"
	"github.com/rivo/tview"
)

type Footer struct {
	*responsive.ResponsiveTextView
}

func (f *Footer) SetBorder(b bool) *Footer {
	f.ResponsiveTextView.SetBorder(b)
	return f
}
func (f *Footer) Clear() *Footer {
	f.SetText("")
	return f
}
func (f *Footer) SetInfo(s string) *Footer {
	// TODO - use colors and indents, possibly with a grid
	f.SetText(fmt.Sprintf("INFO: %s", s))
	return f
}

func NewFooterView() *Footer {
	textView := responsive.NewResponsiveTextView().SetTextAlign(tview.AlignLeft)
	return &Footer{textView}
}
