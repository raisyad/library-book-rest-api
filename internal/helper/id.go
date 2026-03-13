package helper

import "strconv"

func ParseIDParam(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}