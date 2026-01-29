# ask

A fast CLI tool for single-shot LLM queries directly from your terminal. Supports multiple providers: Anthropic, OpenAI, Gemini, xAI, and Ollama.

## Install

### Homebrew (macOS / Linux)

```bash
brew install --cask laurensent/tap/ask
```

### Binary download

Download from [GitHub Releases](https://github.com/laurensent/ask/releases) and place in your PATH.

### Go install

```bash
go install github.com/laurensent/ask@latest
```

CLI mode requires [Claude Code CLI](https://docs.anthropic.com/en/docs/claude-code) in PATH. API mode calls providers directly.

## Quick start

```bash
ask "how to rebase"                        # query with quotes
ask how to rebase                          # without quotes (simple text)
ask                                        # interactive mode (no shell escaping needed)
ask -m opus "hard question"                # specify model (provider-specific aliases)
ask --raw "question"                       # skip markdown rendering
git diff | ask "review this code"          # pipe input
cat error.log | ask "analyze this error"   # pipe + prompt
ask "question" | pbcopy                    # auto raw when piped
ask --dry-run "question"                   # preview command
ask --system-prompt "be concise" "question" # passthrough to claude (cli mode)
```

Output streams in real-time and is re-rendered with [glamour](https://github.com/charmbracelet/glamour) markdown styling on completion.

## Providers

| Provider | Models (aliases) | Env var |
|----------|-----------------|---------|
| anthropic | sonnet, opus, haiku | `ANTHROPIC_API_KEY` |
| openai | gpt4o, gpt4o-mini, o3-mini, o4-mini | `OPENAI_API_KEY` |
| gemini | flash, pro, flash-lite | `GEMINI_API_KEY` |
| xai | grok3, grok3-mini | `XAI_API_KEY` |
| ollama | llama3, qwen, deepseek | (none) |

Use `ask config` to select your provider and model, or pass any full model ID directly with `-m`.

## Interactive mode

Run `ask` with no arguments to enter interactive mode. Input is read directly from the terminal, bypassing shell parsing entirely -- no quoting needed for special characters like `'`, `?`, `*`, `&&`, `!`, etc.

```
$ ask
> what's the difference between && and ||?
```

Full emacs keybindings (Ctrl+A/E/B/F/K/U/W), ESC or Ctrl+C to cancel.

## History

```bash
ask history          # open interactive history browser
ask h                # alias
ask history clear    # clear all history
```

Opens a full-screen fuzzy finder to browse and re-run past queries. Type `/` to filter, arrow keys or `j`/`k` to navigate, Enter to re-run, ESC to cancel.

## Config

`ask config` launches an interactive TUI wizard to configure settings.

Config file: `~/.config/ask/config.json`

```json
{
  "mode": "cli",
  "provider": "anthropic",
  "api_key": "",
  "base_url": "",
  "default_model": "sonnet",
  "raw_output": false,
  "theme": "auto"
}
```

| Key | Description |
|-----|-------------|
| `mode` | `cli` (uses Claude Code CLI) or `api` (direct API calls) |
| `provider` | `anthropic`, `openai`, `gemini`, `xai`, or `ollama` |
| `api_key` | API key for the selected provider |
| `base_url` | Custom base URL (for OpenAI-compatible endpoints) |
| `default_model` | Default model alias or full model ID |
| `raw_output` | Skip markdown rendering by default |
| `theme` | Glamour theme: `auto`, `dark`, `light`, `dracula`, `pink`, `ascii`, `notty` |

## License

MIT
