package interpreter_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/abhinavtripathy/pharos/interpreter"
)

// Every worked example on the website is marked up as a program followed by
// what it prints:
//
//	<div class="run">
//	<pre><code class="pharos">print "hi"</code></pre>
//	<pre class="out"><code>hi</code></pre>
//	</div>
var runBlock = regexp.MustCompile(
	`(?s)<div class="run">\s*` +
		`<pre><code class="pharos">(.*?)</code></pre>\s*` +
		`<pre class="out"><code>(.*?)</code></pre>\s*` +
		`</div>`)

// TestDocsExamples runs every snippet on the site and checks it really prints
// what the page claims. Without this the guide would quietly become fiction
// the first time an error message is reworded.
func TestDocsExamples(t *testing.T) {
	pages, err := filepath.Glob(filepath.Join("..", "docs", "*.html"))
	if err != nil {
		t.Fatalf("could not look for pages: %v", err)
	}
	if len(pages) == 0 {
		t.Fatal("no pages found; expected .html files in docs/")
	}

	total := 0
	for _, page := range pages {
		body, err := os.ReadFile(page)
		if err != nil {
			t.Fatalf("could not read %s: %v", page, err)
		}

		blocks := runBlock.FindAllStringSubmatch(string(body), -1)
		for i, block := range blocks {
			total++
			source := unescape(block[1])
			want := strings.TrimSpace(unescape(block[2]))

			t.Run(filepath.Base(page), func(t *testing.T) {
				var out strings.Builder
				got := ""
				if err := interpreter.Run(source, &out); err != nil {
					// A snippet may be showing off an error message, in which
					// case the message itself is the expected output.
					got = strings.TrimSpace(err.Error())
				} else {
					got = strings.TrimSpace(out.String())
				}

				if got != want {
					t.Errorf("%s, example %d:\n--- program ---\n%s\n--- page says ---\n%s\n--- actually ---\n%s",
						filepath.Base(page), i+1, source, want, got)
				}
			})
		}
	}

	if total == 0 {
		t.Fatal("found no worked examples in docs/; has the markup changed?")
	}
	t.Logf("checked %d worked examples across %d pages", total, len(pages))
}

func unescape(s string) string {
	r := strings.NewReplacer(
		"&lt;", "<",
		"&gt;", ">",
		"&quot;", `"`,
		"&#39;", "'",
		"&amp;", "&",
	)
	return r.Replace(s)
}
