package internal

import (
	"fmt"
	"time"
)

type newsArticle struct {
	Source         string
	ID             string
	Title          string
	ArticleURL     string
	Score          int
	CommentsURL    string
	CommentsNumber int
	CreatedAt      time.Time
}

func (a newsArticle) String() string {
	return fmt.Sprintf("[%s/%s] 𝚫%d ~ %s", a.Source, a.ID, a.Score, a.Title)
}
