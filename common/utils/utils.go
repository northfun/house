package utils

import (
	"net/http"
	"net/url"
	"strconv"
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
