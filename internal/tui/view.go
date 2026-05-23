package tui

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"kurojs.com/jotoba-tui/internal/jotoba"
)

const headerLines = 8

func modeLabel(m searchMode) string {
	switch m {
	case modeWord:
		return "Words"
	case modeKanji:
		return "Kanji"
	case modeSentence:
		return "Sentences"
	default:
		return "?"
	}
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(dotStyle.Render("●"))
	b.WriteString(" ")
	b.WriteString(titleStyle.Render("Jotoba"))
	b.WriteString(hintStyle.Render(" — Japanese Dictionary"))
	b.WriteString("\n\n")

	b.WriteString("  ")
	for i := range 3 {
		if searchMode(i) == m.mode {
			b.WriteString("[" + tabActiveStyle.Render(modeLabel(searchMode(i))) + "]")
		} else {
			b.WriteString(" " + tabStyle.Render(modeLabel(searchMode(i))) + " ")
		}
		b.WriteString("  ")
	}
	b.WriteString(hintStyle.Render("(Tab to switch)"))
	b.WriteString("\n")

	b.WriteString("  ")
	b.WriteString(tabStyle.Render(languages[m.langIndex]))
	b.WriteString(hintStyle.Render("  (Ctrl+L)"))
	b.WriteString("\n\n")

	if m.showLangMenu {
		m.renderLangMenu(&b)
		return b.String()
	}

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
		b.WriteString("\n\n")
		b.WriteString(hintStyle.Render("  Press Esc to return"))
		return b.String()
	}

	availLines := m.termHeight - headerLines

	switch m.mode {
	case modeWord:
		m.renderWordResults(&b, availLines)
	case modeKanji:
		m.renderKanjiResults(&b, availLines)
	case modeSentence:
		m.renderSentenceResults(&b, availLines)
	}

	b.WriteString("\n  ")
	b.WriteString(glowLine(38, m.frame))
	b.WriteString("\n")
	b.WriteString(pulseIndicator(pulseValue(m.frame)))
	b.WriteString(hintStyle.Render("  Ctrl+C / Esc to quit"))

	return b.String()
}

