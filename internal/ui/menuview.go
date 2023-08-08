package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mustwork/aezeed-recover/pkg/tview/responsive"
	"regexp"
)

const (
	ShowInfoHelp             = "return to info page."
	ShowInfoShortcut         = 'i'
	RandomMnemonicHelp       = "create a random mnemonic."
	RandomMnemonicShortcut   = 'r'
	EditMnemonicHelp         = "edit or paste mnemonic words."
	EditMnemonicShortcut     = 'e'
	ClearMnemonicHelp        = "clear mnemonic words."
	ClearMnemonicShortcut    = 'c'
	FindMissingWordsHelp     = "brute force missing words."
	FindMissingWordsShortcut = 'm'
	FindSwappedWordsHelp     = "brute force swapped words."
	FindSwappedWordsShortcut = 's'
	FindWrongWordHelp        = "find a single wrong word."
	FindWrongWordShortcut    = 'w'
	DisplaySeedHelp          = "display additional seed information."
	DisplaySeedShortcut      = 'd'
	QuitHelp                 = "quit application and return to shell."
	QuitShortcut             = 'q'
)

type MenuCallbacks struct {
	OnMnemonicValid        func()
	OnMnemonicInvalid      func()
	OnMnemonicPartial      func()
	OnMnemonicInsufficient func()
}

type Menu struct {
	*responsive.ResponsiveList
	handlers map[rune]func()
}

func (m *Menu) SetBorder(b bool) *Menu {
	m.ResponsiveList.SetBorder(b)
	return m
}
func (m *Menu) GetItem(shortcut rune) (int, *responsive.ResponsiveListItem) {
	for idx, item := range m.Items {
		if item.Shortcut == shortcut {
			return idx, item
		}
	}
	return -1, nil
}
func (m *Menu) SetPreviousEnabledItem() *Menu {
	l := len(m.Items)
	for i := m.GetCurrentItem(); ; {
		i = (i - 1 + l) % l
		if m.Items[i].Enabled {
			m.SetCurrentItem(i)
			break
		}
	}
	return m
}
func (m *Menu) SetNextEnabledItem() *Menu {
	l := len(m.Items)
	for i := m.GetCurrentItem(); ; {
		i = (i + 1) % l
		if m.Items[i].Enabled {
			m.SetCurrentItem(i)
			break
		}
	}
	return m
}
func (m *Menu) setHandler(r rune, fn func()) *Menu {
	m.handlers[r] = fn
	return m
}
func (m *Menu) SetShowInfoHandler(fn func()) *Menu {
	return m.setHandler(ShowInfoShortcut, fn)
}
func (m *Menu) SetCreateRandomMnemonicHandler(fn func()) *Menu {
	return m.setHandler(RandomMnemonicShortcut, fn)
}
func (m *Menu) SetEditMnemonicHandler(fn func()) *Menu {
	return m.setHandler(EditMnemonicShortcut, fn)
}
func (m *Menu) SetClearMnemonicHandler(fn func()) *Menu {
	return m.setHandler(ClearMnemonicShortcut, fn)
}
func (m *Menu) SetFindMissingWordsHandler(fn func()) *Menu {
	return m.setHandler(FindMissingWordsShortcut, fn)
}
func (m *Menu) SetFindSwappedWordsHandler(fn func()) *Menu {
	return m.setHandler(FindSwappedWordsShortcut, fn)
}
func (m *Menu) SetFindWrongWordHandler(fn func()) *Menu {
	return m.setHandler(FindWrongWordShortcut, fn)
}
func (m *Menu) SetDisplaySeedHandler(fn func()) *Menu {
	return m.setHandler(DisplaySeedShortcut, fn)
}
func (m *Menu) SetQuitHandler(fn func()) *Menu {
	return m.setHandler(QuitShortcut, fn)
}
func (m *Menu) ClearBruteForceResults() *Menu {
	m.MarkNeutral(FindMissingWordsShortcut)
	m.MarkNeutral(FindSwappedWordsShortcut)
	m.MarkNeutral(FindWrongWordShortcut)
	return m
}
func (m *Menu) MarkNeutral(shortcut rune) *Menu {
	_, item := m.GetItem(shortcut)
	if item == nil { // some menu items are conditional on dev mode
		return m
	}
	rx := regexp.MustCompile("\\[(red|green)][✓✗]\\[-] ") // TODO - use theme
	for i, s := range item.Texts {
		item.Texts[i] = rx.ReplaceAllString(s, "")
	}
	return m
}
func (m *Menu) MarkSuccess(shortcut rune) {
	_, item := m.GetItem(shortcut)
	item.SetEnabled(false)
	for i, t := range item.Texts {
		item.Texts[i] = "[green]✓[-] " + t // TODO - use theme
	}
	idx, _ := m.GetItem(EditMnemonicShortcut)
	m.SetCurrentItem(idx)
}
func (m *Menu) MarkFailure(shortcut rune) {
	_, item := m.GetItem(shortcut)
	item.SetEnabled(false)
	for i, t := range item.Texts {
		item.Texts[i] = "[red]✗[-] " + t // TODO - use theme
	}
	idx, _ := m.GetItem(EditMnemonicShortcut)
	m.SetCurrentItem(idx)
}
func (m *Menu) SetFoundMissingWords() {
	m.MarkSuccess(FindMissingWordsShortcut)
}
func (m *Menu) SetExhaustedMissingWords() {
	m.GetItem(FindMissingWordsShortcut)
}
func (m *Menu) SetFoundSwappedWords() {
	m.MarkSuccess(FindSwappedWordsShortcut)
}
func (m *Menu) SetExhaustedSwappedWords() {
	m.MarkFailure(FindSwappedWordsShortcut)
}
func (m *Menu) SetFoundWrongWord() {
	m.MarkSuccess(FindWrongWordShortcut)
}
func (m *Menu) SetExhaustedWrongWord() {
	m.MarkFailure(FindWrongWordShortcut)
}

