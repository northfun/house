package utils

import (
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/northfun/house/common/utils/logger"
	"go.uber.org/zap"
)

func ParseUint32(str string) (uint32, error) {
	iuint, err := strconv.Atoi(str)
	return uint32(iuint), err
}

func LastPart2Uint64(url string) uint64 {
	idx := strings.LastIndex(url, "=")
	if idx == 0 {
		return 0
	}

	ret, err := strconv.Atoi(url[idx+1:])
	if err != nil {
		logger.Warn("[utils],last part 2 uint64",
			zap.String("url", url[idx+1:]), zap.Error(err))
		return 0
	}
	return uint64(ret)
}

func ParseUrlQuery(q, name string) string {
	return url.Values{q: {name}}.Encode()
}

func GenCookies(rawCookies string) []*http.Cookie {
	// rawRequest := fmt.Sprintf("GET / HTTP/1.0\r\nCookie: %s\r\n\r\n", rawCookies)

	// req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(rawRequest)))
	header := http.Header{}
	header.Add("Cookie", rawCookies)
	request := http.Request{
		Header: header,
	}

	return request.Cookies()
}

func Int64Ms2Time(itm int64) (tm time.Time) {
	tm = time.Unix(itm/1000, 0)
	return
}

func TrimForNum(str string) string {
	str = Trim(str)

	str = strings.ReplaceAll(str, "¥", "")
	str = strings.ReplaceAll(str, ":", "")

	return str
}

func Trim(str string) string {
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "\u00a0", "")
	str = strings.ReplaceAll(str, "\u00A0", "")
	str = strings.ReplaceAll(str, "\u0020", "")
	// str = strings.ReplaceAll(str, "\u0030", "")
	str = strings.ReplaceAll(str, "&nbsp;", "")
	str = strings.ReplaceAll(str, "\n", "")

	return str
}

func StructByReflect(data map[string]string, inStructPtr interface{}) {
	rType := reflect.TypeOf(inStructPtr)
	rVal := reflect.ValueOf(inStructPtr)
	if rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
		rVal = rVal.Elem()
	} else {
		panic("inStructPtr must be ptr to struct")
	}

	for i := 0; i < rType.NumField(); i++ {
		tagName := rType.Field(i).Tag.Get("cname")
		if len(tagName) == 0 {
			continue
		}
		f := rVal.Field(i)

		for cName, v := range data {
			if strings.Contains(cName, tagName) {
				// f.Set(reflect.ValueOf(v))
				f.SetString(v)
				break
			}
		}
	}
}

func CovToPrice(str string) float64 {
	str = strings.ReplaceAll(str, "\"", "")
	str = strings.ReplaceAll(str, ",", "")
	str = TrimForNum(str)

	f, err := strconv.ParseFloat(str, 64)
	if len(str) == 0 {
		return 0
	}
	if err != nil {
		logger.Warn("[utils],parse float", zap.String("str", str), zap.Error(err))
	}
	return f
}

func ExtraHousePropertyRight(name string) (string, string) {
	slc := strings.Split(name, "（")
	if len(slc) <= 1 {
		return name, ""
	}

	return slc[0], slc[1]
}