func (m model) renderLangMenu(b *strings.Builder) {
	b.WriteString(menuTitleStyle.Render("  Language"))
	b.WriteString("\n\n")

	for i, lang := range languages {
		cursor := "  "
		if i == m.langCursor {
			cursor = menuCursorStyle.Render(" >")
		}

		label := menuItemStyle.Render(lang)
		if i == m.langIndex {
			label = menuActiveStyle.Render(lang + "  (current)")
		}

		b.WriteString(fmt.Sprintf("  %s  %s", cursor, label))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(menuDimStyle.Render("  Up/Down  select   Enter  confirm   Esc  cancel"))
}

func scrollInfo(total, from, to int) string {
	return fmt.Sprintf("  %s", hintStyle.Render(fmt.Sprintf("(%d-%d/%d  Up/Down)", from+1, to+1, total)))
}

func resultLines(r any) int {
	switch v := r.(type) {
	case jotoba.WordResult:
		return 1 + len(v.Meanings)
	case jotoba.KanjiResult:
		l := 1
		if len(v.Kunyomi) > 0 {
			l++
		}
		if len(v.Onyomi) > 0 {
			l++
		}
		if v.Strokes > 0 || v.Grade > 0 {
			l++
		}
		return l
	case jotoba.SentenceResult:
		l := 1
		if v.Furigana != "" {
			l++
		}
		l++
		return l
	default:
		return 2
	}
}

func (m model) renderWordResults(b *strings.Builder, maxLines int) {
	total := len(m.wordResults)
	if total == 0 {
		b.WriteString("\n")
		b.WriteString(hintStyle.Render("  Enter a Japanese word and press Enter to search"))
		b.WriteString("\n")
		return
	}

	offsets := make([]int, total)
	lineCount := 1
	for i, r := range m.wordResults {
		offsets[i] = lineCount
		lineCount += 1 + len(r.Meanings)
	}

	end := total
	for i := m.scrollOffset; i < total; i++ {
		itemEnd := offsets[i] + 1 + len(m.wordResults[i].Meanings) - offsets[i]
		if itemEnd > maxLines {
			end = i
			break
		}
	}

	visible := m.wordResults[m.scrollOffset:end]

	b.WriteString(scrollInfo(total, m.scrollOffset, m.scrollOffset+len(visible)-1))
	b.WriteString("\n")

	for _, r := range visible {
		showReading := r.Reading != "" && r.Reading != r.Word
		if showReading {
			b.WriteString(fmt.Sprintf(
				"  %s %s %s",
				wordStyle.Render(r.Word),
				hintStyle.Render("→"),
				readingStyle.Render(r.Reading),
			))
		} else {
			b.WriteString("  " + wordStyle.Render(r.Word))
		}
		b.WriteString("\n")
		for _, meaning := range r.Meanings {
			b.WriteString(meaningStyle.Render("  - " + meaning))
			b.WriteString("\n")
		}
	}
}

func (m model) renderKanjiResults(b *strings.Builder, maxLines int) {
	total := len(m.kanjiResults)
	if total == 0 {
		b.WriteString("\n")
		b.WriteString(hintStyle.Render("  Enter a kanji or keyword and press Enter to search"))
		b.WriteString("\n")
		return
	}

	end := total
	used := 0
	for i := m.scrollOffset; i < total; i++ {
		l := resultLines(m.kanjiResults[i])
		if used+l > maxLines {
			end = i
			break
		}
		used += l
	}

	visible := m.kanjiResults[m.scrollOffset:end]

	if m.scrollOffset > 0 || end < total {
		b.WriteString(scrollInfo(total, m.scrollOffset, m.scrollOffset+len(visible)-1))
		b.WriteString("\n")
	}

	for _, r := range visible {
		b.WriteString(fmt.Sprintf(
			"  %s  %s",
			kanjiCharStyle.Render(r.Character),
			readingStyle.Render(r.Meaning),
		))
		b.WriteString("\n")

		if len(r.Kunyomi) > 0 {
			b.WriteString(kanjiInfoStyle.Render(fmt.Sprintf("    Kun: %s", strings.Join(r.Kunyomi, ", "))))
			b.WriteString("\n")
		}
		if len(r.Onyomi) > 0 {
			b.WriteString(kanjiInfoStyle.Render(fmt.Sprintf("    On:  %s", strings.Join(r.Onyomi, ", "))))
			b.WriteString("\n")
		}
		if r.Strokes > 0 || r.Grade > 0 {
			parts := []string{}
			if r.Strokes > 0 {
				parts = append(parts, fmt.Sprintf("%d strokes", r.Strokes))
			}
			if r.Grade > 0 {
				parts = append(parts, fmt.Sprintf("grade %d", r.Grade))
			}
			b.WriteString(kanjiInfoStyle.Render("    " + strings.Join(parts, ", ")))
			b.WriteString("\n")
		}
	}
}

func (m model) renderSentenceResults(b *strings.Builder, maxLines int) {
	total := len(m.sentenceResults)
	if total == 0 {
		b.WriteString("\n")
		b.WriteString(hintStyle.Render("  Enter a word and press Enter to search sentences"))
		b.WriteString("\n")
		return
	}

	end := total
	needed := 0
	for i := m.scrollOffset; i < total; i++ {
		needed += resultLines(m.sentenceResults[i])
		if needed > maxLines {
			end = i
			break
		}
	}

	visible := m.sentenceResults[m.scrollOffset:end]

	if m.scrollOffset > 0 || end < total {
		b.WriteString(scrollInfo(total, m.scrollOffset, m.scrollOffset+len(visible)-1))
		b.WriteString("\n")
	}

	for _, r := range visible {
		b.WriteString("  " + sentenceStyle.Render(r.Content))
		b.WriteString("\n")
		if r.Furigana != "" {
			b.WriteString("  " + kanjiInfoStyle.Render(r.Furigana))
			b.WriteString("\n")
		}
		b.WriteString("  " + translationStyle.Render("→ "+r.Translation))
		b.WriteString("\n")
	}
}

func pulseValue(frame int) float64 {
	return math.Sin(float64(frame)*0.08)*0.5 + 0.5
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
	dot := lipgloss.NewStyle().
		Foreground(lipgloss.Color(pulseGreen)).
		Bold(true)
	return dot.Render("●") + "  "
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