func NewMenuView(theme Theme, devMode bool) (*Menu, *MenuCallbacks) {
	showInfo := responsive.NewResponsiveListItem().
		SetTexts("info").
		SetHelp(ShowInfoHelp).
		SetShortcut(ShowInfoShortcut).
		SetEnabledStyle(theme.ListItem.Enabled.Main).
		SetDisabledStyle(theme.ListItem.Disabled.Main)
	createRandomMnemonic := responsive.NewResponsiveListItem().
		// Do NOT use this mnemonic for real funds, only for exploring this tool!
		// The random generator actually used may or may not be cryptographically secure.
		SetTexts("random mnemonic", "random").
		SetHelp(RandomMnemonicHelp).
		SetShortcut(RandomMnemonicShortcut).
		SetEnabledStyle(theme.ListItem.Enabled.Main).
		SetDisabledStyle(theme.ListItem.Disabled.Main)
	editMnemonic := responsive.NewResponsiveListItem().
		SetTexts("edit mnemonic", "edit").
		SetHelp(EditMnemonicHelp).
		SetShortcut(EditMnemonicShortcut).
		SetEnabledStyle(theme.ListItem.Enabled.Main).
		SetDisabledStyle(theme.ListItem.Disabled.Main)
	clearMnemonic := responsive.NewResponsiveListItem().
		SetTexts("clear mnemonic", "clear").
		SetHelp(ClearMnemonicHelp).
		SetShortcut(ClearMnemonicShortcut).
		SetEnabledStyle(theme.ListItem.Enabled.Main).
		SetDisabledStyle(theme.ListItem.Disabled.Main)
	// The following menu options require a mnemonic:
	findMissingWords := responsive.NewResponsiveListItem().
		SetTexts("find missing words", "missing words", "missing").
		SetHelp(FindMissingWordsHelp).
		SetShortcut(FindMissingWordsShortcut).
		SetEnabled(false).
		SetEnabledStyle(theme.ListItem.Enabled.Main).
		SetDisabledStyle(theme.ListItem.Disabled.Main)
	findSwappedWords := responsive.NewResponsiveListItem().
		SetTexts("find swapped words", "swapped words", "swapped").
		SetHelp(FindSwappedWordsHelp).
		SetShortcut(FindSwappedWordsShortcut).
		SetEnabled(false).
		SetEnabledStyle(theme.ListItem.Enabled.Main).
		SetDisabledStyle(theme.ListItem.Disabled.Main)
	findWrongWord := responsive.NewResponsiveListItem().
		SetTexts("find wrong words", "wrong words", "wrong").
		SetHelp(FindWrongWordHelp).
		SetShortcut(FindWrongWordShortcut).
		SetEnabled(false).
		SetEnabledStyle(theme.ListItem.Enabled.Main).
		SetDisabledStyle(theme.ListItem.Disabled.Main)
	displaySeed := responsive.NewResponsiveListItem().
		SetTexts("display seed", "display").
		SetHelp(DisplaySeedHelp).
		SetShortcut(DisplaySeedShortcut).
		SetEnabledStyle(theme.ListItem.Enabled.Main).
		SetDisabledStyle(theme.ListItem.Disabled.Main)
	quit := responsive.NewResponsiveListItem().
		SetTexts("quit").
		SetHelp(QuitHelp).
		SetShortcut(QuitShortcut).
		SetEnabledStyle(theme.ListItem.Enabled.Main).
		SetDisabledStyle(theme.ListItem.Disabled.Main)

	var items []*responsive.ResponsiveListItem
	if devMode {
		items = []*responsive.ResponsiveListItem{
			createRandomMnemonic,
			editMnemonic,
			clearMnemonic,
			findMissingWords,
			findSwappedWords,
			findWrongWord,
			displaySeed,
			quit,
		}
	} else {
		items = []*responsive.ResponsiveListItem{
			showInfo,
			editMnemonic,
			clearMnemonic,
			findMissingWords,
			findWrongWord,
			quit,
		}
	}

	menu := &Menu{responsive.NewResponsiveList(items...), make(map[rune]func())}
	menu.
		// When setting MainTextStyle, the entire list's items will be styled.
		// Overriding individual items' attributes will not work as expected,
		// because they are nested inside an outer text style.
		// SetMainTextStyle(theme.ListItem.Enabled.Main).
		SetSecondaryTextStyle(theme.ListItem.Enabled.Secondary).
		SetSelectedStyle(theme.ListItem.Enabled.Selected).
		SetShortcutStyle(theme.ListItem.Enabled.Shortcut).
		SetSelectedFunc(func(index int, text string, secondaryText string, shortcut rune) {
			if fn, ok := menu.handlers[shortcut]; ok {
				item := menu.Items[menu.GetCurrentItem()]
				if !item.Enabled {
					panic("Selected disabled item")
				} else {
					fn()
				}
			}
		}).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyUp, tcell.KeyLeft, tcell.KeyBacktab:
				menu.SetPreviousEnabledItem()
				return nil
			case tcell.KeyDown, tcell.KeyRight, tcell.KeyTab:
				menu.SetNextEnabledItem()
				return nil
			case tcell.KeyRune:
				if event.Rune() == ' ' {
					return nil
				}
				for _, item := range menu.Items {
					if item.Shortcut == event.Rune() && !item.Enabled {
						return nil
					}
				}
				return event
			default:
				return event
			}
		})

	return menu, &MenuCallbacks{
		OnMnemonicValid: func() {
			findMissingWords.SetEnabled(false)
			findSwappedWords.SetEnabled(false)
			findWrongWord.SetEnabled(false)
			displaySeed.SetEnabled(true)
		},
		OnMnemonicInvalid: func() {
			findMissingWords.SetEnabled(false)
			findSwappedWords.SetEnabled(true)
			findWrongWord.SetEnabled(true)
			displaySeed.SetEnabled(false)
		},
		OnMnemonicPartial: func() {
			findMissingWords.SetEnabled(true)
			findSwappedWords.SetEnabled(false)
			findWrongWord.SetEnabled(false)
			displaySeed.SetEnabled(false)
		},
		OnMnemonicInsufficient: func() {
			findMissingWords.SetEnabled(false)
			findSwappedWords.SetEnabled(false)
			findWrongWord.SetEnabled(false)
			displaySeed.SetEnabled(false)
		},
	}
}
