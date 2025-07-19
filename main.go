package main

import (
	"context"
	"log"
	"newsagg/internal"
	"os"

	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
)

var (
	_                = godotenv.Load()
	telegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	targetChannel    = os.Getenv("TARGET_CHANNEL")
	bucketEndpoint   = os.Getenv("BUCKET_ENDPOINT")
	bucketAccessKey  = os.Getenv("BUCKET_ACCESS_KEY")
	bucketSecretKey  = os.Getenv("BUCKET_SECRET_KEY")
	bucketRegion     = os.Getenv("BUCKET_REGION")
	bucketName       = os.Getenv("BUCKET_NAME")
)

func main() {
	internal.SentryInit()
	defer internal.SentryFlush()
	defer internal.SentryRecover()

	ctx := context.Background()

	b, err := bot.New(telegramBotToken)
	if err != nil {
		panic(err)
	}
	user, _ := b.GetMe(ctx)
	log.Printf("BOT: id=%d username=%s\n", user.ID, user.Username)

	tracker := internal.NewBucketTracker(ctx, bucketEndpoint, bucketAccessKey, bucketSecretKey, bucketRegion, bucketName)
	tracker.CleanupOldTrackers(ctx)

	articles := internal.ReadHackerNews()
	articles = append(articles, internal.ReadLobsters()...)

	for _, a := range articles {
		if a.Source == internal.HackerNewsSource && a.Score < 100 {
			log.Println("Skipping low score HN article:", a)
		} else if a.Source == internal.LobstersSource && a.Score < 20 {
			log.Println("Skipping low score Lobsters article:", a)
		} else if tracker.IsTracked(ctx, a) {
			log.Println("Skipping tracked article:", a)
		} else if internal.SendArticle(ctx, b, a, targetChannel) {
			tracker.MarkAsTracked(ctx, a)
		}
	}
}
