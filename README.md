[![Go](https://github.com/abhinavtripathy/pharos/actions/workflows/go.yml/badge.svg)](https://github.com/abhinavtripathy/pharos/actions/workflows/go.yml)
[![License](http://img.shields.io/badge/License-MIT-brightgreen.svg)](./LICENSE)
![GitHub stars](https://img.shields.io/github/stars/abhinavtripathy/pharos.svg)

# Pharos

> Pharos *noun* — a lighthouse, to steer beginners by

Pharos is a small programming language that reads like plain sentences. It
exists for the awkward moment after Scratch, when coloured blocks stop being
enough but curly braces and semicolons are still a wall.

**[Read the guide →](https://abhinavtripathy.github.io/pharos/)**

```pharos
num score = 82

if score is greater than 89 then
  print "an A"
otherwise if score is greater than 79 then
  print "a B"
otherwise
  print "keep going"
end if
```

```
a B
```

## Why

Most of what trips up a beginner in a first real language is punctuation, not
programming. A missing semicolon, a stray brace, an `elif` that should have
been `else if`. Pharos spells things out instead: `is greater than` rather than
`>`, `otherwise` rather than `else`, `end if` rather than a closing brace.

The ideas underneath are the ordinary ones — typed values, conditions, loops,
named blocks of reusable program — so nothing has to be unlearned later. Only
the spelling is gentler.

Mistakes are explained rather than punished. Every error names its line in
ordinary words:

```
line 2: "score" was declared as num, so it cannot hold a string
```

A loop whose condition can never be met is stopped and reported rather than
left to spin, and so is a codeblock that calls itself without end.

## Install

Pharos needs [Go 1.23 or newer](https://go.dev/dl/).

```sh
go install github.com/abhinavtripathy/pharos/cmd/pharos@latest
```

Or from a clone:

```sh
git clone https://github.com/abhinavtripathy/pharos.git
cd pharos
go build -o pharos ./cmd/pharos
```

## Use

```sh
pharos run examples/hello.pharos   # run a program
pharos examples/hello.pharos       # the same, with less typing
pharos repl                        # try things one line at a time
pharos version
```

The prompt keeps everything you have typed and shows the value of whatever you
write:

```
$ pharos repl
pharos> num x = 20
pharos> x + 1
21
pharos> print "hi " + x
hi 20
```

## The language

| | |
|---|---|
| Values | `num`, `string`, `bool` |
| Declare | `num score = 7` |
| Print | `print x`, or `out` / `output` |
| Compare | `is greater than`, `is less than`, `is equal to`, `is not equal to` |
| Combine | `and`, `or`, `not` |
| Arithmetic | `+`, `-`, `*`, `/`, where `+` also joins text |
| Choose | `if ... then` / `otherwise if ... then` / `otherwise` / `end if` |
| Repeat | `loop until ...` / `end loop` |
| Reuse | `codeblock name with a, b` / `give back x` / `end codeblock` |
| Note | `-- to the end of the line` |

The declared type is enforced, so `num score` is a promise that `score` stays a
number. Keywords ignore case; variable names do not.

The [reference](https://abhinavtripathy.github.io/pharos/reference.html) is the
complete list, and [examples/](./examples) holds runnable programs.

## Layout

```
cmd/pharos/          the command line tool
interpreter/         lexer -> parser -> evaluator, plus the Run entry point
  token/             the lexical vocabulary, including multi-word keywords
  lexer/             source text to tokens
  ast/               the shape of a parsed program
  parser/            recursive descent, with Pratt precedence for expressions
  object/            runtime values and scopes
  evaluator/         walks the tree and runs it
examples/            sample programs, each beside its expected output
docs/                the website, served by GitHub Pages
_archive/            earlier experiments, kept for history and not built
```

Pharos is a tree-walking interpreter with no dependencies outside the Go
standard library. A bytecode compiler would be the natural next step, reusing
the same AST.

## Development

```sh
go test ./...                        # everything
go test ./interpreter -update        # re-record examples/*.out after a change
```

Two tests are worth knowing about. `TestExamples` runs every program in
`examples/` and compares it against the committed `.out` file beside it.
`TestDocsExamples` pulls every worked example out of `docs/*.html`, runs it, and
checks the page really prints what it claims — so the website cannot drift away
from the language.

The site is plain HTML and CSS with no build step. Preview it with:

```sh
python3 -m http.server -d docs
```

## Contributing

Issues and pull requests are welcome; see [CONTRIBUTING.md](./CONTRIBUTING.md).
Things that would help most: a remainder operator, a way to read input, lists,
and better error messages.

## Licence

MIT — see [LICENSE](./LICENSE).
