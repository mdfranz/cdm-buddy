package ui

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// NewForm applies the shared terminal theme used across wizard and editor flows.
func NewForm(groups ...*huh.Group) *huh.Form {
	return huh.NewForm(groups...).WithTheme(formTheme())
}

func formTheme() *huh.Theme {
	t := huh.ThemeBase()

	var (
		text   = lipgloss.Color("252")
		muted  = lipgloss.Color("245")
		accent = lipgloss.Color("110")
		warm   = lipgloss.Color("215")
		good   = lipgloss.Color("79")
		danger = lipgloss.Color("203")
		strong = lipgloss.Color("24")
		subtle = lipgloss.Color("236")
	)

	t.FieldSeparator = lipgloss.NewStyle().SetString("\n")

	t.Focused.Base = lipgloss.NewStyle().
		PaddingLeft(2)
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = lipgloss.NewStyle().Bold(true).Foreground(accent)
	t.Focused.NoteTitle = t.Focused.Title
	t.Focused.Description = lipgloss.NewStyle().Foreground(muted)
	t.Focused.ErrorIndicator = lipgloss.NewStyle().Foreground(danger).SetString(" •")
	t.Focused.ErrorMessage = lipgloss.NewStyle().Foreground(danger)
	t.Focused.SelectSelector = lipgloss.NewStyle().Foreground(warm).SetString("› ")
	t.Focused.Option = lipgloss.NewStyle().Foreground(text)
	t.Focused.NextIndicator = lipgloss.NewStyle().Foreground(warm).SetString("›")
	t.Focused.PrevIndicator = lipgloss.NewStyle().Foreground(warm).SetString("‹")
	t.Focused.MultiSelectSelector = lipgloss.NewStyle().Foreground(warm).SetString("› ")
	t.Focused.SelectedOption = lipgloss.NewStyle().Foreground(good)
	t.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(good).SetString("● ")
	t.Focused.UnselectedOption = lipgloss.NewStyle().Foreground(text)
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(muted).SetString("○ ")
	t.Focused.FocusedButton = lipgloss.NewStyle().
		Padding(0, 1).
		MarginRight(1).
		Foreground(lipgloss.Color("230")).
		Background(strong).
		Bold(true)
	t.Focused.BlurredButton = lipgloss.NewStyle().
		Padding(0, 1).
		MarginRight(1).
		Foreground(text).
		Background(subtle)
	t.Focused.Next = t.Focused.FocusedButton
	t.Focused.TextInput.Cursor = lipgloss.NewStyle().Foreground(warm)
	t.Focused.TextInput.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)
	t.Focused.TextInput.Prompt = lipgloss.NewStyle().Foreground(warm)
	t.Focused.TextInput.Text = lipgloss.NewStyle().Foreground(text)

	t.Blurred = t.Focused
	t.Blurred.Base = lipgloss.NewStyle().PaddingLeft(2)
	t.Blurred.Card = t.Blurred.Base
	t.Blurred.Title = lipgloss.NewStyle().Foreground(muted)
	t.Blurred.NoteTitle = t.Blurred.Title
	t.Blurred.Description = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	t.Blurred.SelectSelector = lipgloss.NewStyle().SetString("  ")
	t.Blurred.MultiSelectSelector = lipgloss.NewStyle().SetString("  ")
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()
	t.Blurred.TextInput.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	t.Blurred.TextInput.Text = lipgloss.NewStyle().Foreground(text)

	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description

	return t
}
