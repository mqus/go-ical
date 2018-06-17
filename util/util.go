package util

import "github.com/soh335/icalparser"

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
