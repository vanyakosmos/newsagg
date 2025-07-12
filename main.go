package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

var _ = godotenv.Load()
var (
	botToken      = os.Getenv("BOT_TOKEN")
	targetChannel = os.Getenv("TARGET_CHANNEL")
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	b, err := bot.New(botToken)
	if err != nil {
		panic(err)
	}
	user, _ := b.GetMe(ctx)
	fmt.Printf("BOT: id=%d username=%s\n", user.ID, user.Username)

	for {
		// collectNews(ctx, b)
		// time.Sleep(time.Second * 5)
		readHackerNews()
		return
	}
}

func collectNews(ctx context.Context, b *bot.Bot) {
	messageLines := []string{
		"*Some new title* " + bot.EscapeMarkdown("(Score 150+ in 11 hours)"),
		"",
		"*Link*: " + bot.EscapeMarkdown("https://example.com"),
		"*Comments*: " + bot.EscapeMarkdown("https://example.com"),
		"",
		bot.EscapeMarkdown("extra text"),
	}
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Read", URL: "https://example.com"},
				{Text: "76+ Comments", URL: "https://example.com"},
			},
		},
	}
	disableLinks := true
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:             targetChannel,
		Text:               strings.Join(messageLines, "\n"),
		ParseMode:          models.ParseModeMarkdown,
		ReplyMarkup:        kb,
		LinkPreviewOptions: &models.LinkPreviewOptions{IsDisabled: &disableLinks},
	})
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}
