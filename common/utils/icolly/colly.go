package icolly

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/redisstorage"
)

func BatchInitCollector(
	storage *redisstorage.Storage,
	userAgent string,
	limit *colly.LimitRule,
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
	}

	return
}
