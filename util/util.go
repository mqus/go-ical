package util

import (
	"time"

	"github.com/pkg/errors"

	"strconv"
	"strings"
)

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

func ParsePeriod(s string) (from time.Time, dur time.Duration, err error) {
	vals := strings.SplitN(s, "/", 2)
	from, err = time.Parse(ISO8601_2004, strings.Trim(vals[0], "Z"))
	if err != nil {
		return time.Time{}, 0, err
	}
	if strings.Contains(vals[1], "P") {
		dur, err = ParseDuration(vals[1])
	} else {
		date2, err := time.Parse(ISO8601_2004, strings.Trim(vals[1], "Z"))
		if err != nil {
			return time.Time{}, 0, err
		}
		dur = date2.Sub(from)
	}

	return
}

func ParseUTCOffset(s string) (out time.Duration, err error) {
	in := s + "00"
	if len(s) < 5 {
		return 0, errors.New("UTC-Offset must have at least 5 symbols: " + s)
	}
	negpos := time.Duration(1)
	var tmpdur int
	if in[0:1] == "-" {
		negpos = time.Duration(-1)
	} else if in[0:1] == "+" {
		return 0, errors.New("Expected '-' or '+' at start of UTC-Offset-Value: " + s)
	}

	//parse hours
	tmpdur, err = strconv.Atoi(in[1:3])
	out = out + time.Duration(tmpdur)*time.Hour
	//parse minutes
	tmpdur, err = strconv.Atoi(in[3:5])
	out = out + time.Duration(tmpdur)*time.Minute
	//parse seconds
	tmpdur, err = strconv.Atoi(in[5:7])
	out = out + time.Duration(tmpdur)*time.Second

	//include previously parsed sign
	out = negpos * out
	return

}
