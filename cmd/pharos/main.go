// Command pharos runs Pharos programs, either from a file or from an
// interactive prompt.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/abhinavtripathy/pharos/interpreter"
	"github.com/abhinavtripathy/pharos/interpreter/lexer"
	"github.com/abhinavtripathy/pharos/interpreter/token"
)

const version = "0.1.0"

const usage = `Pharos is a small language for people writing their first programs.

usage:
  pharos run <file.pharos>   run a program
  pharos <file.pharos>       the same thing, with less typing
  pharos repl                try things out one line at a time
  pharos version             print the version

The guide lives at https://abhinavtripathy.github.io/pharos/
`

func main() {
	if code := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); code != 0 {
		os.Exit(code)
	}
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return 2
	}

	switch args[0] {
	case "version", "--version", "-v":
		fmt.Fprintf(stdout, "pharos %s\n", version)
		return 0

	case "help", "--help", "-h":
		fmt.Fprint(stdout, usage)
		return 0

	case "repl":
		repl(stdin, stdout)
		return 0

	case "run":
		if len(args) < 2 {
			fmt.Fprintf(stderr, "pharos run needs a file to run, for example: pharos run hello.pharos\n")
			return 2
		}
		return runFile(args[1], stdout, stderr)
	}

	if strings.HasSuffix(args[0], ".pharos") {
		return runFile(args[0], stdout, stderr)
	}

	fmt.Fprintf(stderr, "pharos does not have a command called %q\n\n%s", args[0], usage)
	return 2
}

func runFile(path string, stdout, stderr io.Writer) int {
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(stderr, "pharos could not read %s: %v\n", path, err)
		return 1
	}

	if err := interpreter.Run(string(source), stdout); err != nil {
		fmt.Fprintf(stderr, "%s\n%s\n", path, err)
		return 1
	}
	return 0
}

func repl(stdin io.Reader, stdout io.Writer) {
	fmt.Fprintf(stdout, "pharos %s -- type a line and press enter, or press ctrl-d to leave\n", version)

	session := interpreter.NewSession(stdout)
	scanner := bufio.NewScanner(stdin)

	var pending []string
	for {
		if len(pending) == 0 {
			fmt.Fprint(stdout, "pharos> ")
		} else {
			fmt.Fprint(stdout, "   ...  ")
		}

		if !scanner.Scan() {
			fmt.Fprintln(stdout)
			return
		}

		line := scanner.Text()
		if len(pending) == 0 {
			switch strings.TrimSpace(line) {
			case "":
				continue
			case "exit", "quit":
				return
			}
		}

		pending = append(pending, line)
		source := strings.Join(pending, "\n")

		// Blocks span several lines, so hold off until every "if", "loop" and
		// "codeblock" has been closed.
		if unclosedBlocks(source) > 0 {
			continue
		}
		pending = nil

		value, err := session.Eval(source)
		if err != nil {
			fmt.Fprintln(stdout, err)
			continue
		}
		if value != nil && value.Type() != "nothing" {
			fmt.Fprintln(stdout, value.Inspect())
		}
	}
}

// unclosedBlocks counts block openers that have not been closed yet. It can go
// negative on a stray "end if", which the parser will report properly once the
// entry is submitted.
func unclosedBlocks(source string) int {
	depth := 0
	for _, tok := range lexer.New(source).Tokens() {
		switch tok.Type {
		case token.IF, token.LOOP, token.CODEBLOCK:
			depth++
		case token.ENDIF, token.ENDLOOP, token.ENDCODEBLOCK:
			depth--
		}
	}
	return depth
}
