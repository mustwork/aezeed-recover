package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mustwork/aezeed-recover/internal/util"
	"github.com/mustwork/aezeed-recover/pkg/tview/responsive"
	"github.com/rivo/tview"
)

const (
	// Url = "https://recovery.mustwork.de" // TODO - implement LnUrl
	Url               = "bc1qfkv06pjcgq280f42nu66ktdfeer7f06lmtgerp" // TODO - parse from README.md
	ReturnButtonLabel = "Return"
	Plea              = `Success! 

Please consider a donation to keep this tool awesome.`
)

type ResultView struct {
	*tview.Grid
	returnButton *tview.Button
	quitButton   *tview.Button
}

func (v *ResultView) SetBorder(b bool) *ResultView {
	v.Grid.SetBorder(b)
	return v
}
func (v *ResultView) SetQuitButtonSelectedFunc(f func()) *ResultView {
	v.quitButton.SetSelectedFunc(f)
	return v
}
func (v *ResultView) SetReturnButtonSelectedFunc(f func()) *ResultView {
	v.returnButton.SetSelectedFunc(f)
	return v
}

func NewResultView(setFocus func(tview.Primitive) *tview.Application) *ResultView {

	qrCodeView := NewQRCodeView()

	textView := responsive.NewResponsiveTextView().
		SetText(Plea)

	returnButton := tview.NewButton(ReturnButtonLabel)
	quitButton := tview.NewButton(QuitButtonLabel)

	form := tview.NewGrid().
		SetRows(0, 1, 1, 1).
		SetColumns(1, 0, 1).
		AddItem(returnButton, 1, 1, 1, 1, 0, 0, true).
		AddItem(quitButton, 3, 1, 1, 1, 0, 0, true)
	form.
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyTab, tcell.KeyBacktab, tcell.KeyDown, tcell.KeyUp, tcell.KeyLeft, tcell.KeyRight:
				if returnButton.HasFocus() {
					setFocus(quitButton)
					return nil
				}
				if quitButton.HasFocus() {
					setFocus(returnButton)
					return nil
				}
			}
			return event
		})

	grid := tview.NewGrid().
		SetRows(qrCodeView.Height-5, 5, 1, 0).
		SetColumns(qrCodeView.Width, 1, 0).
		AddItem(qrCodeView, 0, 0, 2, 1, 0, 0, false).
		AddItem(tview.NewTextView().SetText(Url), 3, 0, 1, 1, 0, 0, false).
		AddItem(textView, 0, 2, 1, 1, 0, 0, false).
		AddItem(form, 1, 2, 3, 1, 0, 0, true)

	return &ResultView{grid, returnButton, quitButton}
}

type QRCodeView struct {
	*tview.TextView
	Height int
	Width  int
}

func NewQRCodeView() *QRCodeView {
	code, width, height, err := util.CreateQRCode(Url, 1)
	if err != nil {
		panic(err)
	}
	textView := tview.NewTextView().
		SetTextStyle(tcell.StyleDefault.Bold(true).Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)).
		SetText(code)
	return &QRCodeView{textView, height, width}
}
