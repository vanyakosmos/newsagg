package internal

import (
	"testing"
	"time"
)

func TestTimeSincePost(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{duration: time.Second * 10, expected: "1 minute 🔥"},
		{duration: time.Minute * 1, expected: "1 minute 🔥"},
		{duration: time.Minute * 10, expected: "10 minutes 🔥"},
		{duration: time.Minute * 70, expected: "1 hour 🔥"},
		{duration: time.Minute * 120, expected: "2 hours"},
		{duration: time.Hour * 25, expected: "1 day"},
		{duration: time.Hour * 49, expected: "2 days"},
		{duration: time.Hour * 72, expected: "3 days"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			actual := timeSincePost(test.duration)
			if actual != test.expected {
				t.Errorf("%s != %s", actual, test.expected)
			}
		})
	}
}

func TestFormatMessage(t *testing.T) {
	article := newsArticle{
		ID:          "id",
		Title:       "title",
		ArticleURL:  "aurl",
		CommentsURL: "curl",
		Score:       123,
		CreatedAt:   time.Now(),
	}
	actual := formatMessage(article)
	if actual != "*title* \\(Score 123\\+ in 1 minute 🔥\\)\n\n*Link*: aurl\n*Comments*: curl" {
		t.Errorf("message is not right:\n%s", actual)
	}
}
