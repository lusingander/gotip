# gotip

Go Test Interactive Picker

<img src="./img/demo.gif" width=600>

## About

gotip is a TUI application for interactively selecting and running Go tests.

Key features:

- Fuzzy filtering of test cases
- Detection of subtest names defined via table-driven tests (partial support)
- List discovered tests in text or JSON format
- Run individual subtests or grouped subtests
- View and re-run tests from execution history

## Installation

```
go install github.com/lusingander/gotip/cmd/gotip@latest
```

## Usage

### Basic usage

In a directory containing a `go.mod` file, run:

```
gotip
```

While a test is selected, press <kbd>Enter</kbd> to run it using `go test`.

To show only tests and histories from a package and its subpackages, use `--package`:

```
gotip --package ./internal
```

You can also use Go-style package patterns:

```
gotip --package ./internal/...
```

### Passing additional arguments

You can pass extra flags directly to `go test` by appending them after `--`:

```
gotip -- -v -count=1
```

### Running a parent test group

While a test is selected, press <kbd>Backspace</kbd> to move up to its parent test group.

This allows you to execute all subtests under that group.  
For example, if you have `TestFoo/Bar/Baz` selected, pressing <kbd>Backspace</kbd> will select `TestFoo/Bar`, and running it will execute all tests under that prefix.

If subtest names could not be automatically discovered, gotip defaults to selecting the nearest available parent test.

<img src="./img/group.gif" width=600>

### Using test history

Press <kbd>Tab</kbd> to switch to History view, or launch directly with the `--view=history` option.

In this view, you can select and run tests from your previous execution history, just like in the regular view.

The history data is stored under `~/.local/state/gotip/history/`.

<img src="./img/history.gif" width=600>

### Rerunning the last test

You can rerun the last executed test without showing the UI by using the `--rerun` option:

```
gotip --rerun
```

This will immediately execute the most recent test from your history.

When used with `--package`, `--rerun` executes the most recent history entry that matches the package filter:

```
gotip --rerun --package ./internal/...
```

### Listing discovered tests

You can inspect the statically discovered test tree without opening the UI:

```
gotip list
```

Example text output:

```text
# ./foo/foo_test.go
- TestFoo
  - case1
  - case2
- TestBar
```

For machine-readable output, use JSON:

```
gotip list --format=json
```

To list tests from a package and its subpackages, use `--package`. `internal` is treated as `./internal`, and `./internal` matches packages such as `./internal/parse` and `./internal/tip`.

```
gotip list --package ./internal
```

Go-style package patterns are also supported:

```
gotip list --package ./internal/...
gotip list --package ./...
```

Example JSON output:

```json
{
  "files": [
    {
      "path": "./foo/foo_test.go",
      "tests": [
        {
          "name": "TestFoo",
          "subtests": [
            { "name": "case1", "resolved": true, "subtests": [] }
          ]
        }
      ]
    }
  ]
}
```

The JSON output follows [`schema/list.schema.json`](./schema/list.schema.json).

The `list` command uses the same discovery rules as the TUI, including subtest inference and `--skip-subtests`.

### Options

```
Usage:
  gotip [OPTIONS] [list]

Application Options:
  -v, --view=[all|history]      Default view (default: all)
  -f, --filter=[fuzzy|exact]    Default filter type (default: fuzzy)
  -p, --package=PACKAGE         Filter by package path or pattern
  -s, --skip-subtests           Skip subtest detection
  -r, --rerun                   Rerun the last test without showing the UI
  -V, --version                 Print version

Help Options:
  -h, --help                    Show this help message

Available commands:
  list  List discovered tests
```

`gotip list --help` shows options specific to the non-interactive listing command:

```
Usage:
  gotip [OPTIONS] list [list-OPTIONS]

[list command options]
      -p, --package=PACKAGE       Filter by package path or pattern
      -s, --skip-subtests         Skip subtest detection
          --format=[text|json]    Output format (default: text)
```

### Config

gotip supports both global and project-specific configuration.

- Global config
  - Place your global config at `~/.config/gotip/gotip.toml`. This applies to all projects.
