# jotoba-tui

A terminal-based Japanese dictionary TUI powered by the [Jotoba API](https://jotoba.de).

Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss).

## Features

- Three search modes: Words, Kanji, and Sentences (switch with Tab)
- Word search: kanji display, kana reading, definitions in English (with 8 additional languages)
- Kanji search: meanings, kunyomi, onyomi, stroke count, grade level
- Sentence search: example sentences with furigana and translations
- Language selector (Ctrl+L) with 9 languages — English, German, Spanish, French, Russian, Swedish, Dutch, Hungarian, Slovenian
- Language preference persisted across sessions
- Scrollable results with Up/Down
- Automatic retry with exponential backoff on network errors

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

![Words](https://i.imgur.com/Dzceyhx.png)

Search by kanji or kana. Shows reading, definitions in the selected language, and scrollable results.

### Kanji search

![Kanji](https://i.imgur.com/O02KNM8.png)

Look up kanji by character or keyword. Displays meaning, kunyomi, onyomi, stroke count, and grade level.

### Sentence search

![Sentences](https://i.imgur.com/zgxjO6c.png)

Find example sentences with furigana and translations in the selected language.

## Language selection

Press Ctrl+L to open the language menu. 9 languages available:

**English** (default), German, Spanish, French, Russian, Swedish, Dutch, Hungarian, Slovenian

The selected language persists across sessions via `~/.config/jotoba-tui/config.json`.

Definitions default to English if no translation is available for the selected language.

## Keys

| Key | Action |
|-----|--------|
| Enter | Search |
| Tab | Switch mode (Words/Kanji/Sentences) |
| Ctrl+L | Language selector |
| Up/Down | Scroll results |
| Esc / Ctrl+C | Quit |

## Project structure

```
cmd/jotoba/main.go       Entry point
internal/config/          Persistent config (language)
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
