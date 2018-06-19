package go_ical

import (
	"github.com/mqus/go-ical/im"
)

type FreeBusy struct {
	//required:
	DTStamp DateTimeVal
	Uid     StringVal

	//optional
	//Warning: Look at param VALUE, maybe only date is important (VALUE=DATE)
	Contact   *TextVal
	DTStart   *DateTimeVal
	DTEnd     *DateTimeVal
	Organizer *PersonVal
	URL       *URIVal

	Attendee   []*PersonVal
	Categories []*Categories
	Comment    []*TextVal
	FreeBusy   []*FBVal
	RStatus    []*ReqStatusVal

	OtherProperties []*im.Property
	OtherComponents []*im.Component
}

func (cp *CalParser) parseVFREEBUSY(comp *im.Component) (out *FreeBusy, err error, err2 error) {
	out = &FreeBusy{}
	for _, prop := range comp.Properties {
		switch prop.Name {
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

		case prContact:
			x, err := ToTextVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding altrep
			}
			out.Contact = &x

		case prDTStart:
			x, err := cp.ToDateTimeVal(prop, false)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding DTStart
			} else {
				out.DTStart = &x
			}

		case prDTEnd:
			x, err := cp.ToDateTimeVal(prop, false)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding DTEnd
			} else {
				out.DTEnd = &x
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

		case prURL:
			x, err := ToURIVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding URL
			} else {
				out.URL = &x
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

		case prFB:
			x, err := cp.ToFBVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding FreeBusy
			} else {
				out.FreeBusy = append(out.FreeBusy, &x)
			}

		case prReqStat:
			x, err := ToReqStatusVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding ReqStatus
			} else {
				out.RStatus = append(out.RStatus, &x)
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
