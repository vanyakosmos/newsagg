package internal

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func SendArticle(ctx context.Context, b *bot.Bot, article newsArticle, targetChannel string) bool {
	disableLinks := false
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      targetChannel,
		ParseMode:   models.ParseModeMarkdown,
		Text:        formatMessage(article),
		ReplyMarkup: formatKeyboard(article),
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: &disableLinks,
			URL:        &article.ArticleURL,
		},
	})
	if err != nil {
		if rateErr, ok := err.(*bot.TooManyRequestsError); ok {
			log.Printf("Too many telegram requests. Retrying after %d seconds\n", rateErr.RetryAfter)
			time.Sleep(time.Second * time.Duration(rateErr.RetryAfter))
			return SendArticle(ctx, b, article, targetChannel)
		}
		log.Println("Bot error:", err)
		return false
	}
	log.Println("Published article:", article)
	return true
}

func formatMessage(article newsArticle) string {
	duration := time.Since(article.CreatedAt)
	messageLines := []string{
		fmt.Sprintf("*%s* \\(Score %d\\+ in %s\\)",
			bot.EscapeMarkdown(article.Title),
			article.Score,
			bot.EscapeMarkdown(timeSincePost(duration)),
		),
		"",
		"*Link*: " + bot.EscapeMarkdown(article.ArticleURL),
		"*Comments*: " + bot.EscapeMarkdown(article.CommentsURL),
	}
	return strings.Join(messageLines, "\n")
}

func formatKeyboard(article newsArticle) *models.InlineKeyboardMarkup {
	keyboardRow := make([]models.InlineKeyboardButton, 0)
	keyboardRow = append(keyboardRow, models.InlineKeyboardButton{Text: "Read", URL: article.ArticleURL})
	if article.CommentsURL != "" {
		keyboardRow = append(keyboardRow, models.InlineKeyboardButton{Text: fmt.Sprintf("%d+ Comments", article.CommentsNumber), URL: article.CommentsURL})
	}
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{keyboardRow},
	}
}

func timeSincePost(duration time.Duration) string {
	if duration.Minutes() < 2 {
		return "1 minute 🔥"
	}
	if duration.Minutes() < 60 {
		return fmt.Sprintf("%d minutes 🔥", int(duration.Minutes()))
	}
	if duration.Hours() < 2 {
		return "1 hour 🔥"
	}
	if duration.Hours() < 24 {
		return fmt.Sprintf("%d hours", int(duration.Hours()))
	}
	if duration.Hours() < 48 {
		return "1 day"
	}
	return fmt.Sprintf("%d days", int(duration.Hours()/24))
}
