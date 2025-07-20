package internal

import (
	"testing"
)

func TestReadLobsters(t *testing.T) {
	articles := ReadLobsters()
	for _, a := range articles {
		t.Log(a)
	}
	if len(articles) != 125 {
		t.Errorf("expected 125 articles, got %d", len(articles))
	}
}
