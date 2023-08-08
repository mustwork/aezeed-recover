package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/lightningnetwork/lnd/aezeed"
	"github.com/lightningnetwork/lnd/keychain"
	"github.com/mustwork/aezeed-recover/pkg/crack"
	"github.com/rivo/tview"
	"golang.org/x/exp/slices"
	"time"
)

const (
	Title                   = "aezeed-recover"
	InfoPageName            = "infoPage"
	MenuPageName            = "menuPage"
	ResultPageName          = "resultPage"
	PasswordInputFooterHint = "The empty password and 'aezeed' (default) are synonymous."
	// The string concatenation is necessary to prevent IntelliJ from reformatting the string as SQL:
	MnemonicInputFooterHint = "Select mnemonic" + " words or paste from clipboard (or ESC to menu)."
)

func CreateApplication(theme Theme, config *Config, info string) *tview.Application {

	tview.Styles = theme.Theme

	app := tview.NewApplication()

	headerView := NewHeaderView()
	footerView := NewFooterView()
	// left column
	pages := tview.NewPages()
	infoPage := NewInfoView(info)
	menuPage, menuCallbacks := NewMenuView(theme, config.DevMode)
	resultPage := NewResultView(app.SetFocus)
	// right column
	seedView := NewSeedView(theme, config, app.SetFocus)

	pages.
		AddPage(InfoPageName, infoPage, true, true).
		AddPage(MenuPageName, menuPage, true, false).
		AddPage(ResultPageName, resultPage, true, false)

	navigateTo := func(pageName string) {
		for _, page := range []string{InfoPageName, MenuPageName, ResultPageName} {
			if page != pageName {
				pages.HidePage(page)
			}
		}
		pages.ShowPage(pageName)
		go func() {
			app.QueueUpdateDraw(func() {
				switch pageName {
				case InfoPageName:
					app.SetFocus(infoPage)
				case MenuPageName:
					app.SetFocus(menuPage)
				case ResultPageName:
					app.SetFocus(resultPage)
				}
			})
		}()
	}
	navigateToMenu := func() {
		navigateTo(MenuPageName)
	}
	navigateToInfo := func() {
		navigateTo(InfoPageName)
	}
	navigateToResult := func() {
		navigateTo(ResultPageName)
	}

	headerHeight := 1
	footerHeight := 1
	// perfect ratio at 80
	leftColumnWidth := -50
	rightColumnWidth := -30
	mainGrid := tview.NewGrid().
		SetRows(headerHeight, 0, footerHeight).
		SetColumns(leftColumnWidth, rightColumnWidth).
		AddItem(headerView, 0, 0, 1, 2, 0, 0, false).
		AddItem(pages, 1, 0, 1, 1, 0, 0, false).
		AddItem(seedView, 1, 1, 1, 1, 0, 0, false).
		AddItem(footerView, 2, 0, 1, 2, 0, 0, false)

	headerView.
		SetTitle(Title).
		SetBorder(false)

	footerView.
		SetBorder(false)

	showFooterInfo := func(s string) func() {
		return func() {
			go func() {
				app.QueueUpdateDraw(func() {
					footerView.SetInfo(s)
				})
			}()
		}
	}

	progressHandler := func(onSuccess func(), onFailure func()) func(progress crack.Result) {
		// menu must be disabled, in order not to interfere with the brute force attempt
		app.SetFocus(footerView)
		ch := make(chan crack.Result)
		tick := time.NewTicker(100 * time.Millisecond)
		go func() {
			progress := crack.Result{}
			defer close(ch)
			for {
				select {
				case <-tick.C:
					showFooterInfo(fmt.Sprintf("tried %d/%d (%.2f%%) elapsed: %v remaining: ~%v (max)",
						progress.Trials.Tried,
						progress.Trials.Total,
						progress.Trials.Percentage(),
						progress.Trials.Elapsed().Round(time.Second),
						progress.Trials.Remaining().Round(time.Minute)))()
				case p, ok := <-ch:
					if !ok {
						return
					}
					progress = p
				}
			}
		}()
		return func(result crack.Result) {
			ch <- result
			if result.CipherSeed != nil {
				mnemonic, err := result.CipherSeed.ToMnemonic([]byte(result.Password))
				if err != nil {
					panic(err)
				}
				onSuccess()
				seedView.SetWords(mnemonic[:])
				seedView.HighlightWords(result.Highlight)
				tick.Stop()
				navigateToResult()
				go func() { app.QueueUpdateDraw(func() {}) }() // force redraw for proper line breaks in result view
			} else if result.Exhausted() {
				onFailure()
				showFooterInfo(fmt.Sprintf("tried %d/%d - exhausted, no result", result.Trials.Tried, result.Trials.Total))()
				navigateToMenu()
				tick.Stop()
			}
		}
	}

	infoPage.
		SetQuitButtonSelectedFunc(app.Stop).
		SetConfirmButtonSelectedFunc(navigateToMenu).
		SetBorder(true)

	menuPage.
		SetShowInfoHandler(navigateToInfo).
		SetDisplaySeedHandler(navigateToResult).
		SetQuitHandler(func() {
			_, err := seedView.GetCipherSeed()
			if err != nil { // there is no valid seed
				app.Stop()
			} else {
				navigateToResult()
			}
		}).
		SetCreateRandomMnemonicHandler(func() {
			menuPage.ClearBruteForceResults()
			password := seedView.GetPassword()
			cipherSeed, err := aezeed.New(keychain.CurrentKeyDerivationVersion, nil, time.Now())
			mnemnoic, err := cipherSeed.ToMnemonic([]byte(password))
			if err != nil {
				panic(err)
			}
			seedView.SetWords(mnemnoic[:])
		}).
		SetEditMnemonicHandler(func() {
			app.SetFocus(seedView)
		}).
		SetClearMnemonicHandler(func() {
			menuPage.ClearBruteForceResults()
			seedView.SetWords(make([]string, 24))
		}).
		SetFindMissingWordsHandler(func() {
			crack.FindMissingWords(seedView.GetSettings(), progressHandler(menuPage.SetFoundMissingWords, menuPage.SetExhaustedMissingWords))
		}).
		SetFindSwappedWordsHandler(func() {
			crack.FindSwappedWords(seedView.GetSettings(), progressHandler(menuPage.SetFoundSwappedWords, menuPage.SetExhaustedSwappedWords))
		}).
		SetFindWrongWordHandler(func() {
			crack.FindWrongWord(seedView.GetSettings(), progressHandler(menuPage.SetFoundWrongWord, menuPage.SetExhaustedWrongWord))
		}).
		SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			item := menuPage.Items[index]
			if item.HasHelp() {
				footerView.SetInfo(item.Help)
			} else {
				footerView.Clear()
			}
		}).
		SetFocusFunc(func() {
			item := menuPage.Items[menuPage.GetCurrentItem()]
			if item.HasHelp() {
				footerView.SetInfo(item.Help)
			} else {
				footerView.Clear()
			}
		}).
		SetBorder(true)

	resultPage.
		SetQuitButtonSelectedFunc(app.Stop).
		SetReturnButtonSelectedFunc(navigateToMenu).
		SetBorder(true)

	seedView.
		SetOnMnemonicChangedFunc(func(mnemonic []string, valid bool) {
			menuPage.ClearBruteForceResults() // TODO - only clear if mnemonic actually changed
			if valid {
				menuCallbacks.OnMnemonicValid()
			} else {
				invalidWords := false
				for _, word := range mnemonic {
					if !slices.Contains(config.Wordlist, word) {
						invalidWords = true
						break
					}
				}
				if invalidWords {
					menuCallbacks.OnMnemonicInsufficient()
				} else if len(mnemonic) == config.MnemonicLength {
					menuCallbacks.OnMnemonicInvalid()
				} else if len(mnemonic) >= config.MnemonicMinLength {
					menuCallbacks.OnMnemonicPartial()
				} else {
					menuCallbacks.OnMnemonicInsufficient()
				}
			}
		}).
		SetPasswordInputFocusFunc(showFooterInfo(PasswordInputFooterHint)).
		SetMnemonicInputFocusFunc(showFooterInfo(MnemonicInputFooterHint)).
		SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEscape:
				app.SetFocus(menuPage)
			}
		})

	app.SetRoot(mainGrid, true)

	if config.Consent {
		navigateToMenu()
		seedView.SetWords([]string{}) // this triggers the mnemonic hint in seed view
	} else {
		app.SetFocus(infoPage)
	}

	return app
}
