package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

type tracker interface {
	IsTracked(ctx context.Context, article newsArticle) bool
	MarkAsTracked(ctx context.Context, article newsArticle)
	CleanupOldTrackers(ctx context.Context)
}

type fileTracker struct {
	rootDir      string
	fileAgeLimit time.Duration
}

func NewFileTracker() tracker {
	return &fileTracker{
		rootDir:      ".trackers",
		fileAgeLimit: 1 * time.Minute,
	}
}

func (t *fileTracker) IsTracked(ctx context.Context, article newsArticle) bool {
	// setup
	os.Mkdir(t.rootDir, 0755)
	t.CleanupOldTrackers(ctx)
	// check
	filename := t.getFilename(article)
	_, err := os.Stat(filename)
	exists := err == nil
	return exists
}

func (t *fileTracker) MarkAsTracked(ctx context.Context, article newsArticle) {
	filename := t.getFilename(article)
	file, _ := os.Create(filename)
	file.Close()
	log.Println("Tracked new article:", filename)
}

func (t *fileTracker) CleanupOldTrackers(ctx context.Context) {
	entries, _ := os.ReadDir(t.rootDir)
	for _, entry := range entries {
		info, _ := entry.Info()
		if time.Since(info.ModTime()) > t.fileAgeLimit {
			filename := fmt.Sprintf("%s/%s", t.rootDir, entry.Name())
			os.Remove(filename)
			log.Println("Cleaned up old tracker file:", filename)
		}
	}
}

func (t *fileTracker) getFilename(a newsArticle) string {
	return fmt.Sprintf("%s/%s_%s.txt", t.rootDir, a.Source, a.ID)
}
