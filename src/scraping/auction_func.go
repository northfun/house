package scraping

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/northfun/house/common/typedef/tbtype"
	"github.com/northfun/house/common/utils"
	"github.com/northfun/house/common/utils/ihttp"
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

func (la *LoadAuctionManager) LoadDoDetail() {

	var tb tbtype.TableSubjectMatterInfo

	la.clc.DetailC.OnHTML("#J_HoverShow",
		func(e *colly.HTMLElement) {
			e.ForEach("td", func(_ int, tdE *colly.HTMLElement) {
				var name string
				tdE.ForEach("span.pay-mark", func(_ int, nE *colly.HTMLElement) {
					name = nE.Text

				})

				tdE.ForEach("span.pay-price", func(_ int, pE *colly.HTMLElement) {
					price := utils.CovToPrice(pE.Text)

					if strings.Contains(name, "保证金") {
						tb.GuaranteeDeposit = price
					} else if strings.Contains(name, "起拍价") {
						tb.InitialPrice = price
					} else if strings.Contains(name, "评估价") {
						tb.ConsultPrice = price
					} else if strings.Contains(name, "加价") {
						tb.PriceRaise = price
					} else if strings.Contains(name, "优先购买") {
						tb.PayFirst = pE.Text
					}

				})

				tdE.ForEach("span.prior-td", func(_ int, pE *colly.HTMLElement) {
					if strings.Contains(pE.Text, "无") {
						tb.PayFirst = "无"
					} else if strings.Contains(pE.Text, "有") {
						tb.PayFirst = "有"
					} else {
						tb.PayFirst = utils.Trim(pE.Text)
					}
				})

			})
		})

	la.clc.DetailC.OnHTML("#J_NotifyNum", func(e *colly.HTMLElement) {
		var err error
		tb.AlarmTimes, err = utils.ParseUint32(e.Text)
		if err != nil {
			logger.Warn("[scraping],parse uint32", zap.String("text", e.Text), zap.Error(err))
		}
	})

	la.clc.DetailC.OnHTML("#J_desc", func(e *colly.HTMLElement) {
		dataFrom := filterAuctionUrl(e.Attr("data-from"))

		data, err := ihttp.Get(dataFrom)
		if err != nil {
			logger.Warn("[scraping],LoadDoDetail",
				zap.Error(err))
			return
		}

		retMap, err := ihttp.DealSubjectMatterTable(data)
		if err != nil {
			return
		}
		fmt.Println(retMap)

		utils.StructByReflect(retMap, &tb)

		tb.Raw = fmt.Sprintf("%v", retMap)
	})

	la.sp.InsertSubjectMatter(&tb)
}
