package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type wordResult struct {
	Word     string
	Reading  string
	Meanings []string
}

type searchResultMsg struct {
	results []wordResult
}

type errorMsg struct {
	err error
}

type model struct {
	textInput textinput.Model
	results   []wordResult
	loading   bool
	err       error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter a Japanese word..."
	ti.Focus()
	ti.CharLimit = 60
	ti.Width = 50

	return model{
		textInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			query := strings.TrimSpace(m.textInput.Value())
			if query != "" && !m.loading {
				m.textInput.SetValue("")
				m.loading = true
				m.err = nil
				return m, searchWord(query)
			}
		}
	case searchResultMsg:
		m.loading = false
		m.results = msg.results
		return m, nil
	case errorMsg:
		m.loading = false
		m.err = msg.err
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFD700"))

	wordStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFD700"))

	readingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#87FF87"))

	meaningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555"))

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	lineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#585858"))
)

func (m model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Jotoba — Japanese Dictionary"))
	b.WriteString("\n\n")
	b.WriteString(m.textInput.View())
	b.WriteString("\n")

	if m.loading {
		b.WriteString("\n  Searching...")
		return b.String()
	}

	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render("  Error: " + m.err.Error()))
		return b.String()
	}

	if len(m.results) == 0 {
		b.WriteString("\n")
		b.WriteString(hintStyle.Render("  Enter a Japanese word and press Enter to search"))
	}

	for _, r := range m.results {
		b.WriteString("\n")
		b.WriteString(wordStyle.Render(r.Word))
		b.WriteString("\n")
		b.WriteString(lineStyle.Render(strings.Repeat("━", 28)))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf(
			"%s %s %s",
			wordStyle.Render(r.Word),
			hintStyle.Render("→"),
			readingStyle.Render(r.Reading),
		))
		b.WriteString("\n")
		for _, m := range r.Meanings {
			b.WriteString(meaningStyle.Render("  - " + m))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(hintStyle.Render("  Ctrl+C / Esc to quit"))

	return b.String()
}

type jotobaWord struct {
	Reading struct {
		Kana  string `json:"kana"`
		Kanji string `json:"kanji"`
	} `json:"reading"`
	Senses []struct {
		Glosses  []string `json:"glosses"`
		Language string   `json:"language"`
	} `json:"senses"`
}

type jotobaResponse struct {
	Words []jotobaWord `json:"words"`
}

func searchWord(query string) tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{Timeout: 10 * time.Second}

		payload := map[string]interface{}{
			"query":    query,
			"language": "Spanish",
		}

		jsonBody, err := json.Marshal(payload)
		if err != nil {
			return errorMsg{err}
		}

		req, err := http.NewRequest(
			"POST",
			"https://jotoba.de/api/search/words",
			strings.NewReader(string(jsonBody)),
		)
		if err != nil {
			return errorMsg{err}
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return errorMsg{err}
		}
		defer resp.Body.Close()

		var jotobaResp jotobaResponse
		if err := json.NewDecoder(resp.Body).Decode(&jotobaResp); err != nil {
			return errorMsg{err}
		}

		if len(jotobaResp.Words) == 0 {
			return searchResultMsg{[]wordResult{{
				Word:     query,
				Reading:  "—",
				Meanings: []string{"No results found"},
			}}}
		}

		var results []wordResult
		for _, w := range jotobaResp.Words {
			display := w.Reading.Kanji
			if display == "" {
				display = w.Reading.Kana
			}
			r := wordResult{
				Word:    display,
				Reading: w.Reading.Kana,
			}

			var esGlosses, otherGlosses []string
			for _, s := range w.Senses {
				if len(s.Glosses) == 0 {
					continue
				}
				if s.Language == "Spanish" {
					esGlosses = append(esGlosses, strings.Join(s.Glosses, ", "))
				} else {
					otherGlosses = append(otherGlosses, strings.Join(s.Glosses, ", "))
				}
			}

			if len(esGlosses) > 0 {
				r.Meanings = esGlosses
			} else {
				r.Meanings = otherGlosses
			}

			results = append(results, r)
		}

		return searchResultMsg{results}
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
