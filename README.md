# jotoba-tui

A terminal-based Japanese dictionary powered by the [Jotoba API](https://jotoba.de).

Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss).

## Features

- Search Japanese words by kanji or kana
- Displays kana reading and Spanish meanings
- Spanish glosses are preferred; falls back to English when unavailable
- Keep searching without restarting the app

## Installation

```bash
git clone https://github.com/anomalyco/jotoba-tui.git
cd jotoba-tui
go build -o jotoba-tui .
```

## Usage

```bash
./jotoba-tui
```

Type a Japanese word and press Enter. Results show the kana reading and definitions in Spanish.

```
Jotoba — Japanese Dictionary

> 濁る

濁る
━━━━━━━━━━━━━━━━━━━━━━━━
濁る → にごる
  - enturbiarse, ensuciarse, contaminarse, emborrascarse

  Ctrl+C / Esc to quit
```

## How it works

1. The TUI sends a POST request to `https://jotoba.de/api/search/words` with `language: "Spanish"`
2. The response is parsed and displayed inline
3. Glosses tagged as `"Spanish"` are shown first; if none exist, English glosses are used as fallback

## License

MIT
