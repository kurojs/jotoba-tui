package tui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"kurojs.com/jotoba-tui/internal/jotoba"
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyTab:
			m.mode = (m.mode + 1) % 3
			m.wordResults = nil
			m.kanjiResults = nil
			m.sentenceResults = nil
			m.err = nil
			m.textInput.SetValue("")
			return m, tickCmd()
		case tea.KeyCtrlL:
			m.langIndex = (m.langIndex + 1) % len(languages)
			m.wordResults = nil
			m.kanjiResults = nil
			m.sentenceResults = nil
			m.err = nil
			m.textInput.SetValue("")
			return m, tickCmd()
		case tea.KeyEnter:
			query := strings.TrimSpace(m.textInput.Value())
			if query != "" && !m.loading {
				m.textInput.SetValue("")
				m.wordResults = nil
				m.kanjiResults = nil
				m.sentenceResults = nil
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
		return m, tickCmd()
	case errorMsg:
		m.loading = false
		m.err = msg.err
		return m, tickCmd()
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
