package responsive

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ResponsiveListItem struct {
	parent         *ResponsiveList
	Texts          []string
	SecondaryTexts []string
	Help           string
	Shortcut       rune
	Selected       func()
	Enabled        bool
	EnabledStyle   tcell.Style
	DisabledStyle  tcell.Style
}

func (i *ResponsiveListItem) SetTexts(texts ...string) *ResponsiveListItem {
	// TODO - must manipulate parent
	i.Texts = texts
	return i
}
func (i *ResponsiveListItem) SetSecondaryTexts(texts []string) *ResponsiveListItem {
	// TODO - must manipulate parent
	i.SecondaryTexts = texts
	return i
}
func (i *ResponsiveListItem) SetHelp(help string) *ResponsiveListItem {
	i.Help = help
	return i
}
func (i *ResponsiveListItem) HasHelp() bool {
	return i.Help != ""
}
func (i *ResponsiveListItem) SetShortcut(shortcut rune) *ResponsiveListItem {
	// TODO - must manipulate parent
	i.Shortcut = shortcut
	return i
}
func (i *ResponsiveListItem) SetEnabled(b bool) *ResponsiveListItem {
	// TODO - must manipulate parent
	i.Enabled = b
	return i
}
func (i *ResponsiveListItem) SetEnabledStyle(s tcell.Style) *ResponsiveListItem {
	// TODO - must manipulate parent
	i.EnabledStyle = s
	return i
}
func (i *ResponsiveListItem) SetDisabledStyle(s tcell.Style) *ResponsiveListItem {
	// TODO - must manipulate parent
	i.DisabledStyle = s
	return i
}
func (i *ResponsiveListItem) SetSelected(fn func()) *ResponsiveListItem {
	// TODO - must manipulate parent
	i.Selected = fn
	return i
}
func (i *ResponsiveListItem) getText(w int) string {
	if len(i.Texts) == 0 {
		return ""
	}
	for _, s := range i.Texts {
		if len(s) <= w {
			var fg tcell.Color
			var attrs tcell.AttrMask
			if i.Enabled {
				fg, _, attrs = i.EnabledStyle.Decompose()
			} else {
				fg, _, attrs = i.DisabledStyle.Decompose()
			}
			r, g, b := fg.RGB()
			if attrs&tcell.AttrBold == tcell.AttrBold {
				return fmt.Sprintf("[#%02x%02x%02x::b]%s[-::-]", r, g, b, s) // FIXME - bold is not disabled
			} else {
				return fmt.Sprintf("[#%02x%02x%02x::]%s[-::]", r, g, b, s)
			}
		}
	}
	return i.Texts[len(i.Texts)-1]
}
func (i *ResponsiveListItem) getSecondaryText(w int) string {
	if len(i.SecondaryTexts) == 0 {
		return ""
	}
	for _, s := range i.SecondaryTexts {
		if len(s) <= w {
			return s
		}
	}
	return i.SecondaryTexts[len(i.SecondaryTexts)-1]
}
func NewResponsiveListItem() *ResponsiveListItem {
	return &ResponsiveListItem{Enabled: true}
}

type ResponsiveList struct {
	*tview.List
	Items []*ResponsiveListItem
}

func (rl *ResponsiveList) SetEnabledColor(s tcell.Style) *ResponsiveList {
	for _, item := range rl.Items {
		item.SetEnabledStyle(s)
	}
	return rl
}
func (rl *ResponsiveList) SetDisabledColor(s tcell.Style) *ResponsiveList {
	for _, item := range rl.Items {
		item.SetDisabledStyle(s)
	}
	return rl
}
func (rl *ResponsiveList) Draw(screen tcell.Screen) {
	rl.List.Draw(screen)
	l, t, r, b := rl.List.GetInnerRect()
	w := r - l - 4 // not sure how this offset is calculated
	h := b - t - 2 // not sure how this offset is calculated
	for i, item := range rl.Items {
		rl.SetItemText(i, item.getText(w), item.getSecondaryText(w))
	}
	rl.ShowSecondaryText(h > len(rl.Items))
}

// TODO - consider being more idiomatic, using AddItem
func NewResponsiveList(items ...*ResponsiveListItem) *ResponsiveList {
	l := tview.NewList().ShowSecondaryText(false).SetMainTextStyle(tcell.StyleDefault)
	list := &ResponsiveList{l, items}
	for _, item := range items {
		item.parent = list
		// As Draw is called later, we don't need to set the text here.
		l.AddItem("", "", item.Shortcut, item.Selected)
	}
	return list
}
