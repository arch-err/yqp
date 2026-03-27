# yqp

A TUI playground for exploring [yq](https://github.com/mikefarah/yq) — the YAML equivalent of [jqp](https://github.com/noahgorstein/jqp).

Write yq expressions interactively against YAML data with live results. Queries are copy-pasteable straight into `yq` on the command line.

Built with [mikefarah/yq](https://github.com/mikefarah/yq) under the hood. Forked from [noahgorstein/jqp](https://github.com/noahgorstein/jqp).

<img alt="yqp demo" src="./demo.gif" width="1200" />

## Installation

### go install

```bash
go install github.com/arch-err/yqp@latest
```

### Build from source

```bash
git clone https://github.com/arch-err/yqp.git
cd yqp && go build
mv ./yqp /usr/local/bin
```

## Usage

```
$ yqp --help
yqp is a terminal user interface (TUI) for exploring the yq command line utility.

You can use it to run yq queries interactively. If no query is provided, the interface will prompt you for one.

The command accepts an optional query argument which will be executed against the input YAML.
You can provide the input YAML either through a file or via standard input (stdin).

Usage:
  yqp [query] [flags]

Flags:
      --config string   path to config file (default is $HOME/.yqp.yaml)
  -f, --file string     path to the input YAML file
  -h, --help            help for yqp
  -t, --theme string    yqp theme
  -v, --version         version for yqp
```

`yqp` supports input from STDIN or a file. STDIN takes precedence. You can also pass an initial query argument:

```bash
# from file
yqp -f deployment.yaml

# from stdin
cat values.yaml | yqp

# with an initial query
yqp '.spec.containers[].image' -f pod.yaml

# pipe from kubectl
kubectl get deployment my-app -o yaml | yqp
```

> [!NOTE]
> Queries use [yq expression syntax](https://mikefarah.gitbook.io/yq/). Anything you type in `yqp` can be copied and used directly with the `yq` CLI. Multi-document YAML (with `---` separators) is supported.

## Keybindings

| **Keybinding** | **Action** |
|:---------------|:-----------|
| `tab` | cycle through sections |
| `shift-tab` | cycle through sections in reverse |
| `ctrl-y` | copy query to system clipboard[^1] |
| `ctrl-s` | save output to file (copy to clipboard if file not specified) |
| `ctrl-t` | toggle showing/hiding input panel |
| `ctrl-c` | quit program / kill long-running query |

### Query Mode

| **Keybinding** | **Action** |
|:---------------|:-----------|
| `enter` | execute query |
| `↑`/`↓` | cycle through query history |
| `ctrl-a` | go to beginning of line |
| `ctrl-e` | go to end of line |
| `←`/`ctrl-b` | move cursor one character to left |
| `→`/`ctrl-f`| move cursor one character to right |
| `ctrl-k` | delete text after cursor line |
| `ctrl-u` | delete text before cursor |
| `ctrl-w` | delete word to left |
| `ctrl-d` | delete character under cursor |

### Input Preview and Output Mode

| **Keybinding** | **Action** |
|:---------------|:-----------|
| `↑/k` | up |
| `↓/j` | down |
| `ctrl-u` | page up |
| `ctrl-d` | page down |

## Configuration

`yqp` can be configured with a configuration file. By default, `yqp` will search your home directory for a YAML file named `.yqp.yaml`. A path to a configuration file can also be provided with the `--config` flag.

```bash
yqp --config ~/my_yqp_config.yaml -f data.yaml
```

If a configuration option is present in both the configuration file and the command-line, the command-line option takes precedence.

### Available Configuration Options

```yaml
theme:
  name: "nord" # controls the color scheme
  chromaStyleOverrides: # override parts of the chroma style
    kc: "#009900 underline" # keys use the chroma short names
```

## Themes

Themes can be specified on the command-line via the `-t/--theme <themeName>` flag or in your [configuration file](#configuration).

```yaml
theme:
  name: "monokai"
```

### Chroma Style Overrides

Overrides to the chroma styles used for a theme can be configured in your [configuration file](#configuration).

For the list of short keys, see [`chroma.StandardTypes`](https://github.com/alecthomas/chroma/blob/d38b87110b078027006bc34aa27a065fa22295a1/types.go#L210-L308). To see which token to use for a value, see the [YAML lexer](https://github.com/alecthomas/chroma/blob/master/lexers/embedded/yaml.xml) (look for `<token>` tags). To see the color and what's used in the style you're using, look for your style in the chroma [styles directory](https://github.com/alecthomas/chroma/tree/master/styles).

```yaml
theme:
  name: "monokai"
  chromaStyleOverrides:
    kc: "#009900 underline"
```

You can change non-syntax colors using the `styleOverrides` key:
```yaml
theme:
  styleOverrides:
    primary: "#c4b28a"
    secondary: "#8992a7"
    error: "#c4746e"
    inactive: "#a6a69c"
    success: "#87a987"
```

### Light Themes

`abap`, `algol`, `arduino`, `autumn`, `borland`, `catppuccin-latte`, `colorful`, `emacs`, `friendly`, `github`, `gruvbox-light`, `hrdark`, `igor`, `lovelace`, `manni`, `monokai-light`, `murphy`, `onesenterprise`, `paraiso-light`, `pastie`, `perldoc`, `pygments`, `solarized-light`, `tango`, `trac`, `visual_studio`, `vulcan`, `xcode`

### Dark Themes

`average`, `base16snazzy`, `catppuccin-frappe`, `catppuccin-macchiato`, `catppuccin-mocha`, `doom-one`, `doom-one2`, `dracula`, `fruity`, `github-dark`, `gruvbox`, `monokai`, `native`, `paraiso-dark`, `rrt`, `solarized-dark`, `solarized-dark256`, `swapoff`, `vim`, `witchhazel`, `xcode-dark`

## Built with

- [Bubbletea](https://github.com/charmbracelet/bubbletea)
- [Bubbles](https://github.com/charmbracelet/bubbles)
- [Lipgloss](https://github.com/charmbracelet/lipgloss)
- [yq](https://github.com/mikefarah/yq) (Go library)
- [chroma](https://github.com/alecthomas/chroma)

## Credits

- [jqp](https://github.com/noahgorstein/jqp) by Noah Gorstein — the original TUI playground for jq, which this project is forked from
- [jqq](https://github.com/jcsalterego/jqq) for inspiration

--------

[^1]: `yqp` uses [https://github.com/atotto/clipboard](https://github.com/atotto/clipboard) for clipboard functionality. Things should work as expected with OSX and Windows. Linux, Unix require `xclip` or `xsel` to be installed.
