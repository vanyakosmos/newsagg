package internal

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const lobsterRootUrl = "https://lobste.rs"
const LobstersSource = "lobsters"

func ReadLobsters() []newsArticle {
	articles := make([]newsArticle, 0)
	for page := range lobsetersPages() {
		for n := range page.Descendants() {
			if n.Type == html.ElementNode && n.Data == "div" && hasClass(n, "story_liner") {
				article := extractArticle(n)
				articles = append(articles, article)
			}
		}
	}
	return articles
}

func lobsetersPages() chan *html.Node {
	pages := make(chan *html.Node)
	go func() {
		for i := range 5 {
			query := fmt.Sprintf("%s?page=%d", lobsterRootUrl, i+1)
			doc := loadPage(query)
			if doc != nil {
				pages <- doc
			}
		}
		close(pages)
	}()
	return pages
}

func extractArticle(node *html.Node) newsArticle {
	article := newsArticle{Source: LobstersSource}
	for n := range node.Descendants() {
		if n.Data == "a" && hasClass(n, "u-url") {
			href, _ := getAttr(n, "href")
			article.ArticleURL = href
			article.Title = n.FirstChild.Data
		} else if n.Data == "span" && hasClass(n, "comments_label") {
			anchor := findChild(n, "a")
			href, _ := getAttr(anchor, "href")
			parts := strings.Split(href, "/")
			parts = parts[:3]
			article.ID = parts[2]
			article.CommentsURL = lobsterRootUrl + strings.Join(parts, "/")

			text := anchor.FirstChild.Data
			commentsMatches := regexp.MustCompile(`\s*(\d+)[\s ]comments?\s*`).FindStringSubmatch(text)
			if len(commentsMatches) > 0 {
				comments, _ := strconv.Atoi(commentsMatches[1])
				article.CommentsNumber = comments
			}
		} else if n.Data == "a" && hasClass(n, "upvoter") {
			text := n.FirstChild.Data
			score, _ := strconv.Atoi(text)
			article.Score = score
		} else if n.Data == "time" {
			title, _ := getAttr(n, "title")
			age, _ := time.Parse("2006-01-02 15:04:05 -0700", title)
			article.CreatedAt = age
		}
	}
	return article
}
