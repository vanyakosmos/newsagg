package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	TRACKER_DIR              = ".trackers"
	CLEANUP_DURATION_SECONDS = 48 * time.Hour
)

func IsTracked(article HackerNewArticle) bool {
	// setup
	os.Mkdir(TRACKER_DIR, 0755)
	cleanupOldTrackers()
	// check
	filename := getFilename(article.ID)
	_, err := os.Stat(filename)
	exists := err == nil
	return exists
}

func MarkAsTracked(article HackerNewArticle) {
	filename := getFilename(article.ID)
	file, _ := os.Create(filename)
	file.Close()
	log.Println("Tracked new article:", filename)
}

func cleanupOldTrackers() {
	entries, _ := os.ReadDir(TRACKER_DIR)
	for _, entry := range entries {
		info, _ := entry.Info()
		if time.Since(info.ModTime()) > CLEANUP_DURATION_SECONDS {
			filename := fmt.Sprintf("%s/%s", TRACKER_DIR, entry.Name())
			os.Remove(filename)
			log.Println("Cleaned up old tracker file:", filename)
		}
	}
}

func getFilename(id string) string {
	return fmt.Sprintf("%s/hn_%s.txt", TRACKER_DIR, id)
}
