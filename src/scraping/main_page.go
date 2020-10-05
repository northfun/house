package scraping

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/northfun/house/common/utils/logger"
	"go.uber.org/zap"
)

func (m *Manager) MainPage(iurl []string, mainQ, detailQ *queue.Queue) {

	m.mainC.OnHTML("ul.listContent", func(e *colly.HTMLElement) {
		e.DOM.Find("a.img").Each(
			func(_ int, s *goquery.Selection) {
				detailUrl := s.AttrOr("href", "")
				fmt.Println("d==", detailUrl)

				if len(detailUrl) == 0 {
					return
				}
				detailQ.AddURL(detailUrl)

				logger.Info("[scrapping],detailUrl", zap.String("url", detailUrl))
			})
	})

	m.mainC.OnHTML("div.page-box.fr div.page-box.house-lst-page-box",
		func(e *colly.HTMLElement) {
			e.DOM.Find("a").Each(
				func(_ int, s *goquery.Selection) {
					if !strings.Contains(
						s.Text(), "下一页") {
						return
					}
					nextMainUri := s.AttrOr("href", "")
					fmt.Println("m==", nextMainUri)
					if len(nextMainUri) == 0 {
						return
					}

					nextMainUrl := ChengjiaoPage + nextMainUri
					mainQ.AddURL(nextMainUrl)
					logger.Info("[scrapping],nextMainUrl", zap.String("url", nextMainUrl))

				})
		})
}
