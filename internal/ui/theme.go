package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Theme struct {
	tview.Theme
	Button     ButtonStyle
	ListItem   ListItemStyles
	InputField InputFieldStyles
	TextView   TextViewStyles
}

type ButtonStyle struct {
	Style     tcell.Style
	Activated tcell.Style
	Disabled  tcell.Style
}

func (s ButtonStyle) SetForeground(c tcell.Color) ButtonStyle {
	s.Style = s.Style.Foreground(c)
	return s
}
func (s ButtonStyle) SetBackground(c tcell.Color) ButtonStyle {
	s.Style = s.Style.Background(c)
	return s
}
func (s ButtonStyle) SetBold(b bool) ButtonStyle {
	s.Style = s.Style.Bold(b)
	return s
}
func (s ButtonStyle) SetActivatedForeground(c tcell.Color) ButtonStyle {
	s.Activated = s.Activated.Foreground(c)
	return s
}
func (s ButtonStyle) SetActivatedBackground(c tcell.Color) ButtonStyle {
	s.Activated = s.Activated.Background(c)
	return s
}
func (s ButtonStyle) SetActivatedBold(b bool) ButtonStyle {
	s.Activated = s.Activated.Bold(b)
	return s
}
func (s ButtonStyle) SetDisabledForeground(c tcell.Color) ButtonStyle {
	s.Disabled = s.Disabled.Foreground(c)
	return s
}
func (s ButtonStyle) SetDisabledBackground(c tcell.Color) ButtonStyle {
	s.Disabled = s.Disabled.Background(c)
	return s
}
func (s ButtonStyle) SetDisabledBold(b bool) ButtonStyle {
	s.Disabled = s.Disabled.Bold(b)
	return s
}

func NewButtonStyle(theme tview.Theme) ButtonStyle {
	return ButtonStyle{
		// as in tview/button.go
		Style:     tcell.StyleDefault.Background(theme.ContrastBackgroundColor).Foreground(theme.PrimaryTextColor),
		Activated: tcell.StyleDefault.Background(theme.PrimaryTextColor).Foreground(theme.InverseTextColor),
		Disabled:  tcell.StyleDefault.Background(theme.ContrastBackgroundColor).Foreground(theme.ContrastSecondaryTextColor),
	}
}

type ListItemStyles struct {
	Enabled  ListItemStyle
	Disabled ListItemStyle
}
type ListItemStyle struct {
	Main      tcell.Style
	Secondary tcell.Style
	Shortcut  tcell.Style
	Selected  tcell.Style
}

func (s ListItemStyle) SetMainForeground(c tcell.Color) ListItemStyle {
	s.Main = s.Main.Foreground(c)
	return s
}
func (s ListItemStyle) SetSecondaryForeground(c tcell.Color) ListItemStyle {
	s.Secondary = s.Secondary.Foreground(c)
	return s
}
func (s ListItemStyle) SetShortcutForeground(c tcell.Color) ListItemStyle {
	s.Shortcut = s.Shortcut.Foreground(c)
	return s
}
func (s ListItemStyle) SetSelectedForeground(c tcell.Color) ListItemStyle {
	s.Selected = s.Selected.Foreground(c)
	return s
}
func (s ListItemStyle) SetMainBackground(c tcell.Color) ListItemStyle {
	s.Main = s.Main.Background(c)
	return s
}
func (s ListItemStyle) SetSecondaryBackground(c tcell.Color) ListItemStyle {
	s.Secondary = s.Secondary.Background(c)
	return s
}
func (s ListItemStyle) SetShortcutBackground(c tcell.Color) ListItemStyle {
	s.Shortcut = s.Shortcut.Background(c)
	return s
}
func (s ListItemStyle) SetSelectedBackground(c tcell.Color) ListItemStyle {
	s.Selected = s.Selected.Background(c)
	return s
}
func (s ListItemStyle) SetMainBold(b bool) ListItemStyle {
	s.Main = s.Main.Bold(b)
	return s
}
func (s ListItemStyle) SetSecondaryBold(b bool) ListItemStyle {
	s.Secondary = s.Secondary.Bold(b)
	return s
}
func (s ListItemStyle) SetShortcutBold(b bool) ListItemStyle {
	s.Shortcut = s.Shortcut.Bold(b)
	return s
}
func (s ListItemStyle) SetSelectedBold(b bool) ListItemStyle {
	s.Selected = s.Selected.Bold(b)
	return s
}

