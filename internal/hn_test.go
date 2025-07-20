package internal

import (
	"testing"
)

func TestReadHackerNews(t *testing.T) {
	articles := ReadHackerNews()
	for _, a := range articles {
		t.Log(a)
	}
	if len(articles) != 150 {
		t.Errorf("expected 150 articles, got %d", len(articles))
	}
}
