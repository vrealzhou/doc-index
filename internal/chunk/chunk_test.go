package chunk

import (
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	content := `# Main Title

Some intro text here.

## Section One

This is the first section with some content.

## Section Two

This is the second section with more content.
And a second paragraph.`

	result := Split(content, "test.md", 1000, 100)

	if result.Meta.DocID != "test.md" {
		t.Errorf("expected doc_id test.md, got %s", result.Meta.DocID)
	}

	if result.Meta.Version != 1 {
		t.Errorf("expected version 1, got %d", result.Meta.Version)
	}

	if len(result.Chunks) == 0 {
		t.Error("expected at least one chunk")
	}
}

func TestSplitShortContent(t *testing.T) {
	content := "Short content"
	result := Split(content, "short.md", 1000, 100)

	if len(result.Chunks) != 1 {
		t.Errorf("expected 1 chunk for short content, got %d", len(result.Chunks))
	}

	if result.Chunks[0].Title != "Overview" {
		t.Errorf("expected Overview title, got %s", result.Chunks[0].Title)
	}
}

func TestSplitLongSection(t *testing.T) {
	var paragraphs []string
	for i := 0; i < 100; i++ {
		paragraphs = append(paragraphs, "This is a paragraph that adds to the content.")
	}
	longText := strings.Join(paragraphs, "\n\n")

	content := "## Long Section\n\n" + longText
	result := Split(content, "long.md", 200, 50)

	if len(result.Chunks) < 2 {
		t.Errorf("expected multiple chunks for long section, got %d", len(result.Chunks))
	}
}

func TestComputeHash(t *testing.T) {
	content := "test content"
	hash := ComputeHash(content)

	if len(hash) != 16 {
		t.Errorf("expected hash length 16, got %d", len(hash))
	}

	sameHash := ComputeHash(content)
	if hash != sameHash {
		t.Error("hash should be deterministic")
	}

	differentHash := ComputeHash("different content")
	if hash == differentHash {
		t.Error("different content should produce different hash")
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLen    int
		wantLen   int
		wantTrail bool
	}{
		{"short", "abc", 10, 3, false},
		{"exact", "abcdefghij", 10, 10, false},
		{"truncate", "abcdefghijklmnop", 10, 13, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Truncate(tt.input, tt.maxLen)
			if len(result) != tt.wantLen {
				t.Errorf("expected length %d, got %d", tt.wantLen, len(result))
			}
			hasTrail := len(result) > 0 && result[len(result)-3:] == "..."
			if hasTrail != tt.wantTrail {
				t.Errorf("expected trailing '...' = %v, got %v", tt.wantTrail, hasTrail)
			}
		})
	}
}
