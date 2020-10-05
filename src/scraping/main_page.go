package scraping

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/northfun/house/common/utils/logger"
	"go.uber.org/zap"
)

func (m *Manager) MainPage(mainQ, innerQ *queue.Queue) {
	m.mainC.OnError(func(resp *colly.Response, err error) {
		logger.Error("[scraping],main,on err",
			zap.String("url", resp.Request.ProxyURL), zap.Error(err))
	})

	m.mainC.OnHTML("body", func(e *colly.HTMLElement) {
		// fmt.Println("======url", string(e.Response.Body))

		e.DOM.Find("ul.listContent a.img").Each(
			func(_ int, s *goquery.Selection) {
				detailUrl := s.AttrOr("href", "")
				fmt.Println("d==", detailUrl)

				if len(detailUrl) == 0 {
					return
				}
				err := innerQ.AddURL(detailUrl)

				logger.Info("[scrapping],detailUrl", zap.String("url", detailUrl), zap.Error(err))
			})

		pageBox := e.DOM.Find("div.page-box.house-lst-page-box")
		nextUrl := GenPageUrl(
			DOMAIN_NAME,
			pageBox.AttrOr("page-url", ""),
			pageBox.AttrOr("page-data", ""))

		fmt.Println("m==", nextUrl)

		if len(nextUrl) == 0 {
			return
		}

		err := mainQ.AddURL(nextUrl)

		logger.Info("[scrapping],nextMainUrl", zap.String("url", nextUrl), zap.Error(err))
	})
}

func (m *Manager) InnerPage(moreQ *queue.Queue) {
	m.innerC.OnError(func(resp *colly.Response, err error) {
		logger.Error("[scraping],inner,on err",
			zap.String("url", resp.Request.ProxyURL), zap.Error(err))
	})

	m.innerC.OnHTML("body", func(e *colly.HTMLElement) {
		// fmt.Println("======url", string(e.Response.Body))

		// moreUri := e.DOM.Find("a.getMoreHouse").
		// 	AttrOr("href", "")

		communityID := e.DOM.Find("div.house-title.LOGVIEWDATA.LOGVIEW").AttrOr("data-lj_action_housedel_id", "")
		// fmt.Println("=com=========", communityID)
		if len(communityID) == 0 {
			return
		}

		moreUri := "/chengjiao/c" + communityID
		nextUrl := DOMAIN_NAME + moreUri
		if len(nextUrl) == 0 {
			return
		}

		err := moreQ.AddURL(nextUrl)

		logger.Info("[scrapping],nextMoreUrl", zap.String("url", nextUrl), zap.Error(err))
	})
}

func (m *Manager) MorePage(moreQ, detailQ *queue.Queue) {
	m.moreC.OnError(func(resp *colly.Response, err error) {
		logger.Error("[scraping],main,on err",
			zap.String("url", resp.Request.ProxyURL), zap.Error(err))
	})

	m.moreC.OnHTML("body", func(e *colly.HTMLElement) {
		// fmt.Println("======url", string(e.Response.Body))

		e.DOM.Find("ul.listContent a.img").Each(
			func(_ int, s *goquery.Selection) {
				detailUrl := s.AttrOr("href", "")
				fmt.Println("d==", detailUrl)

				if len(detailUrl) == 0 {
					return
				}
				err := detailQ.AddURL(detailUrl)

				logger.Info("[scrapping],detailUrl", zap.String("url", detailUrl), zap.Error(err))
			})

		pageBox := e.DOM.Find("div.page-box.house-lst-page-box")
		nextUrl := GenPageUrl(
			DOMAIN_NAME,
			pageBox.AttrOr("page-url", ""),
			pageBox.AttrOr("page-data", ""))

		fmt.Println("mo==", nextUrl)

		if len(nextUrl) == 0 {
			return
		}

		err := moreQ.AddURL(nextUrl)

		logger.Info("[scrapping],nextMoreUrl", zap.String("url", nextUrl), zap.Error(err))
	})
}