- Project config
  - Place a `gotip.toml` file in your current working directory. This applies only to the current project.

If both global and project configs exist, they are merged.  
For overlapping keys, the project config takes precedence.

The format is as follows:

```toml
# Specifies the command used to run tests.
# If omitted, the default command is used.
# type: list of strings
command = []
# Specify file path patterns to exclude from processing using the .gitignore format.
# https://git-scm.com/docs/gitignore/en#_pattern_format
# type: list of strings
ignore = []

[filter]
# Selects the fuzzy matching implementation.
# Supported values are "gotip" and "legacy".
# type: string
fuzzy_matcher = "gotip"

[history]
# Limits the number of test executions to keep in history.
# type: integer
limit = 100
# Format used to display timestamps in the history view.
# Uses Go's time format syntax.
# type: string
date_format = "2006-01-02 15:04:05"

[keybindings]
# Keys use Bubble Tea's canonical key names.
# Each field replaces the corresponding global or default binding.
select_previous = ["up", "k"]
select_next = ["down", "j"]
previous_page = ["left", "h", "pgup", "b", "u"]
next_page = ["right", "l", "pgdown", "f", "d"]
go_to_start = ["home", "g"]
go_to_end = ["end", "G"]
run = ["enter"]
parent = ["backspace", "ctrl+h"]
start_filter = ["/"]
clear_filter = ["esc"]
cancel_filter = ["esc"]
confirm_filter = ["enter"]
toggle_filter_type = ["ctrl+x"]
switch_view = ["tab", "shift+tab"]
show_help = ["?"]
close_help = ["?", "backspace", "ctrl+h"]
quit = ["q", "esc"]
force_quit = ["ctrl+c"]

[theme]
# Colors accept #RGB, #RRGGBB, or ANSI color numbers from 0 to 255.
# Color values must be strings.
# Primary UI text color.
# type: string
text = "#DDDDDD"
# Selected target and filter cursor color.
# type: string
accent = "#00ADD8"
# Selected list item and help key color.
# type: string
highlight = "#5DC9E2"
# Secondary and dimmed text color.
# type: string
muted = "#777777"
# Description color for dimmed items.
# type: string
dimmed = "#4D4D4D"
# Header and footer border color.
# type: string
border = "240"
# Filter query match color.
# type: string
match = "#CE3262"
# Test command color.
# type: string
command = "#00A29C"
```

#### `filter`

The `filter.fuzzy_matcher` field selects the matcher used when the UI is in fuzzy filtering mode.

| Value | Behavior |
| --- | --- |
| `gotip` | Uses gotip's matcher, which prefers compact, contiguous matches. This is the default. |
| `legacy` | Uses the original `sahilm/fuzzy`-based matcher. |

This setting does not change the initial `fuzzy` or `exact` filter mode selected by the `--filter` option.

#### `keybindings`

The `keybindings` table customizes UI actions. Each field is a list because an action can have multiple keys.
When global and project configs both define an action, the project list replaces the global list.
Set an optional action to an empty list to disable it:

```toml
[keybindings]
run = []
```

`force_quit` is the only required action and must contain at least one key.

