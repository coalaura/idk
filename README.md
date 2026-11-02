# idk

![Terminal demo: idk finds the process using port 4096](.github/terminal.svg)

Turn plain-English questions into a single, ready-to-run shell command with a short explanation. idk detects your OS/shell and common tools, renders a structured response and (optionally) asks if you want to execute the command.

- Powered by OpenRouter chat completions
- Detects parent shell automatically (Windows/macOS/Linux)
- Tailors suggestions to commands available on your system
- One command per answer, plus a concise explanation

## Install

There are prebuilt, stable releases available [here](https://github.com/coalaura/idk/releases). Alternatively you can also build idk from source.

You can bootstrap **idk** with a single command. This script will detect your OS and CPU (`amd64`/`arm64`), download the correct binary and install it to `/usr/local/bin/idk`.

```bash
curl -sL https://src.ws2.sh/idk/install.sh | sh
```

## Configure

On first run, idk creates a config at `~/.config/idk.yml` with defaults and exits so you can add your API key.

```yaml
openrouter:
  api_key: your-api-key
  model: google/gemini-3.5-flash-lite
  fallbacks:
    - google/gemini-3.8-flash
  providers:
  - google-vertex/global
commands:
  ask_to_run: true
```

- `openrouter.api_key`: your OpenRouter API key
- `openrouter.model`: model slug (defaults to `google/gemini-3.5-flash-lite`). Models and providers, including fallbacks, must support structured outputs (`json_schema`).
- `openrouter.fallbacks`: optional model slugs in fallback order. OpenRouter receives the primary model followed by this list and handles fallback routing natively. Omit it or use `[]` to disable model fallbacks.
- `openrouter.providers`: allowed providers, in preference order, for all models including fallbacks. The generated config defaults to `[google-vertex/global]`; set this to `[]` to allow any provider when using models from other providers.
- `commands.ask_to_run`: if true, idk prompts to execute commands

## Usage

```
# Ask a question
idk list the largest files in this directory
```

idk will display a response like:

```
$ du -sh * | sort -h
  Summarizes sizes of entries and sorts ascending.
```

If `ask_to_run` is enabled, you'll be prompted:

```
> run command? [y/n]:
```

Press `y` to execute in your current shell (Windows uses `/C`, others use `-c`).

## Build from source

- Requirements: Go 1.26.5+
- Build: `go build -o idk` (Windows will produce `idk.exe`)
