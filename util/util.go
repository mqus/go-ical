package util

import (
	"time"

	"errors"

	"strconv"
	"strings"

	"github.com/soh335/icalparser"
)

//it is assumed that all parts of a value in a parameter string are blindly separated by semicoli, even when they are in Quotes,
//this function reassembles such strings blindly.
/*func ReassembleParts(valueParts []*icalparser.Ident) string {
	out :=""
	first:=true
	for _,part :=range valueParts{
		if part==nil{
			continue
		}
		if first{
			out = part.
		}
	}
}

MAYBE do this... It has to go over parameter boundaries :/ */

func GetParam(list []*icalparser.Param, key string) int {
	for i, p := range list {
		if p.ParamName.C == key {
			return i
		}
	}
	return -1
}

const (
	ISO8601_2004_DATE = "20060102"
	ISO8601_2004_TIME = "150405"
	ISO8601_2004      = "20060102T150405"
)

//MAYBE theoretically not a good equivalent, 1 Week may not be exactly 7*24 hours if this week e.g. crosses over a dst-boundary
func ParseDuration(in string) (out time.Duration, err error) {
	negative := false
	parts := strings.SplitN(in, "P", 2)
	out = time.Duration(0)
	if parts[0] == "-" {
		negative = true
	} else if parts[0] != "+" || len(parts) < 2 {
		return 0, errors.New("Not Expected: '" + parts[0] + "', expected '+P','-P' or 'P'.")
	}
	dt := strings.SplitN(parts[1], "T", 2)

	if len(dt) == 1 {
		//dur-day or dur-week
		if strings.HasSuffix(dt[0], "D") {
			out, err = parseDay(dt[0])
			if err != nil {
				return 0, err
			}
		} else if strings.HasSuffix(dt[0], "W") {
			i, err := strconv.Atoi(dt[0][0 : len(dt[0])-1])
			if err != nil {
				return 0, err
			} else {
				out = time.Hour * time.Duration(7*24*i)
			}
		}
	} else {
		days := time.Duration(0)
		dur := time.Duration(0)
		if dt[0] != "" {
			days, err = parseDay(dt[0])
		}
		dur, err = time.ParseDuration(strings.ToLower(dt[1]))
		if err != nil {
			return 0, err
		}
		out = dur + days
	}
	if negative {
		out = -out
	}
	return
}
func parseDay(s string) (time.Duration, error) {
	if !strings.HasSuffix(s, "D") {
		return 0, errors.New("Expected D as Suffix: " + s)
	}
	i, err := strconv.Atoi(s[0 : len(s)-1])
	if err != nil {
		return 0, err
	} else {
		return time.Hour * time.Duration(24*i), nil
	}

}
