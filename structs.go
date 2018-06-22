package go_ical

import (
	"net/url"
	"time"

	"strings"

	"encoding/base64"

	"github.com/pkg/errors"

	"strconv"

	cl "github.com/mqus/go-contentline"
	"github.com/mqus/go-ical/util"
)

type Calendar struct {
	//Components
	Events          []*Event
	ToDos           []*ToDo
	Journals        []*Journal
	FreeBusy        []*FreeBusy
	TimeZones       map[string]*TimeZone
	OtherComponents []*cl.Component
	//Properties
	//required
	ProdID  StringVal
	Version Version
	//optional
	CalScale        *StringVal
	Method          *StringVal
	OtherProperties []*cl.Property
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

func ToTextVal(line *cl.Property) (TextVal, error) {
	var err error
	out := TextVal{}
	for name, values := range line.Parameters {
		switch name {
		case paAltRep:
			out.AltRep, err = url.Parse(strings.Trim(values[0], "\""))
		case paLang:
			out.Lang = values[0]
		default:
			out.OtherParam[name] = values
		}
	}
	return out, err

}

type StringVal struct {
	Value      string
	OtherParam cl.Parameters
}

func ToStringVal(line *cl.Property) StringVal {
	return StringVal{line.Value, line.Parameters}
}

type DateTimeVal struct {
	Value      time.Time
	Floating   bool
	OnlyDate   bool
	OtherParam cl.Parameters
}

//TODO Timezones!!  (get the Timezone specified by the TZID param, found in VTIMEZONES
func (cp *CalParser) ToDateTimeVal(line *cl.Property, isUtc bool) (DateTimeVal, error) {
	tstr := strings.Trim(line.Value, "Z")
	var err error
	out := DateTimeVal{}
	out.OnlyDate = false
	for name, values := range line.Parameters {
		switch name {
		case paValue:
			out.OnlyDate = strings.ToUpper(values[0]) == tDate
		default:
			out.OtherParam[name] = values
		}
	}

	out.Floating = !isUtc && !strings.HasSuffix(line.Value, "Z")
	if out.OnlyDate {
		out.Value, err = time.Parse(util.ISO8601_2004_DATE, tstr)
	} else {
		out.Value, err = time.Parse(util.ISO8601_2004, tstr)
	}
	return out, err
}

type URIVal struct {
	URI        *url.URL
	OtherParam cl.Parameters
}

func ToURIVal(line *cl.Property) (URIVal, error) {
	var err error
	out := URIVal{}
	out.OtherParam = line.Parameters

	out.URI, err = url.Parse(line.Value)
	return out, err

}

type DurVal struct {
	Value      time.Duration
	OtherParam cl.Parameters
}

func ToDurVal(line *cl.Property) (DurVal, error) {
	var err error
	out := DurVal{}
	out.OtherParam = line.Parameters
	out.Value, err = util.ParseDuration(line.Value)
	//if err!=nil && len(line.Parameters[paValue])==1 || line.Parameters[paValue][0] !="DURATION"{
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
func ToDataVal(line *cl.Property) (DataVal, error, error) {
	var err, err2 error
	out := DataVal{}
	out.OtherParam = line.Parameters
	var valtype, enc string
	for name, values := range line.Parameters {
		switch name {
		case paAltRep:
			out.AltRep, err2 = url.Parse(strings.Trim(values[0], "\""))
		case paValue:
			valtype = values[0]
		case paFmtType:
			out.MediaType = values[0]
		case paDisp:
			out.Display = values[0]
		case paEnc:
			enc = values[0]
		default:
			out.OtherParam[name] = values
		}
	}
	//assertions: ENCODING=BASE64
	//			  VALUE is either BINARY or URL

	if valtype == tBinary {
		out.Data, err = base64.StdEncoding.DecodeString(line.Value)
		if err != nil && enc != "" && enc != "BASE64" {
			err = errors.New("Unrecognized Encoding Format: " + line.BeforeParsing())
		}
	} else if valtype == tURI {
		out.URI, err = url.Parse(line.Value)
	} else if valtype == "" {
		err = errors.New("No Value Type Defined: " + line.BeforeParsing())
	} else {
		err = errors.New("Unrecognized Value Type: " + line.BeforeParsing())
	}
	return out, err, err2

}

//MAYBE:eventually we will have to add minver and maxver, but for now this is enough.
type Version StringVal

func ToVersion(line *cl.Property) Version {
	return Version{line.Value, line.Parameters}
}

type ColorVal StringVal

func ToColorVal(line *cl.Property) ColorVal {
	return ColorVal{line.Value, line.Parameters}
}

type Categories struct {
	Values     []string
	OtherParam cl.Parameters
}

func ToCategVal(line *cl.Property) Categories {
	return Categories{
		strings.Split(line.Value, ","),
		line.Parameters,
	}
}

type BoolVal struct {
	Value      bool
	OtherParam cl.Parameters
}

type FloatVal struct {
	Value      float64
	OtherParam cl.Parameters
}

type IntVal struct {
	Value      int32
	OtherParam cl.Parameters
}

func ToIntVal(line *cl.Property) (out IntVal, err error) {
	out = IntVal{}
	out.OtherParam = line.Parameters
	i, err := strconv.ParseInt(line.Value, 10, 32)
	out.Value = int32(i)
	return
}

type RecurrenceSetVal struct {
	From         []time.Time
	For          []time.Duration
	Floating     bool
	OnlyDate     bool
	OnlyDateTime bool

	OtherParam cl.Parameters
}

//TODO Timezones!!  (get the Timezone specified by the TZID param, found in VTIMEZONES
func (cp *CalParser) ToRecurrenceSetVal(line *cl.Property) (out RecurrenceSetVal, err error) {
	//todo 	tstr := strings.Trim(line.Value, "Z")
	out = RecurrenceSetVal{
		From:         nil,
		For:          nil,
		Floating:     !strings.Contains(line.Value, "Z"),
		OnlyDate:     false,
		OnlyDateTime: false,
		OtherParam:   nil,
	}
	vtype := tDateTime
	for name, values := range line.Parameters {
		switch name {
		case paValue:
			vtype = strings.ToUpper(values[0])
		default:
			out.OtherParam[name] = values
		}
	}
	values := strings.Split(line.Value, ",")
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
			err = errors.New("Not a valid Type: " + vtype + " in: " + line.BeforeParsing())
			return
		}
	}
	return out, err
}

