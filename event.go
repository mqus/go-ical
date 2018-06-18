package go_ical

import (
	"strings"

	"github.com/mqus/go-ical/im"
	"github.com/soh335/icalparser"
)

type Event struct {
	Alarms []*Alarm
	//required:
	DTStamp DateTimeVal
	Uid     StringVal

	//required if Calendar.Method is nil or if RRULE is specified
	//Warning: Look at param VALUE, maybe only date is important (VALUE=DATE)
	DTStart DateTimeVal
	//optional
	Class     *StringVal
	Created   *DateTimeVal
	Desc      *TextVal
	Geo       *GeoVal
	LastMod   *DateTimeVal
	Location  *TextVal
	Organizer *PersonVal
	Priority  *IntVal
	Seq       *IntVal
	Status    *StatusVal
	Summary   *TextVal
	Transp    *TTranspVal
	URL       *URIVal
	//more complexity hidden here
	RecurID *RecurIDVal

	RRule *RecurrentRule

	//use only either one of DTEnd or DTDuration

	//Warning: Look at param VALUE, maybe only date is important (VALUE=DATE)
	DTEnd      *DateTimeVal
	DTDuration *DurVal

	//optional, multiple
	Attachments []*DataVal
	//Look at param for the type of attendee
	Attendee   []*PersonVal
	Categories []*Categories
	Comment    []*TextVal
	Contact    []*TextVal
	//Warning: Look at param VALUE, maybe only date is important (VALUE=DATE)
	ExDate []*DateTimeVal

	RStatus []*ReqStatusVal

	//maybe important param: RELTYPE, specified in RFC5545/3.2.15
	Related []*RelationVal
	//in this case a single TextVal can also be multiple comma-separated resources.
	Resources []*TextVal
	//Warning: Look at param VALUE, maybe only date/datetime is important (VALUE=DATE)
	RDate           []*RecurrenceSetVal
	OtherProperties []*icalparser.ContentLine
	OtherComponents []*im.Component
	//since RFC 7986:
	Color      *ColorVal
	Images     []*DataVal
	Conference []*ConferenceVal
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
			out.Uid = ToStringVal(prop)

		case prDTStart:
			x, err := ToDateTimeVal(prop, false)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding DTStart
			} else {
				out.DTStart = x
			}

		case prClass:
			x := ToStringVal(prop)
			out.Class = &x

		case prCreated:
			x, err := ToDateTimeVal(prop, false)
			if err != nil {
				//Error while parsing
				//MAYBE return err
				// instead of silently discarding Created
			} else {
				out.Created = &x
			}

		case prDesc:
			x, err := ToTextVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding altrep
			}
			out.Desc = &x
			//out.Desc = append(out.Desc, x)

		case prGeo:
			x, err := ToGeoVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding Geo
			} else {
				out.Geo = &x
			}

		case prLastMod:
			x, err := ToDateTimeVal(prop, true)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding LastMod
			} else {
				out.LastMod = &x
			}

		case prLoc:
			x, err := ToTextVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding altrep
			}
			out.Location = &x

		case prOrganizer:
			x, err, err2 := ToPersonVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding Organizer
			} else {
				if err2 != nil {
					//MAYBE return err2
					// instead of silently discarding parameters
				}
				out.Organizer = &x
			}

		case prPrio:
			x, err := ToIntVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding Priority
			} else {
				out.Priority = &x
			}

		case prSeq:
			x, err := ToIntVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding Seq
			} else {
				out.Seq = &x
			}

		case prStat:
			x := StatusVal(ToStringVal(prop))
			out.Status = &x

		case prSumm:
			x, err := ToTextVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding altrep
			}
			out.Summary = &x

		case prTTransp:
			x := TTranspVal(ToStringVal(prop))
			out.Transp = &x

		case prURL:
			x, err := ToURIVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding URL
			} else {
				out.URL = &x
			}

		case prRid:
			x, err := ToDateTimeVal(prop, false)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding RecurrenceID
			} else {
				y := RecurIDVal(x)
				out.RecurID = &y
			}

		case prRRule:
			x := RecurrentRule(ToStringVal(prop))
			out.RRule = &x

		case prDTEnd:
			x, err := ToDateTimeVal(prop, false)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding DTEnd
			} else {
				out.DTEnd = &x
			}

		case prDur:
			x, err := ToDurVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding Duration
			} else {
				out.DTDuration = &x
			}

		case prAttach:
			x, err, err2 := ToDataVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding Attachment
			} else {
				if err2 != nil {
					//MAYBE return err
					// instead of silently discarding altrep
				}
				out.Attachments = append(out.Attachments, &x)
			}

		case prAttendee:
			x, err, err2 := ToPersonVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding Attendee
			} else {
				if err2 != nil {
					//MAYBE return err2
					// instead of silently discarding parameters
				}
				out.Attendee = append(out.Attendee, &x)
			}

		case prCat:
			x := ToCategVal(prop)
			out.Categories = append(out.Categories, &x)

		case prComm:
			x, err := ToTextVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding altrep
			}
			out.Comment = append(out.Comment, &x)

		case prContact:
			x, err := ToTextVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding altrep
			}
			out.Contact = append(out.Contact, &x)

		case prExcDate:
			x, err := ToDateTimeVal(prop, false)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding ExceptionDate
			} else {
				out.ExDate = append(out.ExDate, &x)
			}

		case prReqStat:
			x, err := ToReqStatusVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding ReqStatus
			} else {
				out.RStatus = append(out.RStatus, &x)
			}
		case prRelTo:
			x := RelationVal(ToStringVal(prop))
			out.Related = append(out.Related, &x)

		case prRes:
			x, err := ToTextVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding altrep
			}
			out.Resources = append(out.Resources, &x)

		case prRDate:
			x, err := ToRecurrenceSetVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding ExceptionDate
			} else {
				out.RDate = append(out.RDate, &x)
			}

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
			x, err, err2 := ToConferenceVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding ConferenceVal
			} else {
				if err2 != nil {
					//MAYBE return err2
					// instead of silently ignoring value type
				}
				out.Conference = append(out.Conference, &x)
			}

		default:
			out.OtherProperties = append(out.OtherProperties, prop)
		}

	}

	for _, subcomp := range comp.Comps {
		switch strings.ToLower(subcomp.Name) {
		case vAl:
			al, e1, e2 := parseVALARM(subcomp)
			if e1 != nil {
				//MAYBE return err
				// instead of silently discarding Alarm component
			} else {
				if e2 != nil {
					//MAYBE return err2
					// instead of silently ignoring property
				}
				out.Alarms = append(out.Alarms, &al)
			}
		default:
			out.OtherComponents = append(out.OtherComponents, subcomp)
		}
	}
	//TODO Conformance Checking
	return
}
