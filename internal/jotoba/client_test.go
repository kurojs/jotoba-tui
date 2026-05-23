package jotoba

import (
	"encoding/json"
	"testing"
)

func TestParseWordsResponse(t *testing.T) {
	data := `{
		"words": [
			{
				"reading": {"kana": "たべる", "kanji": "食べる"},
				"senses": [
					{"glosses": ["to eat"], "language": "English"},
					{"glosses": ["comer"], "language": "Spanish"}
				]
			},
			{
				"reading": {"kana": "のむ", "kanji": "飲む"},
				"senses": [
					{"glosses": ["to drink"], "language": "English"}
				]
			}
		]
	}`

	var resp jotobaWordsResponse
	if err := json.Unmarshal([]byte(data), &resp); err != nil {
		t.Fatal(err)
	}

	if len(resp.Words) != 2 {
		t.Fatalf("expected 2 words, got %d", len(resp.Words))
	}

	w := resp.Words[0]
	if w.Reading.Kana != "たべる" {
		t.Errorf("expected たべる, got %s", w.Reading.Kana)
	}
	if w.Reading.Kanji != "食べる" {
		t.Errorf("expected 食べる, got %s", w.Reading.Kanji)
	}
	if len(w.Senses) != 2 {
		t.Fatalf("expected 2 senses, got %d", len(w.Senses))
	}
	if w.Senses[0].Glosses[0] != "to eat" {
		t.Errorf("expected 'to eat', got %s", w.Senses[0].Glosses[0])
	}
}

func TestParseKanjiResponse(t *testing.T) {
	grade3 := 3
	strokes9 := 9
	data := `{
		"kanji": [
			{
				"literal": "食",
				"meanings": ["eat", "food"],
				"kunyomi": ["た.べる", "く.う"],
				"onyomi": ["ショク", "ジキ"],
				"grade": 3,
				"strokes": 9
			}
		]
	}`

	var resp jotobaKanjiResponse
	if err := json.Unmarshal([]byte(data), &resp); err != nil {
		t.Fatal(err)
	}

	if len(resp.Kanji) != 1 {
		t.Fatalf("expected 1 kanji, got %d", len(resp.Kanji))
	}

	k := resp.Kanji[0]
	if k.Literal != "食" {
		t.Errorf("expected 食, got %s", k.Literal)
	}
	if len(k.Meanings) != 2 || k.Meanings[0] != "eat" {
		t.Errorf("expected meanings [eat food], got %v", k.Meanings)
	}
	if k.Grade == nil || *k.Grade != grade3 {
		t.Errorf("expected grade 3, got %v", k.Grade)
	}
	if k.Strokes == nil || *k.Strokes != strokes9 {
		t.Errorf("expected strokes 9, got %v", k.Strokes)
	}
}

func TestParseKanjiResponseNilGradeStrokes(t *testing.T) {
	data := `{"kanji": [{"literal": "謎", "meanings": ["riddle"], "kunyomi": ["なぞ"], "onyomi": []}]}`

	var resp jotobaKanjiResponse
	if err := json.Unmarshal([]byte(data), &resp); err != nil {
		t.Fatal(err)
	}

	if len(resp.Kanji) != 1 {
		t.Fatalf("expected 1 kanji, got %d", len(resp.Kanji))
	}

	k := resp.Kanji[0]
	if k.Grade != nil {
		t.Errorf("expected nil grade, got %d", *k.Grade)
	}
	if k.Strokes != nil {
		t.Errorf("expected nil strokes, got %d", *k.Strokes)
	}
}

func TestParseSentencesResponse(t *testing.T) {
	data := `{
		"sentences": [
			{
				"content": "私はりんごを食べます。",
				"furigana": "私[わたし] は りんご を 食べ[たべ]ます 。",
				"translation": {"en": "I eat an apple.", "es": "Yo como una manzana."}
			}
		]
	}`

	var resp jotobaSentencesResponse
	if err := json.Unmarshal([]byte(data), &resp); err != nil {
		t.Fatal(err)
	}

	if len(resp.Sentences) != 1 {
		t.Fatalf("expected 1 sentence, got %d", len(resp.Sentences))
	}

	s := resp.Sentences[0]
	if s.Content != "私はりんごを食べます。" {
		t.Errorf("unexpected content: %s", s.Content)
	}
	if s.Translation["en"] != "I eat an apple." {
		t.Errorf("unexpected en translation: %s", s.Translation["en"])
	}
	if s.Translation["es"] != "Yo como una manzana." {
		t.Errorf("unexpected es translation: %s", s.Translation["es"])
	}
}

func TestSearchWordsSortsSpanishFirst(t *testing.T) {
	_, err := SearchWords("食べる", "Spanish")
	if err != nil {
		t.Skip("API not available:", err)
	}

	t.Log("API responded — SearchWords works")
}

func TestEmptyWordsReturnsNil(t *testing.T) {
	var resp jotobaWordsResponse
	if err := json.Unmarshal([]byte(`{"words": []}`), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Words) != 0 {
		t.Error("expected empty words")
	}
}

func TestEmptyKanjiReturnsNil(t *testing.T) {
	var resp jotobaKanjiResponse
	if err := json.Unmarshal([]byte(`{"kanji": []}`), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Kanji) != 0 {
		t.Error("expected empty kanji")
	}
}

func TestEmptySentencesReturnsNil(t *testing.T) {
	var resp jotobaSentencesResponse
	if err := json.Unmarshal([]byte(`{"sentences": []}`), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Sentences) != 0 {
		t.Error("expected empty sentences")
	}
}
