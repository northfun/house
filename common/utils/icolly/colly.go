package icolly

import (
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/redisstorage"
)

func BatchInitCollector(
	storage *redisstorage.Storage,
	userAgent string,
	limit *colly.LimitRule,
	onRequest func(r *colly.Request),
	clts ...**colly.Collector) (err error) {

	for i := range clts {

		newC := colly.NewCollector(
			colly.UserAgent(userAgent),
		)

		(*clts[i]) = newC

		if err = (*clts[i]).SetStorage(storage); err != nil {
			return err
		}

		if limit != nil {
			(*clts[i]).Limit(limit)
		}

		(*clts[i]).SetRequestTimeout(time.Second * 30)

		if onRequest != nil {
			(*clts[i]).OnRequest(onRequest)
		}
	}
	return
}

func OnRequestAliPai(r *colly.Request) {
	r.Headers.Set("sec-fetch-dest", "document")
	r.Headers.Set("sec-fetch-mode", "navigate")
	r.Headers.Set("sec-fetch-site", "none")
	r.Headers.Set("sec-fetch-user", "?1")
	r.Headers.Set("upgrade-insecure-requests", "1")
	r.Headers.Set("authority", "zc-paimai.taobao.com")
	r.Headers.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed  -exchange;v=b3;q=0.9")
}
