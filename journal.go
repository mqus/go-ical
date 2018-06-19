package go_ical

import (
	"strings"

	"github.com/mqus/go-ical/im"
	"github.com/soh335/icalparser"
)

type Journal struct {
	//required:
	DTStamp DateTimeVal
	Uid     StringVal

	//optional
	//Warning: Look at param VALUE, maybe only date is important (VALUE=DATE)
	DTStart DateTimeVal

	Class     *StringVal
	Created   *DateTimeVal
	LastMod   *DateTimeVal
	Organizer *PersonVal
	RecurID   *RecurIDVal
	Seq       *IntVal
	Status    *StatusVal
	Summary   *TextVal
	URL       *URIVal

	//TODO must be merged
	RRule *RecurrentRule

	Attachments []*DataVal
	Attendee    []*PersonVal
	Categories  []*Categories
	Comment     []*TextVal
	Contact     []*TextVal
	Desc        []*TextVal
	ExDate      []*DateTimeVal

	//maybe important param: RELTYPE, specified in RFC5545/3.2.15
	Related []*RelationVal
	RDate   []*RecurrenceSetVal
	RStatus []*ReqStatusVal

	OtherProperties []*icalparser.ContentLine
	OtherComponents []*im.Component

	//since RFC 7986:
	Color  *ColorVal
	Images []*DataVal
}

func (cp *CalParser) parseVJOURNAL(comp *im.Component) (out *Journal, err error, err2 error) {
	out = &Journal{}
	for _, prop := range comp.Properties {
		switch strings.ToLower(prop.Name.C) {
		case prDTStamp:
			x, err := cp.ToDateTimeVal(prop, true)
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

		case prCreated:
			x, err := cp.ToDateTimeVal(prop, true)
			if err != nil {
				//Error while parsing
				//MAYBE return err
				// instead of silently discarding Created
			} else {
				out.Created = &x
			}

		case prLastMod:
			x, err := cp.ToDateTimeVal(prop, true)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding LastMod
			} else {
				out.LastMod = &x
			}

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

		case prRid:
			//RANGE param? TODO
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

		case prDesc:
			x, err := ToTextVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding altrep
			}
			out.Desc = append(out.Desc, &x)

		case prExcDate:
			x, err := cp.ToDateTimeVal(prop, false)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding ExceptionDate
			} else {
				out.ExDate = append(out.ExDate, &x)
			}

		case prRelTo:
			x := ToRelationVal(prop)
			out.Related = append(out.Related, &x)

		case prRDate:
			x, err := cp.ToRecurrenceSetVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding ExceptionDate
			} else {
				out.RDate = append(out.RDate, &x)
			}

		case prReqStat:
			x, err := ToReqStatusVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding ReqStatus
			} else {
				out.RStatus = append(out.RStatus, &x)
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

		default:
			out.OtherProperties = append(out.OtherProperties, prop)
		}

	}
	//MAYBE don't silently discard other Components
	// out.OtherComponents = comp.Comps
	//TODO Conformance Checking
	return
}
