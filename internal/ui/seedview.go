package ui

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/lightningnetwork/lnd/aezeed"
	"github.com/mustwork/aezeed-recover/pkg/crack"
	"github.com/rivo/tview"
	"golang.org/x/exp/slices"
	"strings"
	"unicode"
)

// TODO - detect, whether terminal supports character range:
// check mark: '\u2713' =  (✓)
// cross mark: '\u2717' = (✗)
const (
	InvalidWords            = "✗ invalid words"
	IncompleteMnemonic      = "✗ incomplete mnemonic"
	InvalidMnemonic         = "✗ invalid mnemonic"
	ValidMnemonic           = "✓ valid mnemonic"
	PasswordInputFieldWidth = 16
)

type SeedView struct {
	*tview.Grid
	validationView   *tview.TextView
	Theme            Theme
	config           *Config
	setFocus         func(tview.Primitive) *tview.Application
	mnemonicCallback func([]string, bool)
	seedInputs       [24]*MnemonicWordInputField
	passwordInput    *PasswordInputField
}

func (v *SeedView) SetDoneFunc(fn func(tcell.Key)) *SeedView {
	for _, field := range v.seedInputs {
		field.SetDoneFunc(fn)
	}
	v.passwordInput.SetDoneFunc(fn)
	return v
}
func (v *SeedView) SetPasswordInputFocusFunc(fn func()) *SeedView {
	v.passwordInput.SetFocusFunc(fn)
	return v
}
func (v *SeedView) SetMnemonicInputFocusFunc(fn func()) *SeedView {
	for _, field := range v.seedInputs {
		field.SetFocusFunc(fn)
	}
	return v
}
func (v *SeedView) SetOnMnemonicChangedFunc(fn func(words []string, valid bool)) *SeedView {
	v.mnemonicCallback = fn
	return v
}
func (v *SeedView) SetItemText(idx int, text string) {
	v.seedInputs[idx].SetText(text)
}
func (v *SeedView) GetPassword() string {
	text := v.passwordInput.GetText()
	if len(text) == 0 {
		return v.config.DefaultPassword
	} else {
		return text
	}
}
func (v *SeedView) GetWords() (words []string) {
	for _, field := range v.seedInputs {
		word := strings.TrimSpace(field.GetText())
		if len(word) == 0 {
			continue
		}
		words = append(words, word)
	}
	return words
}
func (v *SeedView) SetWords(words []string) {
	for i, word := range words {
		v.seedInputs[i].SetText(word)
	}
	v.mnemonicChanged()
}
func (v *SeedView) GetMnemonic() (mnemonic aezeed.Mnemonic, err error) {
	words := v.GetWords()
	if len(words) != v.config.MnemonicLength {
		err = fmt.Errorf("mnemonic must have %d words", v.config.MnemonicLength)
	} else {
		mnemonic = aezeed.Mnemonic(words)
	}
	return
}
func (v *SeedView) mnemonicChanged() {
	for _, word := range v.GetWords() {
		if !slices.Contains(v.config.Wordlist, word) {
			v.validationView.
				SetTextStyle(v.Theme.TextView.Invalid.Text).
				SetText("\n " + InvalidWords)
			v.mnemonicCallback(v.GetWords(), false)
			return
		}
	}
	if len(v.GetWords()) < v.config.MnemonicLength {
		v.validationView.
			SetTextStyle(v.Theme.TextView.Invalid.Text).
			SetText("\n " + IncompleteMnemonic)
	} else {
		_, err := v.GetCipherSeed()
		if err != nil {
			v.validationView.
				SetTextStyle(v.Theme.TextView.Invalid.Text).
				SetText("\n " + InvalidMnemonic)
		} else {
			v.validationView.
				SetTextStyle(v.Theme.TextView.Valid.Text).
				SetText("\n " + ValidMnemonic)
			v.mnemonicCallback(v.GetWords(), true)
			return
		}
	}
	v.mnemonicCallback(v.GetWords(), false)
}
func (v *SeedView) GetCipherSeed() (seed *aezeed.CipherSeed, err error) {
	mnemonic, err := v.GetMnemonic()
	if err != nil {
		return
	}
	return mnemonic.ToCipherSeed([]byte(v.GetPassword()))
}
func (v *SeedView) GetSettings() crack.MnemonicSettings {
	return crack.MnemonicSettings{
		Password:       v.GetPassword(),
		Mnemonic:       v.GetWords(),
		Wordlist:       v.config.Wordlist,
		MnemonicLength: v.config.MnemonicLength,
	}
}

func WordCompleter(wordlist []string) func(string) []string {
	return func(currentText string) (entries []string) {
		if len(currentText) == 0 {
			return wordlist
		}
		entries = []string{}
		for _, word := range wordlist {
			if strings.HasPrefix(word, currentText) {
				entries = append(entries, word)
			}
		}
		return entries
	}
}

