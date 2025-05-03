// O(V + E)，其中 V 是节点数（URL 数量），E 是边数（链接数量）

package main

import (
	"strings"
)

type HtmlParser interface {
	GetUrls(url string) []string
}

func crawl(startUrl string, htmlParser HtmlParser) []string {
	visited := make(map[string]struct{})
	queue := []string{}
	res := []string{}

	targetHostName := getHostName(startUrl)

	push := func(url string) {
		visited[url] = struct{}{}
		queue = append(queue, url)
		res = append(res, url)
	}
	push(startUrl)

	for len(queue) > 0 {
		url := queue[0]
		queue = queue[1:]
		for _, nextUrl := range htmlParser.GetUrls(url) {
			if _, has := visited[nextUrl]; !has && getHostName(nextUrl) == targetHostName {
				push(nextUrl)
			}
		}
	}

	return res
}

func getHostName(url string) string {
	return strings.Split(url, "/")[2]
}
