package scraping

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/northfun/house/common/typedef/tbtype"
)

func (m *Manager) DetailPage(detailQ *queue.Queue) {

	var houseInfo tbtype.TableHouseDealInfo

	m.detailC.OnHTML("div.introContent div.content", func(e *colly.HTMLElement) {

		e.DOM.Find("ul li").Each(func(_ int, s *goquery.Selection) {
			content := s.Text()

			label := s.Find("span.label").Text()
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
				houseInfo.HouseType = content
			case "梯户比例":
				houseInfo.StairResident = content
			case "配备电梯":
				houseInfo.Elevator = content
			}
		})
	})

	m.detailC.OnHTML("div.transaction div.content", func(e *colly.HTMLElement) {

		e.DOM.Find("ul li").Each(func(_ int, s *goquery.Selection) {
			content := s.Text()

			label := s.Find("span.label").Text()
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
			}
		})
	})

	m.detailC.OnHTML("div.imgContainer img.defaultImg", func(e *colly.HTMLElement) {
		houseInfo.PicUrl = e.Attr("src")
	})

	m.detailC.OnHTML("div.house-title div.wrapper span", func(e *colly.HTMLElement) {
		houseInfo.SoldTime = e.Text
	})

	m.detailC.OnHTML("div.info.fr", func(e *colly.HTMLElement) {
		e.DOM.Find("div.price span.dealTotalPrice").Each(func(_ int, s *goquery.Selection) {
			houseInfo.Price = s.Find("i").Text() + s.Text()
		})

		e.DOM.Find("msg span").Each(
			func(_ int, s *goquery.Selection) {

				content := s.Find("label").Text()

				label := s.Text()

				if strings.Contains(label, "挂牌价格") {
					houseInfo.ListedPrice = content
				} else if strings.Contains(label, "调价") {
					iadj, _ := strconv.Atoi(content)
					houseInfo.AdjustPriceTimes = uint16(iadj)
				} else if strings.Contains(label, "带看") {
					ictt, _ := strconv.Atoi(content)
					houseInfo.VisitTimes = uint16(ictt)
				}
			})
	})

	m.detailC.OnHTML("div.agent div.fr", func(e *colly.HTMLElement) {
		houseInfo.Saler = e.DOM.Find("a b.LOGCLICK").Text()

		houseInfo.SalerPhone = strings.ReplaceAll(
			e.DOM.Find("div.tel").Text(), " ", "")
	})

	fmt.Println("========", houseInfo)

	m.sp.InsertHouse(&houseInfo)

}
