package main

import (
	"fmt"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
)

func main() {
	mainQ, _ := queue.New(
		2, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	urlSlc := []string{
		"http://www.biopurify.cn/search.do?a=s&searchtype=3&psize=10&q=%E7%BE%9F%E4%B8%99%E5%9F%BA%E5%9B%9B%E6%B0%A2%E5%90%A1%E5%96%83%E4%B8%89%E9%86%87&searchtmp=simplesearch&t=-1",
		// 无符合资料
		"http://www.cosdna.com/chs/product.php?q=%E5%A7%97%E5%A8%9C%E8%B1%86%E4%B9%B3%E7%BE%8E%E8%82%8C%E4%BF%9D%E6%B9%BF%E4%B9%B3%E6%B6%B2",
		"http://www.cosdna.com/chs/product.php?%s&sort=click",
	}

	userAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36"

	for i := range urlSlc {
		if i == 2 {
			urlSlc[i] = fmt.Sprintf(urlSlc[i], url.Values{"q": {"欧莱雅(L'OREAL PARIS,欧莱雅)"}}.Encode())
		}
		mainQ.AddURL(urlSlc[i])

		uri, err := url.ParseRequestURI(urlSlc[i])
		if err != nil {
			fmt.Println("parse url err", err)
			continue
		}

		if v, err := url.ParseQuery(uri.RawQuery); true {
			fmt.Println("=parse==欧莱雅(L'OREAL PARIS,欧莱雅)==", v["q"], err)
		}
	}

	c := url.Values{"q": {"羟丙基四氢吡喃三醇"}, "a": {"s"}}
	fmt.Println(c.Encode())

	mainC := colly.NewCollector(
		colly.UserAgent(userAgent),
	)

	mainC.OnHTML("body", func(e *colly.HTMLElement) {
		fmt.Println("======url", e.Request.URL)
		// fmt.Println("======body", e.DOM.Find(".main").Text())
	})

	// mainC.OnHTML("#pro_details p", func(e *colly.HTMLElement) {
	// 	fmt.Println(e.Text)
	// })

	mainQ.Run(mainC)
}
