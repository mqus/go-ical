package go_ical

import (
	"net/url"
	"time"

	"strings"

	"encoding/base64"

	"errors"

	"github.com/mqus/go-ical/im"
	"github.com/mqus/go-ical/util"
	"github.com/soh335/icalparser"
)

type Calendar struct {
	//Components
	Events []*Event
	//TODO
	//	ToDos           []ToDo
	//	Journals        []Journal
	//	FreeBusyTimes   []FreeBusy
	//	TimeZones       []TimeZone
	OtherComponents []*im.Component
	//Properties
	//required
	ProdID  StringVal
	Version Version
	//optional
	CalScale        *StringVal
	Method          *StringVal
	OtherProperties []*icalparser.ContentLine
	//since RFC 7986:
	Uid        *StringVal
	LastMod    *DateTimeVal
	URL        *URIVal
	Refresh    *DurVal
	Source     *URIVal
	Color      *ColorVal
	Name       []TextVal
	Desc       []TextVal
	Categories []Categories
	Images     []DataVal
}

//Property Value Types

type TextVal struct {
	Value  string
	AltRep *url.URL

	//must conform with RFC5646
	Lang       string
	OtherParam []*icalparser.Param
}

func ToTextVal(line *icalparser.ContentLine) (TextVal, error) {
	var err error
	out := TextVal{}
	for _, param := range line.Param {
		switch strings.ToLower(param.ParamName.C) {
		case paAltRep:
			out.AltRep, err = url.Parse(strings.Trim(param.ParamValues[0].C, "\""))
		case paLang:
			out.Lang = param.ParamValues[0].C
		default:
			out.OtherParam = append(out.OtherParam, param)
		}
	}
	return out, err

}

type StringVal struct {
	Value      string
	OtherParam []*icalparser.Param
}

func ToStringVal(line *icalparser.ContentLine) StringVal {
	return StringVal{line.Value.C, line.Param}
}

type DateTimeVal struct {
	Value      time.Time
	Floating   bool
	OnlyDate   bool
	OtherParam []*icalparser.Param
}

//TODO Timezones!!  (get the Timezone specified by the TZID param, found in VTIMEZONES
func ToDateTimeVal(line *icalparser.ContentLine, isUtc bool) (DateTimeVal, error) {
	tstr := strings.Trim(line.Value.C, "Z")
	var err error
	out := DateTimeVal{}
	out.OtherParam = line.Param
	idx := util.GetParam(line.Param, paValue)
	out.OnlyDate = idx >= 0 && line.Param[idx].ParamValues[0].C == tDate
	out.Floating = !isUtc && !strings.HasSuffix(line.Value.C, "Z")
	if out.OnlyDate {
		out.Value, err = time.Parse(util.ISO8601_2004_DATE, tstr)
	} else {
		out.Value, err = time.Parse(util.ISO8601_2004, tstr)
	}
	return out, err

}

type URIVal struct {
	URI        *url.URL
	OtherParam []*icalparser.Param
}

func ToURIVal(line *icalparser.ContentLine) (URIVal, error) {
	var err error
	out := URIVal{}
	out.OtherParam = line.Param

	out.URI, err = url.Parse(line.Value.C)
	return out, err

}

type DurVal struct {
	Value      time.Duration
	OtherParam []*icalparser.Param
}

func ToDurVal(line *icalparser.ContentLine) (DurVal, error) {
	var err error
	out := DurVal{}
	out.OtherParam = line.Param
	out.Value, err = time.ParseDuration(line.Value.C)

	//idx := util.GetParam(line.Param, "VALUE")
	//if err!=nil && idx<0 || line.Param[idx].ParamValues[0].C !="DURATION"{
	//	//MAYBE VALUE=DURATION is theoretically required
	//}
	return out, err
}

type DataVal struct {
	//binary data must always be encoded with base64!
	Data   []byte
	URI    *url.URL
	AltRep *url.URL

	OtherParam []*icalparser.Param
}

//TODO DisplayParam (new RFC)
func ToDataVal(line *icalparser.ContentLine) (DataVal, error, error) {
	var err, err2 error
	out := DataVal{}
	out.OtherParam = line.Param
	var valtype, enc string
	for _, param := range line.Param {
		switch strings.ToLower(param.ParamName.C) {
		case paAltRep:
			out.AltRep, err2 = url.Parse(strings.Trim(param.ParamValues[0].C, "\""))
		case paValue:
			valtype = param.ParamValues[0].C
		case paEnc:
			enc = param.ParamValues[0].C
		default:
			out.OtherParam = append(out.OtherParam, param)
		}
	}
	//assertions: ENCODING=BASE64
	//			  VALUE is either BINARY or URL

	if valtype == tBinary {
		out.Data, err = base64.StdEncoding.DecodeString(line.Value.C)
		if err != nil && enc != "" && enc != "BASE64" {
			err = errors.New("Unrecognized Encoding Format:" + line.String())
		}
	} else if valtype == tURI {
		out.URI, err = url.Parse(line.Value.C)
	} else if valtype == "" {
		err = errors.New("No Value Type Defined: " + line.String())
	} else {
		err = errors.New("Unrecognized Value Type: " + line.String())
	}
	return out, err, err2

}

//MAYBE:eventually we will have to add minver and maxver, but for now this is enough.
type Version StringVal

func ToVersion(line *icalparser.ContentLine) Version {
	return Version{line.Value.C, line.Param}
}

type ColorVal StringVal

func ToColorVal(line *icalparser.ContentLine) ColorVal {
	return ColorVal{line.Value.C, line.Param}
}

type Categories struct {
	Values     []string
	OtherParam []*icalparser.Param
}

func ToCategVal(line *icalparser.ContentLine) Categories {
	return Categories{
		strings.Split(line.Value.C, ","),
		line.Param,
	}
}

type BoolVal struct {
	Value      bool
	OtherParam []*icalparser.Param
}

type FloatVal struct {
	Value      float64
	OtherParam []*icalparser.Param
}

type IntVal struct {
	Value      int32
	OtherParam []*icalparser.Param
}

type PeriodVal struct {
	From       time.Time
	For        time.Duration
	OtherParam []*icalparser.Param
}

type TimeVal DateTimeVal

type GeoVal struct {
	Lat, Long  float64
	OtherParam []*icalparser.Param
}

type PersonVal URIVal

type StatusVal struct {
	Code, Text string
	OtherParam []*icalparser.Param
}
