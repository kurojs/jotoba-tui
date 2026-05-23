package jotoba

type WordResult struct {
	Word     string
	Reading  string
	Meanings []string
}

type KanjiResult struct {
	Character string
	Meaning   string
	Kunyomi   []string
	Onyomi    []string
	Grade     int
	Strokes   int
}

type SentenceResult struct {
	Content   string
	Furigana  string
	Translation string
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

type jotobaWordsResponse struct {
	Words []jotobaWord `json:"words"`
}

type jotobaKanji struct {
	Literal  string   `json:"literal"`
	Meanings []string `json:"meanings"`
	Kunyomi  []string `json:"kunyomi"`
	Onyomi   []string `json:"onyomi"`
	Grade    *int     `json:"grade"`
	Strokes  *int     `json:"strokes"`
}

type jotobaKanjiResponse struct {
	Kanji []jotobaKanji `json:"kanji"`
}

type jotobaSentence struct {
	Content     string            `json:"content"`
	Furigana    string            `json:"furigana"`
	Translation map[string]string `json:"translation"`
}

type jotobaSentencesResponse struct {
	Sentences []jotobaSentence `json:"sentences"`
}