func NewListItemStyle(theme tview.Theme) ListItemStyle {
	return ListItemStyle{
		// as in tview/list.go
		Main:      tcell.StyleDefault.Foreground(theme.PrimaryTextColor),
		Secondary: tcell.StyleDefault.Foreground(theme.TertiaryTextColor),
		Shortcut:  tcell.StyleDefault.Foreground(theme.SecondaryTextColor),
		Selected:  tcell.StyleDefault.Foreground(theme.PrimitiveBackgroundColor).Background(theme.PrimaryTextColor),
	}
}

type InputFieldStyles struct {
	Valid     InputFieldStyle
	Invalid   InputFieldStyle
	Highlight InputFieldStyle
}

type InputFieldStyle struct {
	Label                  tcell.Style
	Field                  tcell.Style
	Placeholder            tcell.Style
	AutocompleteMain       tcell.Style
	AutocompleteSelected   tcell.Style
	AutocompleteBackground tcell.Color
}

func (s InputFieldStyle) SetLabelForeground(c tcell.Color) InputFieldStyle {
	s.Label = s.Label.Foreground(c)
	return s
}
func (s InputFieldStyle) SetLabelBackground(c tcell.Color) InputFieldStyle {
	s.Label = s.Label.Background(c)
	return s
}
func (s InputFieldStyle) SetLabelBold(b bool) InputFieldStyle {
	s.Label = s.Label.Bold(b)
	return s
}
func (s InputFieldStyle) SetFieldForeground(c tcell.Color) InputFieldStyle {
	s.Field = s.Field.Foreground(c)
	return s
}
func (s InputFieldStyle) SetFieldBackground(c tcell.Color) InputFieldStyle {
	s.Field = s.Field.Background(c)
	return s
}
func (s InputFieldStyle) SetFieldBold(b bool) InputFieldStyle {
	s.Field = s.Field.Bold(b)
	return s
}
func (s InputFieldStyle) SetPlaceholderForeground(c tcell.Color) InputFieldStyle {
	s.Placeholder = s.Placeholder.Foreground(c)
	return s
}
func (s InputFieldStyle) SetPlaceholderBackground(c tcell.Color) InputFieldStyle {
	s.Placeholder = s.Placeholder.Background(c)
	return s
}
func (s InputFieldStyle) SetPlaceholderBold(b bool) InputFieldStyle {
	s.Placeholder = s.Placeholder.Bold(b)
	return s
}
func (s InputFieldStyle) SetAutocompleteMainForeground(c tcell.Color) InputFieldStyle {
	s.AutocompleteMain = s.AutocompleteMain.Foreground(c)
	return s
}
func (s InputFieldStyle) SetAutocompleteMainBackground(c tcell.Color) InputFieldStyle {
	s.AutocompleteMain = s.AutocompleteMain.Background(c)
	return s
}
func (s InputFieldStyle) SetAutocompleteMainBold(b bool) InputFieldStyle {
	s.AutocompleteMain = s.AutocompleteMain.Bold(b)
	return s
}
func (s InputFieldStyle) SetAutocompleteSelectedForeground(c tcell.Color) InputFieldStyle {
	s.AutocompleteSelected = s.AutocompleteSelected.Foreground(c)
	return s
}
func (s InputFieldStyle) SetAutocompleteSelectedBackground(c tcell.Color) InputFieldStyle {
	s.AutocompleteSelected = s.AutocompleteSelected.Background(c)
	return s
}
func (s InputFieldStyle) SetAutocompleteSelectedBold(b bool) InputFieldStyle {
	s.AutocompleteSelected = s.AutocompleteSelected.Bold(b)
	return s
}
func (s InputFieldStyle) SetAutocompleteBackground(c tcell.Color) InputFieldStyle {
	s.AutocompleteBackground = c
	return s
}

func NewInputFieldStyle(theme tview.Theme) InputFieldStyle {
	return InputFieldStyle{
		// as in tview/inputfield.go
		Label:                  tcell.StyleDefault.Foreground(theme.SecondaryTextColor),
		Field:                  tcell.StyleDefault.Background(theme.ContrastBackgroundColor).Foreground(theme.PrimaryTextColor),
		Placeholder:            tcell.StyleDefault.Background(theme.ContrastBackgroundColor).Foreground(theme.ContrastSecondaryTextColor),
		AutocompleteMain:       tcell.StyleDefault.Foreground(theme.PrimitiveBackgroundColor),
		AutocompleteSelected:   tcell.StyleDefault.Background(theme.PrimaryTextColor).Foreground(theme.PrimitiveBackgroundColor),
		AutocompleteBackground: theme.MoreContrastBackgroundColor,
	}
}

