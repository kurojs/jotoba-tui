package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
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

type tickMsg struct{}

type model struct {
	textInput textinput.Model
	spinner   spinner.Model
	results   []wordResult
	loading   bool
	err       error
	frame     int
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter a Japanese word..."
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
		FPS: time.Second / 10,
	}

	return model{
		textInput: ti,
		spinner:   s,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick, tickCmd())
}

func tickCmd() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
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
				m.results = nil
				m.loading = true
				m.err = nil
				return m, tea.Batch(
					m.spinner.Tick,
					searchWord(query),
				)
			}
		}
	case searchResultMsg:
		m.loading = false
		m.results = msg.results
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

var accent = lipgloss.Color("#22c55e")

func (m model) View() string {
	var b strings.Builder

	pulse := math.Sin(float64(m.frame)*0.08)*0.5 + 0.5

	b.WriteString(dotStyle.Render("●"))
	b.WriteString(" ")
	b.WriteString(titleStyle.Render("Jotoba"))
	b.WriteString(hintStyle.Render(" — Japanese Dictionary"))
	b.WriteString("\n\n")
	b.WriteString(m.textInput.View())
	b.WriteString("\n")

	if m.loading {
		b.WriteString("\n  ")
		b.WriteString(m.spinner.View())
		b.WriteString("\n")
		return b.String()
	}

	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render("  Error: " + m.err.Error()))
		b.WriteString("\n")
		return b.String()
	}

	if len(m.results) == 0 {
		b.WriteString("\n")
		b.WriteString(hintStyle.Render("  Enter a Japanese word and press Enter to search"))
		b.WriteString("\n")
	}

	for i, r := range m.results {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(fmt.Sprintf(
			"  %s %s %s",
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

	b.WriteString("  ")
	b.WriteString(glowLine(38, m.frame))
	b.WriteString("\n")
	b.WriteString(pulseIndicator(pulse))
	b.WriteString(hintStyle.Render("  Ctrl+C / Esc to quit"))

	return b.String()
}

func glowLine(length, frame int) string {
	pos := (frame / 2) % (length * 2)
	if pos >= length {
		pos = length*2 - pos - 1
	}

	var chars []string
	for i := 0; i < length; i++ {
		d := abs(i - pos)
		switch {
		case d == 0:
			chars = append(chars, glowPeak.Render("█"))
		case d == 1:
			chars = append(chars, glowMid.Render("▓"))
		case d == 2:
			chars = append(chars, glowDim.Render("▒"))
		default:
			chars = append(chars, glowBase.Render("─"))
		}
	}
	return strings.Join(chars, "")
}

func pulseIndicator(pulse float64) string {
	pulseGreen := fmt.Sprintf("#%02x%02x%02x",
		int(0x22-(0x22-0x0a)*(1-pulse)),
		int(0xc5-(0xc5-0x55)*(1-pulse)),
		int(0x5e-(0x5e-0x2d)*(1-pulse)),
	)
	dotStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(pulseGreen)).
		Bold(true)
	return dotStyle.Render("●") + "  "
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accent)

	dotStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true)

	wordStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accent)

	readingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ade80"))

	meaningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555"))

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#737373"))

	glowPeak = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#22c55e")).
			Bold(true)

	glowMid = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#16a34a"))

	glowDim = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#15803d"))

	glowBase = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#444444"))
)

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
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
