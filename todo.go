package go_ical

import (
	"strings"

	"github.com/mqus/go-ical/im"
)

type ToDo struct {
	Alarms []*Alarm
	//required:
	DTStamp DateTimeVal
	Uid     StringVal

	//required if Calendar.Method is nil or if RRULE is specified
	//optional
	Class     *StringVal
	Completed *DateTimeVal
	Created   *DateTimeVal
	Desc      *TextVal
	//Warning: Look at param VALUE, maybe only date is important (VALUE=DATE)
	DTStart   *DateTimeVal
	Geo       *GeoVal
	LastMod   *DateTimeVal
	Location  *TextVal
	Organizer *PersonVal
	Percent   *IntVal
	Priority  *IntVal
	//more complexity hidden here
	RecurID *RecurIDVal
	Seq     *IntVal
	Status  *StatusVal
	Summary *TextVal
	URL     *URIVal

	//TODO must be merged
	RRule *RecurrentRule

	//use only either one of DTDue or DTDuration

	//Warning: Look at param VALUE, maybe only date is important (VALUE=DATE)
	DTDue      *DateTimeVal
	DTDuration *DurVal

	//optional, multiple
	Attachments []*DataVal
	Attendee    []*PersonVal
	Categories  []*Categories
	Comment     []*TextVal
	Contact     []*TextVal
	//Warning: Look at param VALUE, maybe only date is important (VALUE=DATE)
	ExDate []*DateTimeVal

	RStatus []*ReqStatusVal

	//maybe important param: RELTYPE, specified in RFC5545/3.2.15
	Related []*RelationVal
	//in this case a single TextVal can also be multiple comma-separated resources.
	Resources       []*TextVal
	RDate           []*RecurrenceSetVal
	OtherProperties []*im.Property
	OtherComponents []*im.Component

	//since RFC 7986:
	Color      *ColorVal
	Images     []*DataVal
	Conference []*ConferenceVal
}

func (cp *CalParser) parseVTODO(comp *im.Component) (out *ToDo, err error, err2 error) {
	out = &ToDo{}
	for _, prop := range comp.Properties {
		switch prop.Name {
		case prDTStamp:
			x, err := cp.ToDateTimeVal(prop, false)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding DTStamp
			} else {
				out.DTStamp = x
			}

		case prUid:
			out.Uid = ToStringVal(prop)

		case prClass:
			x := ToStringVal(prop)
			out.Class = &x

		case prCompl:
			x, err := cp.ToDateTimeVal(prop, true)
			if err != nil {
				//Error while parsing
				//MAYBE return err
				// instead of silently discarding Completed
			} else {
				out.Completed = &x
			}

		case prCreated:
			x, err := cp.ToDateTimeVal(prop, true)
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

		case prDTStart:
			x, err := cp.ToDateTimeVal(prop, false)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding DTStart
			} else {
				out.DTStart = &x
			}

		case prGeo:
			x, err := ToGeoVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding Geo
			} else {
				out.Geo = &x
			}

		case prLastMod:
			x, err := cp.ToDateTimeVal(prop, true)
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

		case prPctCompl:
			x, err := ToIntVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding Priority
			} else {
				out.Priority = &x
			}

		case prPrio:
			x, err := ToIntVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding Priority
			} else {
				out.Priority = &x
			}

		case prRid:
			//TODO RANGE param
			x, err := cp.ToDateTimeVal(prop, false)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding RecurrenceID
			} else {
				y := RecurIDVal(x)
				out.RecurID = &y
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

		case prURL:
			x, err := ToURIVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding URL
			} else {
				out.URL = &x
			}

		case prRRule:
			x := RecurrentRule(ToStringVal(prop))
			out.RRule = &x

		case prDue:
			x, err := cp.ToDateTimeVal(prop, false)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding DTDue
			} else {
				out.DTDue = &x
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
			x, err := cp.ToDateTimeVal(prop, false)
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
			x := ToRelationVal(prop)
			out.Related = append(out.Related, &x)

		case prRes:
			x, err := ToTextVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding altrep
			}
			out.Resources = append(out.Resources, &x)

		case prRDate:
			x, err := cp.ToRecurrenceSetVal(prop)
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
			al, e1, e2 := cp.parseVALARM(subcomp)
			if e1 != nil {
				//MAYBE return err
				// instead of silently discarding Alarm component
			} else {
				if e2 != nil {
					//MAYBE return err2
					// instead of silently ignoring property
				}
				out.Alarms = append(out.Alarms, al)
			}
		default:
			//MAYBE don't silently discard other Components
			// out.OtherComponents = append(out.OtherComponents, subcomp)
		}
	}
	//TODO Conformance Checking
	return
}
