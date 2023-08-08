package responsive

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mitchellh/go-wordwrap"
	"github.com/rivo/tview"
	"strings"
)

type ResponsiveTextView struct {
	*tview.TextView
	Text string
}

func (view *ResponsiveTextView) GetLineCount() int {
	return view.TextView.GetOriginalLineCount()
}
func (view *ResponsiveTextView) Draw(screen tcell.Screen) {
	view.TextView.Draw(screen)

	// This method wraps text and takes special sequences into account.

	builder := strings.Builder{}
	escape := false
	skip := false
	for _, r := range view.Text {
		if skip {
			if r == ']' {
				skip = false
			}
			continue
		}
		switch r {
		case '\\':
			escape = true
		case '[':
			if escape {
				builder.WriteRune(r)
			} else {
				skip = true
			}
		default:
			builder.WriteRune(r)
			escape = false
		}
	}
	s := builder.String()
	_, _, w, _ := view.GetRect()
	wrapped := wordwrap.WrapString(s, uint(w-1))
	builder = strings.Builder{}
	k := 0
	newline := false
	for i := 0; i < len(wrapped); i++ {
		if wrapped[i] == '\n' {
			// newlines are determined by the wrapped string, not the original one.
			// TODO - one caveat is that newlines are written BEFORE closing special sequences.
			//  it does not seem to make any difference visually, though.
			builder.WriteByte('\n')
			newline = true
			continue
		}
		for view.Text[k] != wrapped[i] {
			// catching up special characters, ignoring spaces after newline
			if newline && view.Text[k] == ' ' {
				newline = false
				k++
				continue
			}
			// also ignoring newlines, as they are determined by the wrapped string
			if view.Text[k] != '\n' {
				builder.WriteByte(view.Text[k])
			}
			k++
		}
		// writing all other characters
		builder.WriteByte(view.Text[k])
		k++
	}
	view.TextView.SetText(builder.String())
}
func (view *ResponsiveTextView) SetTextColor(color tcell.Color) *ResponsiveTextView {
	view.TextView.SetTextColor(color)
	return view
}
func (view *ResponsiveTextView) SetTextAlign(align int) *ResponsiveTextView {
	view.TextView.SetTextAlign(align)
	return view
}
func (view *ResponsiveTextView) SetDynamicColors(b bool) *ResponsiveTextView {
	view.TextView.SetDynamicColors(b)
	return view
}
func (view *ResponsiveTextView) SetText(text string) *ResponsiveTextView {
	view.Text = text
	view.TextView.SetText(text)
	return view
}
func (view *ResponsiveTextView) SetScrollable(b bool) *ResponsiveTextView {
	view.TextView.SetScrollable(b)
	return view
}

func NewResponsiveTextView() *ResponsiveTextView {
	return &ResponsiveTextView{
		TextView: tview.NewTextView(),
	}
}
