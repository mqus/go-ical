package go_ical

import (
	"strings"

	"github.com/mqus/go-ical/im"
	"github.com/soh335/icalparser"
)

type Event struct {
	Alarms []Alarm
	//required:
	DTStamp DateTimeVal
	Uid     StringVal

	//required if Calendar.Method is nil or if RRULE is specified
	//Warning: Look at param VALUE, maybe only date is important (VALUE=DATE)
	DTStart DateTimeVal
	//optional
	Class     *TextVal
	Created   *DateTimeVal
	Desc      *TextVal
	Geo       *GeoVal
	LastMod   *DateTimeVal
	Location  *TextVal
	Organizer *PersonVal
	Priority  *IntVal
	Seq       *IntVal
	Status    *TextVal
	Summary   *TextVal
	Transp    *BoolVal
	URL       *URIVal
	//more complexity hidden here
	RecurID TextVal

	RRule RecurrentRule

	//use only either one of DTEnd or DTDuration

	//Warning: Look at param VALUE, maybe only date is important (VALUE=DATE)
	DTEnd      DateTimeVal
	DTDuration DurVal

	//optional, multiple
	Attachments []DataVal
	//Look at param for the type of attendee
	Attendee   []PersonVal
	Categories []Categories
	Comment    []TextVal
	Contact    []TextVal
	//Warning: Look at param VALUE, maybe only date is important (VALUE=DATE)
	ExDate  []DateTimeVal
	RStatus []StatusVal
	//maybe important param: RELTYPE, specified in RFC5545/3.2.15
	Related []TextVal
	//in this case a single TextVal can also be multiple comma-separated resources.
	Resources []TextVal
	//Warning: Look at param VALUE, maybe only date/datetime is important (VALUE=DATE)
	RDate           []PeriodVal
	OtherProperties []*icalparser.ContentLine
	OtherComponents []*im.Component
	//since RFC 7986:
	Color      *ColorVal
	Images     []*DataVal
	Conference []*URIVal
}

func parseVEVENT(comp *im.Component) (out *Event, err error, err2 error) {
	out = &Event{}
	for _, prop := range comp.Properties {
		switch strings.ToLower(prop.Name.C) {
		case prDTStamp:
			x, err := ToDateTimeVal(prop, false)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding DTStamp
			} else {
				out.DTStamp = x
			}

		case prUid:
			x := ToStringVal(prop)
			out.Uid = x

		case prDTStart:
			x, err := ToDateTimeVal(prop, false)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding DTStart
			} else {
				out.DTStamp = x
			}

		case prClass:
			fallthrough //TODO

		case prCreated:
			fallthrough //TODO

		case prDesc:
			x, err := ToTextVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding altrep
			}
			out.Desc = &x
			//out.Desc = append(out.Desc, x)

		case prGeo:
			fallthrough //TODO

		case prLastMod:
			x, err := ToDateTimeVal(prop, true)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding LastMod
			} else {
				out.LastMod = &x
			}

		case prLoc:
			fallthrough //TODO

		case prOrganizer:
			fallthrough //TODO

		case prPrio:
			fallthrough //TODO

		case prSeq:
			fallthrough //TODO

		case prStat:
			fallthrough //TODO

		case prSumm:
			fallthrough //TODO

		case prTransp:
			fallthrough //TODO

		case prURL:
			x, err := ToURIVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding URL
			} else {
				out.URL = &x
			}

		case prRid:
			fallthrough //TODO

		case prRRule:
			fallthrough //TODO

		case prDTEnd:
			fallthrough //TODO

		case prDur:
			fallthrough //TODO

		case prAttach:
			fallthrough //TODO

		case prAttendee:
			fallthrough //TODO

		case prCat:
			out.Categories = append(out.Categories, ToCategVal(prop))

		case prComm:
			fallthrough //TODO

		case prContact:
			fallthrough //TODO

		case prExcDate:
			fallthrough //TODO

		case prReqStat:
			fallthrough //TODO

		case prRelTo:
			fallthrough //TODO

		case prRes:
			fallthrough //TODO

		case prRDate:
			fallthrough //TODO

		case prColor:
			x := ToColorVal(prop)
			out.Color = &x

		case prImage:
			x, err, err2 := ToDataVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding Image
			} else {
				if err2 != nil {
					//MAYBE return err
					// instead of silently discarding altrep
				}
				out.Images = append(out.Images, &x)
			}

		case prConf:
			fallthrough //TODO
		default:
			out.OtherProperties = append(out.OtherProperties, prop)
		}

	}

	return
}

type Alarm struct {
	//TODO
	OtherProperties []icalparser.ContentLine
}
