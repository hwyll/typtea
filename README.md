<h1 align="center">typtea</h1>

<p align="center">Minimal terminal-based typing speed tester with support for dozens of programming languages</p>

<div align="center">
<img src="assets/example.gif" />
</div>

<br>

---

## Features

- **Terminal-based typing** with WPM and accuracy tracking
- **Multi-language support** including English and 30+ programming languages
- **Infinite word generation** with smooth 3-line scrolling display
- **Minimalist TUI** built with Bubble Tea and Lipgloss
- **Embedded language data** for easy distribution
- **Accurate metrics** following standard typing test calculations

### Supported Languages

| | | | |
|-----------|-----------|-----------|-----------|
| Bash      | C         | C++       | C#        |
| Crystal   | CSS       | Emacs     | English 1k|
| Erlang    | Go        | Haskell   | HTML      |
| Java      | JavaScript| JSON      | Julia     |
| Lisp      | Lua       | OCaml     | Perl      |
| PHP       | PowerShell| Python    | R         |
| Ruby      | Rust      | SCSS      | SQL       |
| Swift     | TeX       | TypeScript| Vala      |
| Vimscript | Wolfram   | YAML      | Zig       |
| | | | |

---

## Installation

### Via brew

```bash
brew install hwyll/tap/typtea
```

### Via `go install`

```bash
go install github.com/hwyll/typtea@latest
```

### Build from Source

```bash
git clone --depth=1 https://github.com/hwyll/typtea
cd typtea/
go build
sudo mv typtea /usr/local/bin/
typtea -h
```

---

## Usage

### Basic Commands

```yaml
# Start a 30-second English typing test (default)
typtea start

# Start a 60-second typing test
typtea start --duration 60

# Start a Rust keywords typing test
typtea start --lang rust

# Combine duration and language
typtea start --duration 45 --lang javascript

# List all available languages
typtea start --list-langs

# Get help
typtea --help
typtea start --help
```

### During the Test

- **The test starts** when you begin typing
- **Backspace** to correct mistakes
- **Enter** to restart after completion
- **Esc** to quit the application

---

## Development

### Prerequisites

[Go 1.19+](https://go.dev/doc/install)

### Setup

```bash
git clone https://github.com/hwyll/typtea.git
cd typtea
go mod tidy
go build
./typtea start
```

### Adding New Languages

1. Create a `JSON` file in `internal/game/data/` with the format:

```json
{
  "name": "Language Name",
  "words": ["word1", "word2", "word3", ...]
}
```

2. Rebuild the application to embed the new language data

---

## Dependencies

- [**Bubble Tea**](https://github.com/charmbracelet/bubbletea) - TUI framework
- [**Lipgloss**](https://github.com/charmbracelet/lipgloss) - Styling and layout
- [**Cobra**](https://github.com/spf13/cobra) - CLI framework

## License & Attribution
<p align="left">
        <i><code>Original work &copy; 2025-present <a href="https://github.com/ashish0kumar">Ashish Kumar</a></code></i>
        <br>
        <i><code>Fork modifications &copy; 2025-present <a href="https://github.com/hwyll">hwyll</a></code></i>
</p>

Forked from [ashish0kumar/typtea](https://github.com/ashish0kumar/typtea)

<div align="center">
<a href="https://github.com/hwyll/typtea/blob/main/LICENSE"><img src="https://img.shields.io/github/license/hwyll/typtea?style=for-the-badge&color=CBA6F7&logoColor=cdd6f4&labelColor=302D41" alt="LICENSE"></a>&nbsp;&nbsp;
</div>
