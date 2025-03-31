package tree_sitter_go_test

import (
	"testing"

	tree_sitter "github.com/rlch/implgen/go-tree-sitter"
	tree_sitter_go "github.com/rlch/implgen/parser/bindings/go"
)

func TestCanLoadGrammar(t *testing.T) {
	language := tree_sitter.NewLanguage(tree_sitter_go.Language())
	if language == nil {
		t.Errorf("Error loading Go grammar")
	}
}