type TextViewStyles struct {
	Valid   TextViewStyle
	Invalid TextViewStyle
}

type TextViewStyle struct {
	Label tcell.Style
	Text  tcell.Style
}

func (s TextViewStyle) SetLabelForeground(c tcell.Color) TextViewStyle {
	s.Label = s.Label.Foreground(c)
	return s
}
func (s TextViewStyle) SetLabelBackground(c tcell.Color) TextViewStyle {
	s.Label = s.Label.Background(c)
	return s
}
func (s TextViewStyle) SetLabelBold(b bool) TextViewStyle {
	s.Label = s.Label.Bold(b)
	return s
}
func (s TextViewStyle) SetTextForeground(c tcell.Color) TextViewStyle {
	s.Text = s.Text.Foreground(c)
	return s
}
func (s TextViewStyle) SetTextBackground(c tcell.Color) TextViewStyle {
	s.Text = s.Text.Background(c)
	return s
}
func (s TextViewStyle) SetTextBold(b bool) TextViewStyle {
	s.Text = s.Text.Bold(b)
	return s
}

func NewTextViewStyle(theme tview.Theme) TextViewStyle {
	return TextViewStyle{
		// as in tview/textview.go
		Label: tcell.StyleDefault.Foreground(theme.SecondaryTextColor),
		Text:  tcell.StyleDefault.Background(theme.PrimitiveBackgroundColor).Foreground(theme.PrimaryTextColor),
	}
}

// UniversalTheme returns a theme that works well with both light and dark terminal settings.
func UniversalTheme() Theme {
	theme := tview.Theme{
		PrimitiveBackgroundColor:    tcell.ColorDefault,   // main bg (`tcell.ColorWhite` is darkened for some reason)
		ContrastBackgroundColor:     tcell.ColorLightCyan, // input field text bg
		MoreContrastBackgroundColor: tcell.ColorLightBlue, // dropdown bg
		BorderColor:                 tcell.ColorGray,      // all borders
		//TitleColor:                  tcell.ColorBlack,       // unused, box titles
		//GraphicsColor:               tcell.ColorRed,         // unused
		PrimaryTextColor:   tcell.ColorDarkBlue,  // normal text fg, dropdown selection bg
		SecondaryTextColor: tcell.ColorDarkGreen, // input field label fg
		//TertiaryTextColor:           tcell.ColorCoral,       // unused, list item secondary text fg
		//InverseTextColor:            tcell.ColorSalmon,      // unused
		//ContrastSecondaryTextColor:  tcell.ColorBrown,       // unused, placeholder text fg
	}
	return Theme{
		Theme:  theme,
		Button: NewButtonStyle(theme),
		ListItem: ListItemStyles{
			Enabled: NewListItemStyle(theme).
				SetMainBold(true).
				SetSelectedBold(true).
				SetSelectedBackground(tcell.ColorLightGreen),
			Disabled: NewListItemStyle(theme).
				SetMainForeground(tcell.ColorGray),
		},
		InputField: InputFieldStyles{
			Valid: NewInputFieldStyle(theme).
				SetFieldForeground(tcell.ColorGreen).
				SetFieldBold(true).
				SetAutocompleteMainForeground(tcell.ColorPurple).
				SetAutocompleteSelectedForeground(tcell.ColorPurple).
				SetAutocompleteSelectedBold(true).
				SetAutocompleteSelectedBackground(tcell.ColorLightYellow),
			Invalid: NewInputFieldStyle(theme).
				SetFieldForeground(tcell.ColorRed).
				SetFieldBold(true),
			Highlight: NewInputFieldStyle(theme).
				SetFieldForeground(tcell.ColorBlack).
				SetFieldBackground(tcell.ColorLightSeaGreen).
				SetFieldBold(true),
		},
		TextView: TextViewStyles{
			Valid: NewTextViewStyle(theme).
				SetTextForeground(tcell.ColorLightGreen).
				SetTextBold(true),
			Invalid: NewTextViewStyle(theme).
				SetTextForeground(tcell.ColorRed).
				SetTextBold(true),
		},
	}
}
