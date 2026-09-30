package interpreter_test

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/abhinavtripathy/pharos/interpreter"
)

var update = flag.Bool("update", false, "rewrite the .out files from what the examples currently print")

// TestExamples runs every program in examples/ and compares what it prints
// against the committed .out file beside it. The guide quotes these same
// files, so this is what stops the website from claiming output the language
// no longer produces.
func TestExamples(t *testing.T) {
	programs, err := filepath.Glob(filepath.Join("..", "examples", "*.pharos"))
	if err != nil {
		t.Fatalf("could not look for examples: %v", err)
	}
	if len(programs) == 0 {
		t.Fatal("no examples found; expected .pharos files in examples/")
	}

	for _, path := range programs {
		t.Run(filepath.Base(path), func(t *testing.T) {
			source, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("could not read %s: %v", path, err)
			}

			var out strings.Builder
			if err := interpreter.Run(string(source), &out); err != nil {
				t.Fatalf("%s failed to run:\n%v", path, err)
			}

			golden := strings.TrimSuffix(path, ".pharos") + ".out"
			if *update {
				if err := os.WriteFile(golden, []byte(out.String()), 0o644); err != nil {
					t.Fatalf("could not write %s: %v", golden, err)
				}
				return
			}

			want, err := os.ReadFile(golden)
			if err != nil {
				t.Fatalf("could not read %s (run `go test ./interpreter -update` to create it): %v", golden, err)
			}

			if out.String() != string(want) {
				t.Errorf("%s printed:\n%s\nbut %s says:\n%s",
					filepath.Base(path), out.String(), filepath.Base(golden), want)
			}
		})
	}
}
