package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gocolly/colly"
)

var visited = map[string]bool{}

func main() {
	// Instanitate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("www.xinhuanet.com"),
		colly.MaxDepth(1),
	)
	// regular expression for detail page url
	detailRegex, _ := regexp.Compile(`c_\d{10}\.htm`)
	// regular expression for list page
	// 暂不要列表页
	fmt.Println("Start a test for regular regression")
	url := "http://www.xinhuanet.com/politics/2021-02/07/c_1127077568.htm"
	match := detailRegex.Match([]byte(url))
	fmt.Println("matched?: ", match)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// skip pages already visited
		if visited[link] && detailRegex.Match([]byte(link)) {
			return
		}
		if !detailRegex.Match([]byte(link)) {
			println("not match(detail page)", link)
			return
		}
		// sleep logic
		time.Sleep(time.Second)
		println("Matched link: ", link)

		visited[link] = true

		time.Sleep(time.Millisecond * 2)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	err := c.Visit("http://www.xinhuanet.com")
	if err != nil {
		fmt.Println("Got error", err)
	}
}
