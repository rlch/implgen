package tree_sitter_go_test

import (
	"testing"

	tree_sitter_go "github.com/rlch/implgen/parser/bindings/go"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func TestCanLoadGrammar(t *testing.T) {
	language := tree_sitter.NewLanguage(tree_sitter_go.Language())
	if language == nil {
		t.Errorf("Error loading Go grammar")
	}
}
