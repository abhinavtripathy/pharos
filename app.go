package main

import (
	"fmt"
	"log"

	"gopkg.in/jdkato/prose.v2"
)

func main() {
	// Create a new document with the default configuration:
	doc, err := prose.NewDocument("Go is an open-source programming language created at Google.")
	if err != nil {
		log.Fatal(err)
	}

	// Iterate over the doc's tokens:
	for _, tok := range doc.Tokens() {
		fmt.Println(tok.Text, tok.Tag, tok.Label)
		// Go NNP B-GPE
		// is VBZ O
		// an DT O
		// ...
	}

	// Iterate over the doc's named-entities:
	for _, ent := range doc.Entities() {
		fmt.Println(ent.Text, ent.Label)
		// Go GPE
		// Google GPE
	}

	// Iterate over the doc's sentences:
	for _, sent := range doc.Sentences() {
		fmt.Println(sent.Text)
		// Go is an open-source programming language created at Google.
	}
}
