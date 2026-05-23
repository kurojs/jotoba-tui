package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"kurojs.com/jotoba-tui/internal/jotoba"
)

type searchMode int

const (
	modeWord searchMode = iota
	modeKanji
	modeSentence
)

type searchResultMsg struct {
	mode    searchMode
	results any
}

type errorMsg struct {
	err error
}

type tickMsg struct{}

var languages = []string{
	"English",
	"German",
	"Spanish",
	"French",
	"Russian",
	"Swedish",
	"Dutch",
	"Hungarian",
	"Slovenian",
}

type model struct {
	textInput   textinput.Model
	spinner     spinner.Model
	mode        searchMode
	langIndex   int
	showLangMenu bool
	langCursor  int
	wordResults []jotoba.WordResult
	kanjiResults []jotoba.KanjiResult
	sentenceResults []jotoba.SentenceResult
	loading     bool
	err         error
	frame       int
}

var (
	accent     = lipgloss.Color("#22c55e")
	darkGreen  = lipgloss.Color("#16a34a")
	darker     = lipgloss.Color("#15803d")
	darkest    = lipgloss.Color("#444444")

	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(accent)
	dotStyle   = lipgloss.NewStyle().Foreground(accent).Bold(true)
	wordStyle  = lipgloss.NewStyle().Bold(true).Foreground(accent)
	readingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4ade80"))
	meaningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555"))
	hintStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#737373"))
	kanjiCharStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffd700"))
	kanjiInfoStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a0a0"))
	sentenceStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#e0e0e0"))
	translationStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#94a3b8"))
	tabStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#737373"))
	tabActiveStyle  = lipgloss.NewStyle().Bold(true).Foreground(accent)

	glowPeak = lipgloss.NewStyle().Bold(true).Foreground(accent)
	glowMid  = lipgloss.NewStyle().Foreground(darkGreen)
	glowDim  = lipgloss.NewStyle().Foreground(darker)
	glowBase = lipgloss.NewStyle().Foreground(darkest)

	menuTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(accent)
	menuItemStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#e0e0e0"))
	menuCursorStyle = lipgloss.NewStyle().Foreground(accent).Bold(true)
	menuActiveStyle = lipgloss.NewStyle().Bold(true).Foreground(accent)
	menuDimStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#737373"))
)

func New() tea.Model {
	ti := textinput.New()
	ti.Placeholder = "Enter Japanese text..."
	ti.Focus()
	ti.CharLimit = 60
	ti.Width = 50

	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(accent)
	s.Spinner = spinner.Spinner{
		Frames: []string{
			"▱▱▱▱▱▱▱▱",
			"▰▱▱▱▱▱▱▱",
			"▰▰▱▱▱▱▱▱",
			"▰▰▰▱▱▱▱▱",
			"▰▰▰▰▱▱▱▱",
			"▰▰▰▰▰▱▱▱",
			"▰▰▰▰▰▰▱▱",
			"▰▰▰▰▰▰▰▱",
			"▰▰▰▰▰▰▰▰",
		},
	}

	return model{
		textInput: ti,
		spinner:   s,
		langIndex: 0,
	}
}
