package main

import (
	"log"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type HackerNewArticle struct {
	Title          string
	ArticleURL     string
	CommentsURL    string
	Score          int
	CommentsNumber int
	Age            string
}

const ROOT_URL = "https://news.ycombinator.com"

func readHackerNews() []HackerNewArticle {
	resp, err := http.Get(ROOT_URL)
	if err != nil {
		log.Println("error:", err)
		return nil
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Println("error:", err)
		return nil
	}
	articles := make([]HackerNewArticle, 0)

	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.Data == "tr" && hasClass(n, "submission") {
			article := extractCoreArticle(n)
			articles = append(articles, article)
		} else if n.Type == html.ElementNode && n.Data == "td" && hasClass(n, "subtext") {
			article := extractMetaArticle(n)
			articles = append(articles, article)
		}
	}
	return mergeArticles(articles)
}

func getAttr(node *html.Node, key string) (string, bool) {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}

func hasClass(node *html.Node, class string) bool {
	classesStr, found := getAttr(node, "class")
	if !found {
		return false
	}
	classes := strings.Split(classesStr, " ")
	return slices.Contains(classes, class)
}

func extractCoreArticle(node *html.Node) HackerNewArticle {
	article := HackerNewArticle{}
	for n := range node.Descendants() {
		if n.Type == html.ElementNode && n.Data == "span" && hasClass(n, "titleline") {
			anchor := n.FirstChild
			href, _ := getAttr(n.FirstChild, "href")
			article.ArticleURL = href

			text := anchor.FirstChild
			article.Title = text.Data
		}
	}
	return article
}

func extractMetaArticle(node *html.Node) HackerNewArticle {
	article := HackerNewArticle{}
	for n := range node.Descendants() {
		if n.Data == "span" && hasClass(n, "score") {
			text := n.FirstChild.Data
			scoreMatches := regexp.MustCompile("([0-9]+) points").FindStringSubmatch(text)
			score, _ := strconv.Atoi(scoreMatches[1])
			article.Score = score
		} else if n.Data == "span" && hasClass(n, "age") {
			title, _ := getAttr(n, "title")
			article.Age = title
		} else if n.Data == "a" {
			href, _ := getAttr(n, "href")
			text := n.FirstChild.Data
			text, _ = url.QueryUnescape(text)
			if !strings.HasPrefix(href, "item?id=") || !strings.HasSuffix(text, "comments") {
				continue
			}
			article.CommentsURL = ROOT_URL + "/" + href
			commentsMatches := regexp.MustCompile(`(\d+)[\s ]comments`).FindStringSubmatch(text)
			comments, _ := strconv.Atoi(commentsMatches[1])
			article.CommentsNumber = comments
		}
	}
	return article
}

func mergeArticles(articles []HackerNewArticle) []HackerNewArticle {
	newArticles := make([]HackerNewArticle, len(articles)/2)
	for i := 0; i < len(articles); i += 2 {
		core := articles[i]
		meta := articles[i+1]
		newArticles[i/2] = HackerNewArticle{
			Title:          core.Title,
			ArticleURL:     core.ArticleURL,
			Score:          meta.Score,
			CommentsURL:    meta.CommentsURL,
			CommentsNumber: meta.CommentsNumber,
			Age:            meta.Age,
		}
	}
	return newArticles
}
