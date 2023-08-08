package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mustwork/aezeed-recover/pkg/tview/responsive"
	"github.com/rivo/tview"
)

const (
	QuitButtonLabel    = "Quit"
	ConsentButtonLabel = "Agree"
)

type Info struct {
	*tview.Grid
	textView *responsive.ResponsiveTextView
	form     *tview.Form
}

func (o *Info) SetBorder(b bool) *Info {
	o.Grid.SetBorder(b)
	return o
}
func (o *Info) GetQuitButton() *tview.Button {
	idx := o.form.GetButtonIndex(QuitButtonLabel)
	return o.form.GetButton(idx)
}
func (o *Info) GetConfirmButton() *tview.Button {
	idx := o.form.GetButtonIndex(ConsentButtonLabel)
	return o.form.GetButton(idx)
}
func (o *Info) SetQuitButtonDisabled(b bool) *Info {
	o.GetQuitButton().SetDisabled(b)
	return o
}
func (o *Info) SetConfirmButtonDisabled(b bool) *Info {
	o.GetConfirmButton().SetDisabled(b)
	return o
}
func (o *Info) SetQuitButtonSelectedFunc(fn func()) *Info {
	o.GetQuitButton().SetSelectedFunc(fn)
	return o
}
func (o *Info) SetConfirmButtonSelectedFunc(fn func()) *Info {
	o.GetConfirmButton().SetSelectedFunc(fn)
	return o
}

func NewInfoView(InfoText string) *Info {

	textView := responsive.NewResponsiveTextView().
		SetScrollable(true).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetText(InfoText)

	form := tview.NewForm().
		SetButtonsAlign(tview.AlignRight).
		AddButton(QuitButtonLabel, func() {}).
		AddButton(ConsentButtonLabel, func() {})
	form.SetBorderPadding(1, 0, 1, 1)

	grid := tview.NewGrid().
		SetRows(0, 2).
		SetColumns(0).
		AddItem(textView, 0, 0, 1, 1, 0, 0, false).
		AddItem(form, 1, 0, 1, 1, 0, 0, true)

	info := Info{grid, textView, form}
	info.SetConfirmButtonDisabled(true)

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Detecting whether the last line is visible
		lineCount := textView.GetLineCount()
		_, _, _, height := textView.GetInnerRect()
		row, _ := textView.GetScrollOffset()
		if row+height >= lineCount {
			info.SetConfirmButtonDisabled(false)
		}
		switch event.Key() {
		case tcell.KeyLeft:
			return tcell.NewEventKey(tcell.KeyBacktab, event.Rune(), event.Modifiers())
		case tcell.KeyRight:
			return tcell.NewEventKey(tcell.KeyTab, event.Rune(), event.Modifiers())
		case tcell.KeyCtrlP, tcell.KeyPgUp:
			textView.InputHandler()(tcell.NewEventKey(tcell.KeyUp, event.Rune(), event.Modifiers()), nil)
			return nil
		case tcell.KeyCtrlN, tcell.KeyPgDn:
			textView.InputHandler()(tcell.NewEventKey(tcell.KeyDown, event.Rune(), event.Modifiers()), nil)
			return nil
		default:
			textView.InputHandler()(event, nil)
			return event
		}
	})

	return &info
}
