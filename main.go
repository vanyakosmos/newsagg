package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

var (
	_               = godotenv.Load()
	botToken        = os.Getenv("BOT_TOKEN")
	targetChannel   = os.Getenv("TARGET_CHANNEL")
	bucketEndpoint  = os.Getenv("BUCKET_ENDPOINT")
	bucketAccessKey = os.Getenv("BUCKET_ACCESS_KEY")
	bucketSecretKey = os.Getenv("BUCKET_SECRET_KEY")
	bucketRegion    = os.Getenv("BUCKET_REGION")
	bucketName      = getEnv("BUCKET_NAME", "newsagg")
)

func main() {
	sentryInit()
	defer sentryFlush()
	defer sentryRecover()

	ctx := context.Background()

	b, err := bot.New(botToken)
	if err != nil {
		panic(err)
	}
	user, _ := b.GetMe(ctx)
	log.Printf("BOT: id=%d username=%s\n", user.ID, user.Username)
	// tracker := NewFileTracker()
	tracker := NewBucketTracker(ctx, bucketEndpoint, bucketAccessKey, bucketSecretKey, bucketRegion, bucketName)

	tracker.CleanupOldTrackers(ctx)
	articles := ReadHackerNews()
	for _, a := range articles {
		if a.Score < 100 {
			log.Println("Skipping low score article:", a)
		}
		if tracker.IsTracked(ctx, a.ID) {
			log.Println("Skipping tracked article:", a)
		}
		if sendArticle(ctx, b, a) {
			tracker.MarkAsTracked(ctx, a.ID)
		}
	}
}

func sendArticle(ctx context.Context, b *bot.Bot, article HackerNewArticle) bool {
	disableLinks := true
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:             targetChannel,
		ParseMode:          models.ParseModeMarkdown,
		Text:               formatMessage(article),
		ReplyMarkup:        formatKeyboard(article),
		LinkPreviewOptions: &models.LinkPreviewOptions{IsDisabled: &disableLinks},
	})
	if err != nil {
		log.Println("Bot error:", err)
		return false
	}
	log.Println("Published article:", article.ID, article.Title)
	return true
}

func formatMessage(article HackerNewArticle) string {
	duration := time.Since(article.CreatedAt)
	messageLines := []string{
		fmt.Sprintf("*%s* \\(Score %d\\+ in %d hours\\)",
			bot.EscapeMarkdown(article.Title),
			article.Score,
			int(duration.Hours())),
		"",
		"*Link*: " + bot.EscapeMarkdown(article.ArticleURL),
		"*Comments*: " + bot.EscapeMarkdown(article.CommentsURL),
	}
	return strings.Join(messageLines, "\n")
}

func formatKeyboard(article HackerNewArticle) *models.InlineKeyboardMarkup {
	keyboardRow := make([]models.InlineKeyboardButton, 0)
	keyboardRow = append(keyboardRow, models.InlineKeyboardButton{Text: "Read", URL: article.ArticleURL})
	if article.CommentsURL != "" {
		keyboardRow = append(keyboardRow, models.InlineKeyboardButton{Text: fmt.Sprintf("%d+ Comments", article.CommentsNumber), URL: article.CommentsURL})
	}
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{keyboardRow},
	}
}
