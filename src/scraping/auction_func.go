package scraping

import (
	"encoding/json"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/northfun/house/common/typedef/tbtype"
	"github.com/northfun/house/common/utils"
	"github.com/northfun/house/common/utils/logger"
	"github.com/northfun/house/src/dao"
	"go.uber.org/zap"
)

func filterAuctionUrl(raw string) string {
	u := strings.ReplaceAll(raw, " ", "")
	u = strings.Replace(u, "\t", "", -1)
	if len(u) == 0 {
		return ""
	}

	if strings.HasPrefix(u, "https") {
		return u
	}

	return "https:" + u
}

func (am *AuctionManager) DoViewer(
	viewerQ, detailQ *queue.Queue) {

	am.clc.ViewerC.OnHTML("#sf-item-list-data", func(e *colly.HTMLElement) {
		// mainC.OnHTML("#J_ImgBooth", func(e *colly.HTMLElement) {
		var viewer tbtype.JsonAuctionReview
		if err := json.Unmarshal(
			[]byte(e.Text),
			&viewer); err != nil {
			logger.Warn("[scraping],unmarshal auction review", zap.Error(err))
			return
		}

		for i := range viewer.Data {
			// nextUrl := filterAuctionUrl(
			// 	viewer.Data[i].ItemUrl)

			// if len(nextUrl) == 0 {
			// 	return
			// }

			// am.clc.DetailC.SetCookies(nextUrl, am.clc.Ck)
			// detailQ.AddURL(nextUrl)

			// logger.Info("[scraping],auction,detail add url",
			//	zap.String("next", nextUrl))

			viewer.Data[i].IStart = utils.Int64Ms2Time(viewer.Data[i].Start)
			viewer.Data[i].IEnd = utils.Int64Ms2Time(viewer.Data[i].End)
		}

		dao.SaveAuctionReview(viewer.Data)
	})

	am.clc.ViewerC.OnHTML("a.next", func(e *colly.HTMLElement) {
		// mainC.OnHTML("#J_ImgBooth", func(e *colly.HTMLElement) {
		nextUrl := filterAuctionUrl(
			e.Attr("href"))

		if len(nextUrl) == 0 {
			return
		}

		am.clc.ViewerC.SetCookies(nextUrl, am.clc.Ck)
		viewerQ.AddURL(nextUrl)

		logger.Info("[scraping],auction,viewer add url",
			zap.String("next", nextUrl))
	})
}

func (am *AuctionManager) DoDetail() {
}
