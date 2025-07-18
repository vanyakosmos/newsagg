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
		ChatID:             targetChannel,
		ParseMode:          models.ParseModeMarkdown,
		Text:               formatMessage(article),
		ReplyMarkup:        formatKeyboard(article),
		LinkPreviewOptions: &models.LinkPreviewOptions{IsDisabled: &disableLinks},
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
