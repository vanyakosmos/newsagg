package internal

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const hackerNewRootURL = "https://news.ycombinator.com"
const HackerNewsSource = "hn"

func ReadHackerNews() []newsArticle {
	articles := make([]newsArticle, 0)

	for page := range loadHackerNewsPages() {
		for n := range page.Descendants() {
			if n.Type == html.ElementNode && n.Data == "tr" && hasClass(n, "submission") {
				article := extractCoreArticle(n)
				articles = append(articles, article)
			} else if n.Type == html.ElementNode && n.Data == "td" && hasClass(n, "subtext") {
				article := extractMetaArticle(n)
				articles = append(articles, article)
			}
		}
	}
	return mergeArticles(articles)
}

func loadHackerNewsPages() <-chan *html.Node {
	pages := make(chan *html.Node)
	go func() {
		for i := range 5 {
			query := fmt.Sprintf("%s?p=%d", hackerNewRootURL, i+1)
			doc := loadPage(query)
			if doc != nil {
				pages <- doc
			}
		}
		close(pages)
	}()
	return pages
}

func extractCoreArticle(node *html.Node) newsArticle {
	article := newsArticle{Source: HackerNewsSource}

	id, _ := getAttr(node, "id")
	article.ID = id

	for n := range node.Descendants() {
		if n.Type == html.ElementNode && n.Data == "span" && hasClass(n, "titleline") {
			anchor := n.FirstChild
			href, _ := getAttr(n.FirstChild, "href")
			if strings.HasPrefix(href, "http") {
				article.ArticleURL = href
			} else {
				article.ArticleURL = hackerNewRootURL + "/" + href
			}
			text := anchor.FirstChild
			article.Title = text.Data
		}
	}
	return article
}

func extractMetaArticle(node *html.Node) newsArticle {
	article := newsArticle{Source: HackerNewsSource}
	for n := range node.Descendants() {
		if n.Data == "span" && hasClass(n, "score") {
			text := n.FirstChild.Data
			scoreMatches := regexp.MustCompile("([0-9]+) points").FindStringSubmatch(text)
			score, _ := strconv.Atoi(scoreMatches[1])
			article.Score = score
		} else if n.Data == "span" && hasClass(n, "age") {
			title, _ := getAttr(n, "title")
			title = strings.Split(title, " ")[0]
			age, _ := time.Parse("2006-01-02T15:04:05", title)
			article.CreatedAt = age
		} else if n.Data == "a" {
			href, _ := getAttr(n, "href")
			text := n.FirstChild.Data
			text, _ = url.QueryUnescape(text)
			if !strings.HasPrefix(href, "item?id=") || !strings.HasSuffix(text, "comments") {
				continue
			}
			article.CommentsURL = hackerNewRootURL + "/" + href
			commentsMatches := regexp.MustCompile(`(\d+)[\s ]comments`).FindStringSubmatch(text)
			comments, _ := strconv.Atoi(commentsMatches[1])
			article.CommentsNumber = comments
		}
	}
	return article
}

func mergeArticles(articles []newsArticle) []newsArticle {
	newArticles := make([]newsArticle, len(articles)/2)
	for i := 0; i < len(articles); i += 2 {
		core := articles[i]
		meta := articles[i+1]
		newArticles[i/2] = newsArticle{
			Source:         HackerNewsSource,
			ID:             core.ID,
			Title:          core.Title,
			ArticleURL:     core.ArticleURL,
			Score:          meta.Score,
			CommentsURL:    meta.CommentsURL,
			CommentsNumber: meta.CommentsNumber,
			CreatedAt:      meta.CreatedAt,
		}
	}
	return newArticles
}
