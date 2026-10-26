package invezgo

import (
	"bytes"
	"strconv"
	"strings"
)

// FlexFloat is a number the API sometimes sends as a numeric string; both decode to float64.
type FlexFloat float64

// UnmarshalJSON accepts 1.5, "1.5", "" and null.
func (f *FlexFloat) UnmarshalJSON(b []byte) error {
	s, ok := flexText(b)
	if !ok {
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*f = FlexFloat(v)
	return nil
}

// FlexInt is an integer the API sometimes sends as a numeric string; both decode to int64.
type FlexInt int64

// UnmarshalJSON accepts 12, "12", "" and null.
func (i *FlexInt) UnmarshalJSON(b []byte) error {
	s, ok := flexText(b)
	if !ok {
		return nil
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fv, ferr := strconv.ParseFloat(s, 64)
		if ferr != nil {
			return err
		}
		v = int64(fv)
	}
	*i = FlexInt(v)
	return nil
}

func flexText(b []byte) (string, bool) {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || string(b) == "null" {
		return "", false
	}
	s := strings.TrimSpace(strings.Trim(string(b), `"`))
	return s, s != ""
}
