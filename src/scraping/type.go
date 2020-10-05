package scraping

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/northfun/house/common/utils/logger"
	"go.uber.org/zap"
)

type StPage struct {
	TotalPage int `json:"totalPage"`
	CurPage   int `json:"curPage"`
}

func GenStPage(str string) (sp StPage) {
	if err := json.Unmarshal([]byte(str), &sp); err != nil {
		logger.Error("[scraping],page unmarshal",
			zap.String("str", str),
			zap.Error(err))
	}

	return
}

func GenPageUrl(dName, pageUrl, pageData string) (url string) {
	// <div class="page-box house-lst-page-box" comp-module='page' page-url="/chengjiao/pg{page}"page-data='{"totalPage":100,"curPage":1}'></div>
	sp := GenStPage(pageData)
	if sp.TotalPage == 0 {
		return
	}

	if sp.CurPage == sp.TotalPage {
		return
	}

	url = dName + strings.Replace(
		pageUrl, "{page}", fmt.Sprintf("%d", sp.CurPage+1), -1)
	return
}
