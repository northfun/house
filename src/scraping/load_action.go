package scraping

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/redisstorage"
	"github.com/northfun/house/common/typedef/tbtype"
	"github.com/northfun/house/common/utils"
	"github.com/northfun/house/common/utils/ihttp"
	"github.com/northfun/house/common/utils/logger"
	"github.com/northfun/house/src/conf"
	"github.com/northfun/house/src/dao"
	"github.com/northfun/house/src/sink"
	"go.uber.org/zap"
)

type LoadAuctionManager struct {
	domainName string
	fromUrls   []string

	sp *sink.SinkPool

	clc *ViewerDetailCollector
}

func NewLoadAuctionManager(
	sp *sink.SinkPool,
	storage *redisstorage.Storage) *LoadAuctionManager {

	rawCookies := conf.C().Cookie

	return &LoadAuctionManager{
		domainName: "",
		sp:         sp,
		clc: NewViewerDetailCollector(
			storage, rawCookies),

		fromUrls: []string{},
	}
}

func (la *LoadAuctionManager) Start() error {
	urls, err := dao.LoadUnflagedAuctionItems()
	if err != nil {
		return err
	}
	return la.clc.Start("load auction", nil,
		la.LoadDoDetail, urls)
}

func (la *LoadAuctionManager) Stop() {}

func (la *LoadAuctionManager) LoadDoDetail() {

	la.clc.DetailC.OnHTML("body",
		func(oe *colly.HTMLElement) {
			var tb tbtype.TableSubjectMatterInfo

			tb.Id = utils.LastPart2Uint64(oe.DOM.Find("#stream-url").AttrOr("value", ""))

			oe.DOM.Find("#J_HoverShow td").
				Each(func(_ int, tdE *goquery.Selection) {
					name := utils.Trim(tdE.Find("span.pay-mark").Text())

					priceText := tdE.Find("span.pay-price").Text()
					price := utils.CovToPrice(priceText)

					if strings.Contains(name, "保证金") {
						tb.GuaranteeDeposit = price
					} else if strings.Contains(name, "起拍价") {
						tb.InitialPrice = price
					} else if strings.Contains(name, "评估价") {
						tb.ConsultPrice = price
					} else if strings.Contains(name, "加价") {
						tb.PriceRaise = price
					} else if strings.Contains(name, "优先购买") {
						tb.PayFirst = priceText
					}

					priorText := tdE.Find("span.prior-td").Text()
					if strings.Contains(priorText, "无") {
						tb.PayFirst = "无"
					} else if strings.Contains(priorText, "有") {
						tb.PayFirst = "有"
					} else {
						tb.PayFirst = utils.Trim(priorText)
					}
				})

			alarmText := oe.DOM.Find("#J_NotifyNum").Text()
			var err error
			tb.AlarmTimes, err = utils.ParseUint32(alarmText)
			if err != nil {
				logger.Warn("[scraping],parse uint32", zap.String("text", alarmText), zap.Error(err))
			}

			eAttr := oe.DOM.Find("#J_desc").AttrOr("data-from", "")
			dataFrom := filterAuctionUrl(eAttr)

			if len(dataFrom) == 0 {
				logger.Warn("[scraping],filter auction url",
					zap.Uint64("id", tb.Id),
					zap.Error(err))
				return
			}

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

			if len(tb.HousePropertyRight) == 0 {
				_, tb.HousePropertyRight = utils.ExtraHousePropertyRight(tb.Name)
				if strings.Contains(tb.HousePropertyRight, "房产证号") {
					tb.HousePropertyRight = "有"
				}
			}

			la.sp.InsertSubjectMatter(&tb)
			logger.Info("[scraping],sm insert", zap.Uint64("id", tb.Id))
		})
}
