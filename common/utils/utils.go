package utils

import (
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
