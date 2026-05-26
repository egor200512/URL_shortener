package links

import (
	"strings"
	"testing"
)

func TestGenerateShortLink(t *testing.T) {
	t.Parallel()

	for range 20 {
		shortLink, err := GenerateShortLink()
		if err != nil {
			t.Fatalf("GenerateShortLink() error = %v", err)
		}

		if shortLink == "" {
			t.Fatal("GenerateShortLink() returned empty string")
		}

		if len(shortLink) != LENGTH {
			t.Fatalf("len(shortLink) = %d, want %d", len(shortLink), LENGTH)
		}

		for _, char := range shortLink {
			if !strings.ContainsRune(CHARSET, char) {
				t.Fatalf("shortLink contains invalid char %q; shortLink=%q", char, shortLink)
			}
		}
	}
}
