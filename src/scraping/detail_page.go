package scraping

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/northfun/house/common/typedef/tbtype"
	"github.com/northfun/house/common/utils/logger"
	"go.uber.org/zap"
)

func houseBase(_ int, s *goquery.Selection,
	houseInfo *tbtype.TableHouseDealInfo) {

	content := strings.Trim(s.Text(), " ")

	label := strings.ReplaceAll(
		s.Find("span.label").Text(), " ", "")

	content = strings.Replace(content, label, "", 1)
	switch label {
	case "房屋户型":
		houseInfo.HouseType = content
	case "所在楼层":
		houseInfo.Floor = content
	case "建筑面积":
		houseInfo.Area = content
	case "套内面积":
		houseInfo.AllArea = content
	case "户型结构":
		houseInfo.RoomType = content
	case "建筑类型":
		houseInfo.BuildingType = content
	case "房屋朝向":
		houseInfo.RoomForward = content
	case "建成年代":
		houseInfo.BuildYearN = content
	case "装修情况":
		houseInfo.DecorationType = content
	case "建筑结构":
		houseInfo.HouseStructure = content
	case "供暖方式":
		houseInfo.HeatingMode = content
	case "梯户比例":
		houseInfo.StairResident = content
	case "配备电梯":
		houseInfo.Elevator = content

	default:
		logger.Error("[scraping],unknown housebase lable", zap.String("label", label))
	}
}

func houseBase2(_ int, s *goquery.Selection,
	houseInfo *tbtype.TableHouseDealInfo) {

	content := strings.Trim(s.Text(), " ")

	label := strings.ReplaceAll(
		s.Find("span.label").Text(), " ", "")

	content = strings.Replace(content, label, "", 1)
	switch label {
	case "链家编号":
		houseInfo.LianjianId = content
	case "交易权属":
		houseInfo.TxTypeN = content
	case "挂牌时间":
		houseInfo.SaleTime = content
	case "房屋用途":
		houseInfo.HouseUsingN = content
	case "房屋年限":
		houseInfo.HouseLimitYearN = content
	case "房权所属":
		houseInfo.HouseBelong = content

	default:
		logger.Error("[scraping],unknown housebase2 lable", zap.String("label", label))
	}
}

func soldInfo(_ int, s *goquery.Selection,
	houseInfo *tbtype.TableHouseDealInfo) {

	content := strings.ReplaceAll(
		s.Find("label").Text(), " ", "")

	label := strings.Trim(s.Text(), " ")

	content = strings.Replace(content, label, "", 1)

	if strings.Contains(label, "挂牌价格") {
		houseInfo.ListedPrice = content
	} else if strings.Contains(label, "调价") {
		iadj, _ := strconv.Atoi(content)
		houseInfo.AdjustPriceTimes = uint16(iadj)
	} else if strings.Contains(label, "带看") {
		ictt, _ := strconv.Atoi(content)
		houseInfo.VisitTimes = uint16(ictt)
	}
}

func getCommunityName(oriStr string) string {
	trimed := strings.Trim(oriStr, " ")

	if slc := strings.Split(trimed, " "); len(slc) > 1 {
		return slc[0]
	}
	return trimed
}

func (m *Manager) DetailPage() {

	m.detailC.OnHTML("body", func(e *colly.HTMLElement) {
		var houseInfo tbtype.TableHouseDealInfo

		houseInfo.Extra = e.Request.URL.RequestURI()

		houseInfo.CommunityName = getCommunityName(
			e.DOM.Find("h1.index_h1").Text())

		houseInfo.CommunityId = e.DOM.Find("div.house-title").
			AttrOr("data-lj_action_housedel_id", "")

		e.DOM.
			Find("div.introContent div.base div.content ul li").
			Each(func(_ int, s *goquery.Selection) {
				houseBase(0, s, &houseInfo)
			})

		e.DOM.
			Find("div.transaction div.content ul li").
			Each(func(_ int, s *goquery.Selection) {
				houseBase2(0, s, &houseInfo)
			})

		houseInfo.PicUrl = e.DOM.
			Find("div.imgContainer img.defaultImg").
			AttrOr("src", "")

		houseInfo.SoldTime = strings.
			Trim(e.DOM.
				Find("div.house-title div.wrapper span").
				Text(), " ")

		houseInfo.Price = strings.
			ReplaceAll(e.DOM.
				Find("div.price span.dealTotalPrice").
				Text(), " ", "")

		e.DOM.Find("div.info.fr span").Each(
			func(_ int, s *goquery.Selection) {
				soldInfo(0, s, &houseInfo)
			})

		houseInfo.Saler = e.DOM.Find("b.LOGCLICK").Text()

		houseInfo.SalerPhone = strings.ReplaceAll(
			e.DOM.Find("div.tel").Text(), " ", "")

		logger.Debug("[scraping],insert house",
			zap.String("url", e.Request.URL.String()),
			zap.Reflect("house", houseInfo))

		m.sp.InsertHouse(&houseInfo)
	})
}
