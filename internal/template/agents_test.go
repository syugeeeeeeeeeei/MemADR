package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadmeContainsAgentsSnippet(t *testing.T) {
	body, err := os.ReadFile(filepath.Join("..", "..", "README.md"))
	if err != nil {
		t.Fatalf("read README: %v", err)
	}

	want := "```md\n" + RenderAgentsSnippet() + "\n```"
	got := strings.ReplaceAll(string(body), "\r\n", "\n")
	if !strings.Contains(got, want) {
		t.Fatalf("README.md does not contain the AGENTS snippet block")
	}
}