Key names are case-sensitive and use the strings returned by Bubble Tea's
[`KeyMsg.String`](https://pkg.go.dev/github.com/charmbracelet/bubbletea@v1.3.10#KeyMsg.String):

- Printable characters are written directly, such as `"j"`, `"G"`, or `"?"`. Use `" "` for Space.
- Special keys include `"enter"`, `"backspace"`, `"tab"`, `"esc"`, arrow keys such as `"up"`, paging keys such as `"pgdown"`, and `"f1"` through `"f20"`.
- Modifiers use `+`, as in `"ctrl+n"`, `"alt+x"`, `"shift+tab"`, or `"ctrl+shift+up"`.
- Use canonical names such as `"tab"`, `"enter"`, and `"esc"` rather than their control-key aliases.

For canonical special-key names, see Bubble Tea v1.3.10's
[key name mapping](https://github.com/charmbracelet/bubbletea/blob/v1.3.10/key.go#L260-L350).

The available actions are:

| Field | Action |
| --- | --- |
| `select_previous` | Select the previous item, or scroll help up |
| `select_next` | Select the next item, or scroll help down |
| `previous_page` | Select the previous page |
| `next_page` | Select the next page |
| `go_to_start` | Select the first item |
| `go_to_end` | Select the last item |
| `run` | Run the selected test |
| `parent` | Select the parent test group |
| `start_filter` | Enter filtering mode |
| `clear_filter` | Clear an applied filter |
| `cancel_filter` | Cancel filter editing |
| `confirm_filter` | Confirm filter editing |
| `toggle_filter_type` | Toggle between fuzzy and exact filtering |
| `switch_view` | Switch between All Tests and History |
| `show_help` | Show help |
| `close_help` | Close help |
| `quit` | Quit while browsing |
| `force_quit` | Quit from any screen |

Keys may be reused for actions that are active in different contexts, such as `run` and `confirm_filter`.
Conflicting keys within the same context are rejected when the config is loaded.
`clear_filter` takes precedence over `quit` when both contain the same key and a filter is applied.
The in-app help is generated from the active keybindings and omits disabled actions.

#### `command`

The `command` field allows you to customize how tests are executed.

You can use this to always pass specific flags or use an external test runner instead of the default.

For example, to use [gotestsum](https://github.com/gotestyourself/gotestsum), you can configure it like this:

```toml
command = ["gotestsum", "--format", "testname", "--", "-run", "${name}", "${package}"]
```

`${name}` and `${package}` are placeholders that will be replaced at runtime with the selected test name pattern and package name, respectively.

If not specified, the following default command is used:

```toml
command = ["go", "test", "-run", "${name}", "${package}"]
```

#### `theme`

The `theme` table controls the colors used by the UI and the displayed test command.
You can specify only the colors you want to change; omitted colors keep their global or default values.

| Field | Usage |
| --- | --- |
| `text` | Primary UI text |
| `accent` | Selected target and filter cursor |
| `highlight` | Selected list items and help keys |
| `muted` | Secondary and dimmed text |
| `dimmed` | Descriptions of dimmed items |
| `border` | Header and footer borders |
| `match` | Characters matching the filter query |
| `command` | Test command printed before execution |

### Keybindings

| Key                                                        | Description                                      |
| ---------------------------------------------------------- | ------------------------------------------------ |
| <kbd>Ctrl-c</kbd>                                          | Quit from any screen                             |
| <kbd>q</kbd>                                                | Quit while browsing                              |
| <kbd>Esc</kbd>                                              | Quit while browsing without an applied filter    |
| <kbd>j</kbd> <kbd>↓</kbd>                                 | Select next item / Scroll help down               |
| <kbd>k</kbd> <kbd>↑</kbd>                                 | Select previous item / Scroll help up             |
| <kbd>l</kbd> <kbd>→</kbd> <kbd>PgDown</kbd> <kbd>f</kbd> <kbd>d</kbd> | Select next page                    |
| <kbd>h</kbd> <kbd>←</kbd> <kbd>PgUp</kbd> <kbd>b</kbd> <kbd>u</kbd>  | Select previous page                |
| <kbd>g</kbd> <kbd>Home</kbd>                               | Select first item                                |
| <kbd>G</kbd> <kbd>End</kbd>                                | Select last item                                 |
| <kbd>Enter</kbd>                                           | Run the selected test                            |
| <kbd>Backspace</kbd> <kbd>Ctrl-h</kbd>                     | Select parent test group                         |
| <kbd>/</kbd>                                               | Enter filtering mode                             |
| <kbd>Enter</kbd>                                           | Confirm filter while editing                     |
| <kbd>Esc</kbd>                                             | Cancel filter editing or clear an applied filter |
| <kbd>Ctrl-x</kbd>                                          | Toggle filtering type                            |
| <kbd>Tab</kbd> <kbd>Shift-Tab</kbd>                        | Switch view                                      |
| <kbd>?</kbd>                                               | Show help                                        |
| <kbd>?</kbd> <kbd>Backspace</kbd> <kbd>Ctrl-h</kbd>        | Close help                                       |

## Planned features

- Launch with initial filter based on package or test name

## License

MIT