type MnemonicWordInputField struct {
	*tview.InputField
	Theme        Theme
	config       *Config
	autocomplete bool
}

func (i *MnemonicWordInputField) updateColor() *MnemonicWordInputField {
	if slices.Contains(i.config.Wordlist, i.GetText()) {
		i.SetFieldStyle(i.Theme.InputField.Valid.Field)
	} else {
		i.SetFieldStyle(i.Theme.InputField.Invalid.Field)
	}
	return i
}
func (i *MnemonicWordInputField) SetText(text string) *MnemonicWordInputField {
	i.InputField.SetText(text)
	i.updateColor()
	return i
}
func (f MnemonicWordInputField) IsEmpty() bool {
	return len(f.GetText()) == 0
}
func (i *MnemonicWordInputField) Highlight() *MnemonicWordInputField {
	i.InputField.SetFieldStyle(i.Theme.InputField.Highlight.Field)
	return i
}

type FocusItem interface {
	FocusPrev(current int)
	FocusNext(current int)
	FocusAdjacent(current int)
}

func (v *SeedView) FocusPrev(current int) {
	if current == 0 {
		v.setFocus(v.passwordInput)
	} else {
		v.setFocus(v.seedInputs[(current+23)%24])
	}
}
func (v *SeedView) FocusNext(current int) {
	if current == 23 {
		v.setFocus(v.passwordInput)
	} else {
		v.setFocus(v.seedInputs[(current+1)%24])
	}
}
func (v *SeedView) FocusAdjacent(current int) {
	v.setFocus(v.seedInputs[(current+12)%24])
}
func (v *SeedView) HighlightWords(xs []int) *SeedView {
	for i, input := range v.seedInputs {
		if slices.Contains(xs, i) {
			input.Highlight()
		}
	}
	return v
}

func NewSeedView(theme Theme, config *Config, setFocus func(p tview.Primitive) *tview.Application) *SeedView {
	// TODO - define min width, and shrink menu once min width is reached
	h := 1
	var maxWordLength int // 8 on default wordlist
	for _, word := range config.Wordlist {
		if len(word) > maxWordLength {
			maxWordLength = len(word)
		}
	}
	wordCompleter := WordCompleter(config.Wordlist)

	validationView := tview.NewTextView().
		SetDynamicColors(true)

	seedView := &SeedView{tview.NewGrid(), validationView, theme, config, setFocus, nil, [24]*MnemonicWordInputField{}, nil}

	passwordInput := NewPasswordInputField(theme, config, seedView)

	seedView.passwordInput = passwordInput
	seedView.
		// TODO - add padding by adding nil components
		SetRows(2, h, h, h, h, h, h, h, h, h, h, h, h, 0).
		SetColumns(0, 0).
		SetBorder(true)
	seedView.
		AddItem(passwordInput, 0, 0, 1, 2, 0, 0, true)

	rowOffset := 1
	for i := 0; i < 24; i++ {
		field := NewMnemonicInputField(seedView, i, wordCompleter)
		if i < 12 {
			seedView.AddItem(field, rowOffset+(i%12), 0, 1, 1, 0, 0, true)
		} else {
			seedView.AddItem(field, rowOffset+(i%12), 1, 1, 1, 0, 0, true)
		}
		seedView.seedInputs[i] = field
	}
	seedView.AddItem(validationView, rowOffset+12, 0, 1, 2, 0, 0, false)

	return seedView
}

type PasswordInputField struct {
	*tview.InputField
}

func (p *PasswordInputField) SetFieldTextColor(color tcell.Color) *PasswordInputField {
	p.InputField.SetFieldTextColor(color)
	return p
}

func NewPasswordInputField(theme Theme, config *Config, seedView *SeedView) *PasswordInputField {
	f := PasswordInputField{tview.NewInputField()} // TODO - create responsive.NewInputField
	f.SetLabelStyle(theme.InputField.Valid.Label).
		SetFieldStyle(theme.InputField.Valid.Field).
		SetAutocompleteStyles(theme.InputField.Valid.AutocompleteBackground,
			theme.InputField.Valid.AutocompleteMain,
			theme.InputField.Valid.AutocompleteSelected).
		SetLabel(" password: ").
		SetFieldWidth(PasswordInputFieldWidth).SetText(config.DefaultPassword).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyLeft, tcell.KeyUp, tcell.KeyBacktab, tcell.KeyCtrlP:
				seedView.FocusPrev(24)
				seedView.mnemonicChanged()
				return nil
			case tcell.KeyRight, tcell.KeyDown, tcell.KeyTab, tcell.KeyCtrlN, tcell.KeyEnter:
				seedView.FocusNext(-1)
				seedView.mnemonicChanged()
				return nil
			default:
				return event
			}
		})
	return &f
}

