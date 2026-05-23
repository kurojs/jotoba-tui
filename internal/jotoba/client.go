package jotoba

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
)

const baseURL = "https://jotoba.de"

var httpClient = &http.Client{Timeout: 15 * time.Second}

func retryDo(req *http.Request, maxRetries int) (*http.Response, error) {
	var lastErr error
	for attempt := range maxRetries {
		if attempt > 0 {
			backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
			time.Sleep(backoff)
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		return resp, nil
	}
	return nil, fmt.Errorf("request failed after %d retries: %w", maxRetries, lastErr)
}

func postJSON(endpoint string, payload any, dest any, language string) error {
	body := map[string]any{
		"query":    payload,
		"language": language,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequest("POST", baseURL+endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := retryDo(req, 3)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	return nil
}

func SearchWords(query, language string) ([]WordResult, error) {
	var resp jotobaWordsResponse
	if err := postJSON("/api/search/words", query, &resp, language); err != nil {
		return nil, err
	}

	if len(resp.Words) == 0 {
		return nil, nil
	}

	results := make([]WordResult, 0, len(resp.Words))
	for _, w := range resp.Words {
		display := w.Reading.Kanji
		if display == "" {
			display = w.Reading.Kana
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

		meanings := otherGlosses
		if len(esGlosses) > 0 {
			meanings = esGlosses
		}

		results = append(results, WordResult{
			Word:     display,
			Reading:  w.Reading.Kana,
			Meanings: meanings,
		})
	}

	return results, nil
}

func SearchKanji(query, language string) ([]KanjiResult, error) {
	var resp jotobaKanjiResponse
	if err := postJSON("/api/search/kanji", query, &resp, language); err != nil {
		return nil, err
	}

	if len(resp.Kanji) == 0 {
		return nil, nil
	}

	results := make([]KanjiResult, 0, len(resp.Kanji))
	for _, k := range resp.Kanji {
		meaning := ""
		if len(k.Meanings) > 0 {
			meaning = k.Meanings[0]
		}

		grade := 0
		if k.Grade != nil {
			grade = *k.Grade
		}

		strokes := 0
		if k.Strokes != nil {
			strokes = *k.Strokes
		}

		results = append(results, KanjiResult{
			Character: k.Literal,
			Meaning:   meaning,
			Kunyomi:   k.Kunyomi,
			Onyomi:    k.Onyomi,
			Grade:     grade,
			Strokes:   strokes,
		})
	}

	return results, nil
}

func langToCode(language string) string {
	switch language {
	case "Spanish":
		return "es"
	case "German":
		return "de"
	case "French":
		return "fr"
	case "Russian":
		return "ru"
	case "Swedish":
		return "sv"
	case "Dutch":
		return "nl"
	case "Hungarian":
		return "hu"
	case "Slovenian":
		return "sl"
	default:
		return "en"
	}
}

func SearchSentences(query, language string) ([]SentenceResult, error) {
	var resp jotobaSentencesResponse
	if err := postJSON("/api/search/sentences", query, &resp, language); err != nil {
		return nil, err
	}

	if len(resp.Sentences) == 0 {
		return nil, nil
	}

	code := langToCode(language)

	results := make([]SentenceResult, 0, len(resp.Sentences))
	for _, s := range resp.Sentences {
		translation := s.Translation["en"]
		if t, ok := s.Translation[code]; ok && t != "" && code != "en" {
			translation = t
		}

		results = append(results, SentenceResult{
			Content:     s.Content,
			Furigana:    s.Furigana,
			Translation: translation,
		})
	}

	return results, nil
}
