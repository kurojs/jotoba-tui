package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"kurojs.com/jotoba-tui/internal/jotoba"
)

func TestInitialModel(t *testing.T) {
	m := New().(model)

	if m.mode != modeWord {
		t.Errorf("expected modeWord, got %v", m.mode)
	}
	if !m.textInput.Focused() {
		t.Error("expected text input to be focused")
	}
	if m.loading {
		t.Error("expected loading to be false")
	}
}

func TestTabSwitchesMode(t *testing.T) {
	m := New().(model)

	keyMsg := tea.KeyMsg{Type: tea.KeyTab}
	for range 3 {
		result, _ := m.Update(keyMsg)
		m = result.(model)
	}

	if m.mode != modeWord {
		t.Errorf("expected back to modeWord after 3 tabs, got %v", m.mode)
	}
}

func TestTabClearsResults(t *testing.T) {
	m := New().(model)
	m.wordResults = []jotoba.WordResult{{Word: "test"}}

	keyMsg := tea.KeyMsg{Type: tea.KeyTab}
	result, _ := m.Update(keyMsg)
	m = result.(model)

	if len(m.wordResults) != 0 {
		t.Error("expected word results to be cleared on tab switch")
	}
}

func TestTabClearsError(t *testing.T) {
	m := New().(model)
	m.err = &modeError{"some error"}

	keyMsg := tea.KeyMsg{Type: tea.KeyTab}
	result, _ := m.Update(keyMsg)
	m = result.(model)

	if m.err != nil {
		t.Error("expected error to be cleared on tab switch")
	}
}

func TestEnterWithEmptyQueryDoesNothing(t *testing.T) {
	m := New().(model)

	result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if result.(model).loading {
		t.Error("expected no loading on empty query")
	}
}

func TestEnterStartsSearch(t *testing.T) {
	m := New().(model)
	m.textInput.SetValue("test")

	result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m2 := result.(model)

	if !m2.loading {
		t.Error("expected loading to be true after enter")
	}
}

func TestCtrlCQuits(t *testing.T) {
	m := New().(model)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	if cmd == nil {
		t.Fatal("expected quit command")
	}
}

func TestEscQuits(t *testing.T) {
	m := New().(model)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if cmd == nil {
		t.Fatal("expected quit command")
	}
}

func TestSearchResultClearsLoading(t *testing.T) {
	m := New().(model)
	m.loading = true

	result, _ := m.Update(searchResultMsg{
		mode:    modeKanji,
		results: []jotoba.KanjiResult{},
	})
	m2 := result.(model)

	if m2.loading {
		t.Error("expected loading to be false after search result")
	}
}

func TestSearchResultStoresWordResults(t *testing.T) {
	m := New().(model)
	results := []jotoba.WordResult{{Word: "食べる", Reading: "たべる"}}

	m2, _ := m.Update(searchResultMsg{mode: modeWord, results: results})
	m3 := m2.(model)

	if len(m3.wordResults) != 1 {
		t.Fatalf("expected 1 word result, got %d", len(m3.wordResults))
	}
	if m3.wordResults[0].Word != "食べる" {
		t.Errorf("expected 食べる, got %s", m3.wordResults[0].Word)
	}
}

func TestSearchResultStoresKanjiResults(t *testing.T) {
	m := New().(model)
	results := []jotoba.KanjiResult{{Character: "食", Meaning: "eat"}}

	m2, _ := m.Update(searchResultMsg{mode: modeKanji, results: results})
	m3 := m2.(model)

	if len(m3.kanjiResults) != 1 {
		t.Fatalf("expected 1 kanji result, got %d", len(m3.kanjiResults))
	}
	if m3.kanjiResults[0].Character != "食" {
		t.Errorf("expected 食, got %s", m3.kanjiResults[0].Character)
	}
}

func TestSearchResultStoresSentenceResults(t *testing.T) {
	m := New().(model)
	results := []jotoba.SentenceResult{
		{Content: "私はりんごを食べます。", Translation: "I eat an apple."},
	}

	m2, _ := m.Update(searchResultMsg{mode: modeSentence, results: results})
	m3 := m2.(model)

	if len(m3.sentenceResults) != 1 {
		t.Fatalf("expected 1 sentence result, got %d", len(m3.sentenceResults))
	}
	if m3.sentenceResults[0].Content != "私はりんごを食べます。" {
		t.Errorf("unexpected content: %s", m3.sentenceResults[0].Content)
	}
}

func TestErrorMsgClearsLoading(t *testing.T) {
	m := New().(model)
	m.loading = true

	result, _ := m.Update(errorMsg{err: &modeError{"network error"}})
	m2 := result.(model)

	if m2.loading {
		t.Error("expected loading to be false after error")
	}
	if m2.err == nil || m2.err.Error() != "network error" {
		t.Errorf("expected 'network error', got %v", m2.err)
	}
}

func TestModeLabel(t *testing.T) {
	cases := []struct {
		mode searchMode
		want string
	}{
		{modeWord, "Words"},
		{modeKanji, "Kanji"},
		{modeSentence, "Sentences"},
	}

	for _, tc := range cases {
		t.Run(tc.want, func(t *testing.T) {
			got := modeLabel(tc.mode)
			if got != tc.want {
				t.Errorf("modeLabel(%d) = %s, want %s", tc.mode, got, tc.want)
			}
		})
	}
}