func NewMnemonicInputField(seedView *SeedView, i int, wordCompleter func(string) []string) *MnemonicWordInputField {
	theme := seedView.Theme
	config := seedView.config
	field := MnemonicWordInputField{tview.NewInputField(), theme, config, false}
	field.SetLabelStyle(theme.InputField.Valid.Label).
		SetFieldStyle(theme.InputField.Valid.Field).
		SetAutocompleteStyles(theme.InputField.Valid.AutocompleteBackground,
			theme.InputField.Valid.AutocompleteMain,
			theme.InputField.Valid.AutocompleteSelected).
		SetLabel(fmt.Sprintf(" %*d ", 2, i+1)).
		SetFieldWidth(9).
		SetAutocompleteFunc(func(currentText string) (entries []string) {
			if field.autocomplete {
				return wordCompleter(currentText)
			}
			return []string{}
		}).
		SetInputCapture(func(event *tcell.EventKey) (result *tcell.EventKey) {
			if field.autocomplete {
				switch event.Key() {
				case tcell.KeyRune:
					// This is for pasting a comma or otherwise separated mnemonic on macOS
					switch event.Rune() {
					case ' ', ';', ',':
						field.updateColor()
						seedView.FocusNext(i)
						return nil
					default:
						return tcell.NewEventKey(event.Key(), unicode.ToLower(event.Rune()), event.Modifiers())
					}
				case tcell.KeyEscape:
					field.autocomplete = false
					field.updateColor()
					seedView.mnemonicChanged()
					result = event
				case tcell.KeyDEL:
					if field.IsEmpty() {
						field.autocomplete = false
						field.updateColor()
						seedView.mnemonicChanged()
						seedView.FocusPrev(i)
						result = nil
					}
					result = event
				case tcell.KeyCtrlP:
					// navigating within autocomplete drop down
					return tcell.NewEventKey(tcell.KeyUp, event.Rune(), event.Modifiers())
				case tcell.KeyCtrlN:
					// navigating within autocomplete drop down
					return tcell.NewEventKey(tcell.KeyDown, event.Rune(), event.Modifiers())
				case tcell.KeyEnter:
					// TODO - maybe we need a tri-state for autocomplete?
					//  complete on enter does not work initially
					//  but works on the second focus.
					//  Also, enter does not show dropdown, immediately but requires a second keypress.
					completion := wordCompleter(field.GetText())
					if len(completion) > 0 {
						field.SetText(completion[0])
					}
					field.autocomplete = false
					// Because enter does not change the fields context,
					// the mnemonicChanged method sees the actual word,
					// even though this event has not been propagated, yet.
					field.updateColor()
					seedView.mnemonicChanged()
					seedView.FocusNext(i)
					result = event
				case tcell.KeyCtrlV:
					_, err := clipboard.ReadAll()
					if err != nil {
						return
					}
					// TODO - implement and test on linux!!!
				default:
					result = event
				}
				return
			}
			if unicode.IsLetter(event.Rune()) {
				field.autocomplete = true
				return tcell.NewEventKey(event.Key(), unicode.ToLower(event.Rune()), event.Modifiers())
			}
			switch event.Key() {
			// we are a bit defensive here with calling #mnemonicChanged()
			case tcell.KeyEscape, tcell.KeyDelete:
				seedView.mnemonicChanged()
				return event
			case tcell.KeyLeft, tcell.KeyRight:
				seedView.mnemonicChanged()
				seedView.FocusAdjacent(i)
			case tcell.KeyBacktab, tcell.KeyUp, tcell.KeyCtrlP:
				seedView.mnemonicChanged()
				seedView.FocusPrev(i)
			case tcell.KeyDown, tcell.KeyTab, tcell.KeyCtrlN:
				seedView.mnemonicChanged()
				seedView.FocusNext(i)
			case tcell.KeyDEL:
				if field.IsEmpty() {
					seedView.mnemonicChanged()
					seedView.FocusPrev(i)
					return nil
				} else if (event.Modifiers() & tcell.ModAlt) == tcell.ModAlt { // this does not seem to work on macOS
					field.SetText("")
					return nil
				} else {
					field.autocomplete = true
					return event
				}
			case tcell.KeyEnter:
				if field.IsEmpty() {
					field.autocomplete = true
					field.Autocomplete()
					return nil
				} else {
					seedView.mnemonicChanged()
					seedView.FocusNext(i)
					return nil
				}
			case tcell.KeyRune:
				if event.Rune() == ' ' {
					field.autocomplete = true
					field.Autocomplete()
					return nil
				}
			}
			// It does not seem to make a difference, whether we return nil or not.
			return event
		})
	return &field
}
