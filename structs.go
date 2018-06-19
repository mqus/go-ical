package go_ical

import (
	"net/url"
	"time"

	"strings"

	"encoding/base64"

	"github.com/pkg/errors"

	"strconv"

	"github.com/mqus/go-ical/im"
	"github.com/mqus/go-ical/util"
	"github.com/soh335/icalparser"
)

type Calendar struct {
	//Components
	Events          []*Event
	ToDos           []*ToDo
	Journals        []*Journal
	FreeBusy        []*FreeBusy
	TimeZones       map[string]*TimeZone
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
	StringVal
	AltRep *url.URL

	//must conform with RFC5646
	Lang string
}

func ToTextVal(line *icalparser.ContentLine) (TextVal, error) {
	var err error
	out := TextVal{}
	for _, param := range line.Param {
		pval := param.ParamValues[0].C
		switch strings.ToLower(param.ParamName.C) {
		case paAltRep:
			out.AltRep, err = url.Parse(strings.Trim(pval, "\""))
		case paLang:
			out.Lang = pval
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
func (cp *CalParser) ToDateTimeVal(line *icalparser.ContentLine, isUtc bool) (DateTimeVal, error) {
	tstr := strings.Trim(line.Value.C, "Z")
	var err error
	out := DateTimeVal{}
	out.OnlyDate = false
	for _, param := range line.Param {
		pval := param.ParamValues[0].C
		switch strings.ToLower(param.ParamName.C) {
		case paValue:
			out.OnlyDate = strings.ToLower(pval) == tDate
		default:
			out.OtherParam = append(out.OtherParam, param)
		}
	}

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
	out.Value, err = util.ParseDuration(line.Value.C)
	//idx := util.GetParam(line.Param, "VALUE")
	//if err!=nil && idx<0 || line.Param[idx].ParamValues[0].C !="DURATION"{
	//	//MAYBE VALUE=DURATION is theoretically required
	//}
	return out, err
}

type DataVal struct {
	//binary data must always be encoded with base64!
	Data []byte
	URIVal
	AltRep    *url.URL
	MediaType string
	Display   string
}

//TO?DO DisplayParam (new RFC)
func ToDataVal(line *icalparser.ContentLine) (DataVal, error, error) {
	var err, err2 error
	out := DataVal{}
	out.OtherParam = line.Param
	var valtype, enc string
	for _, param := range line.Param {
		pval := param.ParamValues[0].C
		switch strings.ToLower(param.ParamName.C) {
		case paAltRep:
			out.AltRep, err2 = url.Parse(strings.Trim(pval, "\""))
		case paValue:
			valtype = pval
		case paFmtType:
			out.MediaType = pval
		case paDisp:
			out.Display = pval
		case paEnc:
			enc = pval
		default:
			out.OtherParam = append(out.OtherParam, param)
		}
	}
	//assertions: ENCODING=BASE64
	//			  VALUE is either BINARY or URL

	if valtype == tBinary {
		out.Data, err = base64.StdEncoding.DecodeString(line.Value.C)
		if err != nil && enc != "" && enc != "BASE64" {
			err = errors.New("Unrecognized Encoding Format: " + line.String())
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

func ToIntVal(line *icalparser.ContentLine) (out IntVal, err error) {
	out = IntVal{}
	out.OtherParam = line.Param
	i, err := strconv.ParseInt(line.Value.C, 10, 32)
	out.Value = int32(i)
	return
}

type RecurrenceSetVal struct {
	From         []time.Time
	For          []time.Duration
	Floating     bool
	OnlyDate     bool
	OnlyDateTime bool

	OtherParam []*icalparser.Param
}

//TODO Timezones!!  (get the Timezone specified by the TZID param, found in VTIMEZONES
func (cp *CalParser) ToRecurrenceSetVal(line *icalparser.ContentLine) (out RecurrenceSetVal, err error) {
	//todo 	tstr := strings.Trim(line.Value.C, "Z")
	out = RecurrenceSetVal{
		From:         nil,
		For:          nil,
		Floating:     !strings.Contains(line.Value.C, "Z"),
		OnlyDate:     false,
		OnlyDateTime: false,
		OtherParam:   nil,
	}
	vtype := tDateTime
	for _, param := range line.Param {
		pval := param.ParamValues[0].C
		switch strings.ToLower(param.ParamName.C) {
		case paValue:
			vtype = strings.ToLower(pval)
		default:
			out.OtherParam = append(out.OtherParam, param)
		}
	}
	values := strings.Split(line.Value.C, ",")
	for _, val := range values {
		switch vtype {
		case tDate:
			out.OnlyDate = true
			date, err := time.Parse(util.ISO8601_2004_DATE, val)
			if err != nil {
				return out, err
			}
			out.From = append(out.From, date)

		case tDateTime:
			out.OnlyDateTime = true
			date, err := time.Parse(util.ISO8601_2004, strings.Trim(val, "Z"))
			if err != nil {
				return out, err
			}
			out.From = append(out.From, date)

		case tPeriod:
			x1, x2, er := util.ParsePeriod(val)
			if er != nil {
				return out, er
			}
			out.From = append(out.From, x1)
			out.For = append(out.For, x2)
		default:
			err = errors.New("Not a valid Type: " + vtype + " in: " + line.String())
			return
		}
	}
	return out, err
}

type TimeVal DateTimeVal

type GeoVal struct {
	Lat, Long  float64
	OtherParam []*icalparser.Param
}

func ToGeoVal(line *icalparser.ContentLine) (out GeoVal, err error) {
	out = GeoVal{}
	out.OtherParam = line.Param
	latlongs := strings.SplitN(line.Value.C, ";", 2)
	out.Lat, err = strconv.ParseFloat(latlongs[0], 64)
	out.Long, err = strconv.ParseFloat(latlongs[1], 64)
	if err == nil {
		if out.Lat > 90 || out.Lat < -90 {
			err = errors.New("latitude out of bounds [-90;90]:" + line.String())
		}
		if out.Long > 180 || out.Long < -180 {
			err = errors.New("longitude out of bounds [-180;180]:" + line.String())
		}
	}
	return
}

type PersonVal struct {
	Address *url.URL
	CN      string
	Dir     *url.URL
	Lang    string
	SentBy  *url.URL
	//Attendees:
	CUType    string
	Member    *url.URL
	Role      string
	PartStat  string
	RSvp      bool
	DelegTo   *url.URL
	DelegFrom *url.URL
	//from RFC7986 (section 6.2), when Address can not be used as an Email-Adress.
	EMail string

	OtherParam []*icalparser.Param
}

func ToPersonVal(line *icalparser.ContentLine) (out PersonVal, err, err2 error) {
	out = PersonVal{}
	out.Address, err = url.Parse(line.Value.C)
	if err != nil {
		return
	}
	//parse parameters
	for _, param := range line.Param {
		pval := param.ParamValues[0].C
		switch strings.ToLower(param.ParamName.C) {
		case paCN:
			out.CN = pval
		case paDir:
			out.Dir, err2 = url.Parse(strings.Trim(pval, "\""))
		case paSentBy:
			out.SentBy, err2 = url.Parse(strings.Trim(pval, "\""))
		case paLang:
			out.Lang = pval

		// attendees
		case paCUType:
			out.CUType = pval
		case paMbr:
			out.Member, err2 = url.Parse(strings.Trim(pval, "\""))
		case paRole:
			out.Role = pval
		case paPartStat:
			out.PartStat = pval
		case paRSVP:
			if strings.ToLower(pval) == "true" {
				out.RSvp = true
			} else if strings.ToLower(pval) == "false" {
				out.RSvp = false
			} else {
				err2 = errors.New("RSVP has to be either TRUE or FALSE :" + line.String())
			}
		case paDelTo:
			out.DelegTo, err2 = url.Parse(strings.Trim(pval, "\""))
		case paDelFrom:
			out.DelegFrom, err2 = url.Parse(strings.Trim(pval, "\""))

		case paEMail:
			out.EMail = pval

		default:
			out.OtherParam = append(out.OtherParam, param)
		}
	}
	return
}

type ReqStatusVal struct {
	Code, Text, Data string
	OtherParam       []*icalparser.Param
}

func ToReqStatusVal(line *icalparser.ContentLine) (out ReqStatusVal, err error) {
	out = ReqStatusVal{}
	out.OtherParam = line.Param
	all := strings.SplitN(line.Value.C, ";", 3)
	if len(all) < 2 {
		err = errors.New("REQUEST-STATUS has to have at least two ';'-separated values: " + line.String())
	} else {
		out.Code = all[0]
		out.Text = all[1]
		if len(all) == 3 {
			out.Data = all[2]
		}
	}
	return
}

type ConferenceVal struct {
	URIVal
	Features []string
	Label    string
	Lang     string
}

func ToConferenceVal(line *icalparser.ContentLine) (out ConferenceVal, err, err2 error) {
	out = ConferenceVal{}
	out.URI, err = url.Parse(line.Value.C)
	if err != nil {
		return
	}
	valtype := ""
	//parse parameters
	for _, param := range line.Param {
		pval := param.ParamValues[0].C
		switch strings.ToLower(param.ParamName.C) {
		case paValue:
			valtype = strings.ToLower(pval)
		case paFeature:
			//out.Features = strings.Split(strings.ToLower(pval), ",")
			for _, v := range param.ParamValues {
				out.Features = append(out.Features, v.C)
			}
		case paLabel:
			out.Label = strings.Trim(pval, "\"")
		case paLang:
			out.Lang = pval

		default:
			out.OtherParam = append(out.OtherParam, param)
		}
	}
	if valtype != "URI" {
		err2 = errors.New("VALUE must be defined as URI: " + line.String())
	}
	return
}

type TriggerVal struct {
	IsRelative bool
	Delay      *time.Duration
	Time       *time.Time
	RelToEnd   bool
	OtherParam []*icalparser.Param
}

func ToTriggerVal(line *icalparser.ContentLine) (out TriggerVal, err, err2 error) {
	out = TriggerVal{
		IsRelative: false,
		Delay:      nil,
		Time:       nil, //IS ALWAYS UTC
		RelToEnd:   false,
		OtherParam: nil,
	}
	valtype := tDur
	//parse parameters
	for _, param := range line.Param {
		pval := param.ParamValues[0].C
		switch strings.ToLower(param.ParamName.C) {
		case paValue:
			valtype = strings.ToLower(pval)
		case paTrRel:
			val := strings.ToLower(pval)
			if val == "end" {
				out.RelToEnd = true
			} else if val != "start" {
				err2 = errors.New("RELATED - property must be START or END but was " + val + ": " + line.String())
			}
		default:
			out.OtherParam = append(out.OtherParam, param)
		}
	}
	if valtype == tDur {
		out.IsRelative = true
		x, err1 := util.ParseDuration(line.Value.C)
		out.Delay = &x
		err = err1
	} else if valtype == tDateTime {
		x, err1 := time.Parse(util.ISO8601_2004, strings.Trim(line.Value.C, "Z"))
		out.Time = &x
		err = err1
	} else {
		err = errors.New("Unexpected Value type '" + valtype + "': " + line.String())
	}
	return
}

type FBVal struct {
	From   []time.Time
	For    []time.Duration
	FBType string

	OtherParam []*icalparser.Param
}

//TODO Timezones!!  (get the Timezone specified by the TZID param, found in VTIMEZONES
func (cp *CalParser) ToFBVal(line *icalparser.ContentLine) (out FBVal, err error) {
	//todo 	tstr := strings.Trim(line.Value.C, "Z")
	out = FBVal{}
	for _, param := range line.Param {
		switch strings.ToLower(param.ParamName.C) {
		case paFBType:
			out.FBType = param.ParamValues[0].C
		default:
			out.OtherParam = append(out.OtherParam, param)
		}
	}
	for _, val := range strings.Split(line.Value.C, ",") {
		x1, x2, er := util.ParsePeriod(val)
		if er != nil {
			return out, er
		}
		out.From = append(out.From, x1)
		out.For = append(out.For, x2)
	}
	return out, err
}

type StatusVal StringVal
type TTranspVal StringVal
type RecurIDVal DateTimeVal

type RelationVal struct {
	StringVal
	RelType string
}

func ToRelationVal(line *icalparser.ContentLine) (out RelationVal) {
	for _, par := range line.Param {
		switch strings.ToLower(par.ParamName.C) {
		case paRelType:
			out.RelType = par.ParamValues[0].C
		default:
			out.OtherParam = append(out.OtherParam, par)
		}
	}
	out.Value = line.Value.C
	return
}
