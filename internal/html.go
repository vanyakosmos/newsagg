package internal

import (
	"log"
	"net/http"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

func loadPage(url string) *html.Node {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("error:", err)
		return nil
	}
	defer resp.Body.Close()
	log.Println("loaded:", url)

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Println("error:", err)
		return nil
	}
	return doc
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

func findChild(node *html.Node, tag string) *html.Node {
	for n := range node.Descendants() {
		if n.Data == tag {
			return n
		}
	}
	return nil
}
