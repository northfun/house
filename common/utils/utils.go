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

func Trim(str string) string {
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "\u00a0", "")
	str = strings.ReplaceAll(str, "\u00A0", "")
	str = strings.ReplaceAll(str, "\u0020", "")
	str = strings.ReplaceAll(str, "\u0030", "")
	str = strings.ReplaceAll(str, "&nbsp;", "")
	str = strings.ReplaceAll(str, "\n", "")

	// fmt.Printf("====%q====", str)
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
	str = Trim(str)

	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		logger.Warn("[utils],parse float", zap.String("str", str), zap.Error(err))
	}
	return f
}
