package ihttp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"github.com/northfun/house/common/utils"
)

func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func DealSubjectMatterTable(data []byte) (map[string]string, error) {
	ret := make(map[string]string)
	str := string(data)

	str = strings.Replace(str, "var desc='", "", 1)
	str = strings.Replace(str, "';", "", -1)

	str = mahonia.NewDecoder("gbk").ConvertString(str)

	doc, err := goquery.NewDocumentFromReader(
		bytes.NewReader([]byte(str)))
	if err != nil {
		return nil, err
	}

	doc.Find("table tr").Each(func(i int, s *goquery.Selection) {
		var content []string

		s.Find("td").Each(func(j int, is *goquery.Selection) {
			content = append(content,
				utils.Trim(is.Text()))
		})

		lenContent := len(content)
		if lenContent > 1 {
			ret[content[lenContent-2]] =
				content[lenContent-1]
		} else if lenContent == 1 {
			if strings.Contains(content[0], "介绍") {
				ret["介绍"] = content[0]
			} else {
				ret[fmt.Sprintf("%d", i)] = content[0]
			}
			// logger.Debug("[subject mater],content", zap.String("raw", content[0]))
		}
	})
	return ret, nil
}