type TimeVal DateTimeVal

type GeoVal struct {
	Lat, Long  float64
	OtherParam cl.Parameters
}

func ToGeoVal(line *cl.Property) (out GeoVal, err error) {
	out = GeoVal{}
	out.OtherParam = line.Parameters
	latlongs := strings.SplitN(line.Value, ";", 2)
	out.Lat, err = strconv.ParseFloat(latlongs[0], 64)
	out.Long, err = strconv.ParseFloat(latlongs[1], 64)
	if err == nil {
		if out.Lat > 90 || out.Lat < -90 {
			err = errors.New("latitude out of bounds [-90;90]:" + line.BeforeParsing())
		}
		if out.Long > 180 || out.Long < -180 {
			err = errors.New("longitude out of bounds [-180;180]:" + line.BeforeParsing())
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

	OtherParam cl.Parameters
}

func ToPersonVal(line *cl.Property) (out PersonVal, err, err2 error) {
	out = PersonVal{}
	out.Address, err = url.Parse(line.Value)
	if err != nil {
		return
	}
	//parse parameters
	for name, values := range line.Parameters {
		switch name {
		case paCN:
			out.CN = values[0]
		case paDir:
			out.Dir, err2 = url.Parse(strings.Trim(values[0], "\""))
		case paSentBy:
			out.SentBy, err2 = url.Parse(strings.Trim(values[0], "\""))
		case paLang:
			out.Lang = values[0]

		// attendees
		case paCUType:
			out.CUType = values[0]
		case paMbr:
			out.Member, err2 = url.Parse(strings.Trim(values[0], "\""))
		case paRole:
			out.Role = values[0]
		case paPartStat:
			out.PartStat = values[0]
		case paRSVP:
			if strings.ToUpper(values[0]) == "TRUE" {
				out.RSvp = true
			} else if strings.ToUpper(values[0]) == "FALSE" {
				out.RSvp = false
			} else {
				err2 = errors.New("RSVP has to be either TRUE or FALSE :" + line.BeforeParsing())
			}
		case paDelTo:
			out.DelegTo, err2 = url.Parse(strings.Trim(values[0], "\""))
		case paDelFrom:
			out.DelegFrom, err2 = url.Parse(strings.Trim(values[0], "\""))

		case paEMail:
			out.EMail = values[0]

		default:
			out.OtherParam[name] = values
		}
	}
	return
}

type ReqStatusVal struct {
	Code, Text, Data string
	OtherParam       cl.Parameters
}

func ToReqStatusVal(line *cl.Property) (out ReqStatusVal, err error) {
	out = ReqStatusVal{}
	out.OtherParam = line.Parameters
	all := strings.SplitN(line.Value, ";", 3)
	if len(all) < 2 {
		err = errors.New("REQUEST-STATUS has to have at least two ';'-separated values: " + line.BeforeParsing())
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

func ToConferenceVal(line *cl.Property) (out ConferenceVal, err, err2 error) {
	out = ConferenceVal{}
	out.URI, err = url.Parse(line.Value)
	if err != nil {
		return
	}
	valtype := ""
	//parse parameters
	for name, values := range line.Parameters {
		switch name {
		case paValue:
			valtype = strings.ToUpper(values[0])
		case paFeature:
			//out.Features = strings.Split(strings.ToLower(pval), ",")
			out.Features = values
		case paLabel:
			out.Label = strings.Trim(values[0], "\"")
		case paLang:
			out.Lang = values[0]

		default:
			out.OtherParam[name] = values
		}
	}
	if valtype != tURI {
		err2 = errors.New("VALUE must be defined as URI: " + line.BeforeParsing())
	}
	return
}

type TriggerVal struct {
	IsRelative bool
	Delay      *time.Duration
	Time       *time.Time
	RelToEnd   bool
	OtherParam cl.Parameters
}

func ToTriggerVal(line *cl.Property) (out TriggerVal, err, err2 error) {
	out = TriggerVal{
		IsRelative: false,
		Delay:      nil,
		Time:       nil, //IS ALWAYS UTC
		RelToEnd:   false,
		OtherParam: nil,
	}
	valtype := tDur
	//parse parameters
	for name, values := range line.Parameters {
		switch name {
		case paValue:
			valtype = strings.ToUpper(values[0])
		case paTrRel:
			val := strings.ToUpper(values[0])
			if val == "END" {
				out.RelToEnd = true
			} else if val != "START" {
				err2 = errors.New("RELATED - property must be START or END but was " + val + ": " + line.BeforeParsing())
			}
		default:
			out.OtherParam[name] = values
		}
	}
	if valtype == tDur {
		out.IsRelative = true
		x, err1 := util.ParseDuration(line.Value)
		out.Delay = &x
		err = err1
	} else if valtype == tDateTime {
		x, err1 := time.Parse(util.ISO8601_2004, strings.Trim(line.Value, "Z"))
		out.Time = &x
		err = err1
	} else {
		err = errors.New("Unexpected Value type '" + valtype + "': " + line.BeforeParsing())
	}
	return
}

type FBVal struct {
	From   []time.Time
	For    []time.Duration
	FBType string

	OtherParam cl.Parameters
}

//TODO Timezones!!  (get the Timezone specified by the TZID param, found in VTIMEZONES
func (cp *CalParser) ToFBVal(line *cl.Property) (out FBVal, err error) {
	//todo 	tstr := strings.Trim(line.Value, "Z")
	out = FBVal{}
	for name, values := range line.Parameters {
		switch name {
		case paFBType:
			out.FBType = values[0]
		default:
			out.OtherParam[name] = values
		}
	}
	for _, val := range strings.Split(line.Value, ",") {
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

func ToRelationVal(line *cl.Property) (out RelationVal) {
	for name, values := range line.Parameters {
		switch name {
		case paRelType:
			out.RelType = values[0]
		default:
			out.OtherParam[name] = values
		}
	}
	out.Value = line.Value
	return
}
