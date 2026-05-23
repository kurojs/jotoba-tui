package tui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/kurojs/jotoba-tui/internal/config"
	"github.com/kurojs/jotoba-tui/internal/jotoba"
)

func tickCmd() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick, tickCmd())
}

func searchCmd(query string, mode searchMode, language string) tea.Cmd {
	return func() tea.Msg {
		switch mode {
		case modeWord:
			results, err := jotoba.SearchWords(query, language)
			if err != nil {
				return errorMsg{err}
			}
			return searchResultMsg{mode: mode, results: results}
		case modeKanji:
			results, err := jotoba.SearchKanji(query, language)
			if err != nil {
				return errorMsg{err}
			}
			return searchResultMsg{mode: mode, results: results}
		case modeSentence:
			results, err := jotoba.SearchSentences(query, language)
			if err != nil {
				return errorMsg{err}
			}
			return searchResultMsg{mode: mode, results: results}
		default:
			return errorMsg{err: ErrUnknownMode}
		}
	}
}

var ErrUnknownMode = &modeError{"unknown search mode"}

type modeError struct{ msg string }

func (e *modeError) Error() string { return e.msg }

func hasResults(m model) bool {
	return len(m.wordResults) > 0 || len(m.kanjiResults) > 0 || len(m.sentenceResults) > 0
}

func maxScroll(m model) int {
	switch m.mode {
	case modeWord:
		if len(m.wordResults) == 0 {
			return 0
		}
		return len(m.wordResults) - 1
	case modeKanji:
		if len(m.kanjiResults) == 0 {
			return 0
		}
		return len(m.kanjiResults) - 1
	case modeSentence:
		if len(m.sentenceResults) == 0 {
			return 0
		}
		return len(m.sentenceResults) - 1
	default:
		return 0
	}
}

func resultCount(m model) int {
	switch m.mode {
	case modeWord:
		return len(m.wordResults)
	case modeKanji:
		return len(m.kanjiResults)
	case modeSentence:
		return len(m.sentenceResults)
	default:
		return 0
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termHeight = msg.Height

	case tea.KeyMsg:
		if m.showLangMenu {
			return m.updateLangMenu(msg)
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyTab:
			m.mode = (m.mode + 1) % 3
			m.wordResults = nil
			m.kanjiResults = nil
			m.sentenceResults = nil
			m.err = nil
			m.scrollOffset = 0
			m.textInput.SetValue("")
		case tea.KeyCtrlL:
			m.showLangMenu = true
			m.langCursor = m.langIndex
		case tea.KeyUp:
			if hasResults(m) && m.scrollOffset > 0 {
				m.scrollOffset--
			}
		case tea.KeyDown:
			if hasResults(m) && m.scrollOffset < maxScroll(m) {
				m.scrollOffset++
			}
		case tea.KeyEnter:
			query := strings.TrimSpace(m.textInput.Value())
			if query != "" && !m.loading {
				m.textInput.SetValue("")
				m.wordResults = nil
				m.kanjiResults = nil
				m.sentenceResults = nil
				m.scrollOffset = 0
				m.loading = true
				m.err = nil
				return m, tea.Batch(
					m.spinner.Tick,
					searchCmd(query, m.mode, languages[m.langIndex]),
				)
			}
		}
	case searchResultMsg:
		m.loading = false
		m.scrollOffset = 0
		switch msg.mode {
		case modeWord:
			if r, ok := msg.results.([]jotoba.WordResult); ok {
				m.wordResults = r
			}
		case modeKanji:
			if r, ok := msg.results.([]jotoba.KanjiResult); ok {
				m.kanjiResults = r
			}
		case modeSentence:
			if r, ok := msg.results.([]jotoba.SentenceResult); ok {
				m.sentenceResults = r
			}
		}
	case errorMsg:
		m.loading = false
		m.err = msg.err
	case tickMsg:
		m.frame++
		return m, tickCmd()
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	var spinnerCmd tea.Cmd
	m.spinner, spinnerCmd = m.spinner.Update(msg)

	return m, tea.Batch(cmd, spinnerCmd)
}

func (m model) updateLangMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.showLangMenu = false
	case tea.KeyEnter:
		m.langIndex = m.langCursor
		m.showLangMenu = false
		m.wordResults = nil
		m.kanjiResults = nil
		m.sentenceResults = nil
		m.err = nil
		m.scrollOffset = 0
		m.textInput.SetValue("")
		config.Save(&config.Config{Language: languages[m.langIndex]})
	case tea.KeyUp:
		if m.langCursor > 0 {
			m.langCursor--
		}
	case tea.KeyDown:
		if m.langCursor < len(languages)-1 {
			m.langCursor++
		}
	}
	return m, nil
}
