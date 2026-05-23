# jotoba-tui

A terminal-based Japanese dictionary TUI powered by the [Jotoba API](https://jotoba.de).

Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss).

## Features

- Three search modes: Words, Kanji, and Sentences (switch with Tab)
- Word search: kanji display, kana reading, definitions in Spanish (with English fallback)
- Kanji search: meanings, kunyomi, onyomi, stroke count, grade level
- Sentence search: example sentences with furigana and translations
- Automatic retry with exponential backoff on network errors
- Spanish-first: shows Spanish glosses when available, falls back to English

## Installation

### Arch Linux (AUR)

```bash
yay -S jotoba-tui
```

### From source

```bash
git clone https://github.com/kurojs/jotoba-tui.git
cd jotoba-tui
go build -o jotoba-tui ./cmd/jotoba
```

## Usage

```bash
jotoba
```

Type a Japanese word, kanji, or phrase and press Enter. Use Tab to cycle between Words, Kanji, and Sentences modes. Press Ctrl+C or Esc to quit.

### Word search

```
  [Words]   Kanji   Sentences  (Tab to switch)

  > 食べる

  食べる → たべる
    - comer, tomar
    - to eat
```

### Kanji search

```
  Words   [Kanji]   Sentences  (Tab to switch)

  > 食

  食  eat, food
    Kun: た.べる, く.う
    On:  ショク, ジキ
    9 strokes, grade 3
```

### Sentence search

```
  Words   Kanji   [Sentences]  (Tab to switch)

  > 食べる

  私はりんごを食べます。
  私[わたし] は りんご を 食べ[たべ]ます 。
  -> I eat an apple.
```

## Language preference

Results default to Spanish translations. If no Spanish gloss is available for a given word, the API falls back to English. This is controlled at the API level via the `language` parameter.

## Project structure

```
cmd/jotoba/main.go       Entry point
internal/jotoba/          API client and types
internal/jotoba/client.go HTTP client with retry logic
internal/jotoba/types.go  Shared response types
internal/tui/             Bubbletea model, update, and view
internal/tui/model.go     State and styles
internal/tui/update.go    Message handling and search dispatch
internal/tui/view.go      Terminal rendering
```

## License

MIT
