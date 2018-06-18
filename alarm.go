package go_ical

import (
	"strings"

	"github.com/mqus/go-ical/im"
	"github.com/soh335/icalparser"
)

// implements VALARM, specified in RFC5545, Section 3.6.6
type Alarm struct {
	//required in any case
	Action  StringVal
	Trigger TriggerVal
	//required for ACTION=DISP and ACTION=EMAIL
	Desc *TextVal
	//required for ACTION=EMAIL
	Summary  *TextVal
	Attendee []*PersonVal

	//only together
	DTDuration *DurVal
	Repeat     *IntVal

	//Optional, at most one if ACTION=AUDIO
	Attachments []*DataVal

	OtherProperties []*icalparser.ContentLine
}

func parseVALARM(comp *im.Component) (out *Alarm, err error, err2 error) {
	out = &Alarm{}
	for _, prop := range comp.Properties {
		switch strings.ToLower(prop.Name.C) {
		case prAction:
			out.Action = ToStringVal(prop)
		case prTrigger:
			x, err, err2 := ToTriggerVal(prop)
			if err != nil {
				return out, err, nil
			} else {
				if err2 != nil {
					//MAYBE return err2
					// instead of silently discarding parameters
				}
				out.Trigger = x
			}
		case prDesc:
			x, err := ToTextVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding altrep
			}
			out.Desc = &x
		case prSumm:
			x, err := ToTextVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding altrep
			}
			out.Summary = &x
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
		case prDur:
			x, err := ToDurVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding Duration
			} else {
				out.DTDuration = &x
			}
		case prRepeat:
			x, err := ToIntVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding Repeat
			} else {
				out.Repeat = &x
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

		default:
			out.OtherProperties = append(out.OtherProperties, prop)
		}
	}

	//TODO Conformance-checks:
	// check if dur is set <=> repeat is set
	// check if Desc/Summary/Attendee is set
	return
}
