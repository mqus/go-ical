package go_ical

import (
	"strings"

	cl "github.com/mqus/go-contentline"

	"github.com/mqus/go-ical/util"
)

type TimeZone struct {
	//required
	ID TZidVal
	//required, at least one with daylight and standard each
	Offsets []*StandardDaylightTime

	//optional
	LastMod *DateTimeVal
	URL     *URIVal

	OtherProperties []*cl.Property
}

type StandardDaylightTime struct {
	IsDaylightTime bool
	//required
	DTStart    *DateTimeVal
	OffsetFrom UTCOffsetVal
	OffsetTo   UTCOffsetVal

	//optional
	//TODO must be merged
	RRule *RecurrentRule

	Comment []*TextVal
	RDate   []*RecurrenceSetVal
	//Constraint: no Altrep
	// for displaying TimeZone (e.g. EST)
	TZName []*TextVal

	OtherProperties []*cl.Property
}

type TZidVal struct {
	StringVal
	IsGloballyUnique bool
}

func ToTZidVal(line *cl.Property) (out TZidVal) {
	return TZidVal{ToStringVal(line), strings.HasPrefix(line.Value, "/")}
}

type UTCOffsetVal DurVal

func ToUTCOffsetVal(line *cl.Property) (out UTCOffsetVal, err error) {
	out = UTCOffsetVal{}
	out.OtherParam = line.Parameters
	out.Value, err = util.ParseUTCOffset(line.Value)
	return out, err
}

func (cp *CalParser) parseVTIMEZONE(comp *cl.Component) (out *TimeZone, err error, err2 error) {
	out = &TimeZone{}
	for _, prop := range comp.Properties {
		switch prop.Name {
		case prTZid:
			out.ID = ToTZidVal(prop)

		case prLastMod:
			x, err := cp.ToDateTimeVal(prop, true)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding LastMod
			} else {
				out.LastMod = &x
			}

		case prTZURL:
			x, err := ToURIVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding URL
			} else {
				out.URL = &x
			}

		default:
			out.OtherProperties = append(out.OtherProperties, prop)
		}

	}

	for _, subcomp := range comp.Comps {
		isDaylight := false
		switch subcomp.Name {
		case vDaylight:
			isDaylight = true
			fallthrough
		case vStandard:
			os, e1, e2 := cp.parseSDTime(subcomp, isDaylight)
			if e1 != nil {
				//MAYBE return err
				// instead of silently discarding DAYLIGHT/STANDARD component
			} else {
				if e2 != nil {
					//MAYBE return err2
					// instead of silently ignoring property
				}
				out.Offsets = append(out.Offsets, os)
			}
		default:
			//MAYBE don't silently discard other Components
			// out.OtherComponents = append(out.OtherComponents, subcomp)
		}
	}
	//TODO Conformance Checking
	return
}

func (cp *CalParser) parseSDTime(comp *cl.Component, isDaylight bool) (out *StandardDaylightTime, err error, err2 error) {
	out = &StandardDaylightTime{}
	for _, prop := range comp.Properties {
		switch prop.Name {
		case prDTStart:
			x, err := cp.ToDateTimeVal(prop, false)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding DTStart
			} else {
				out.DTStart = &x
			}

		case prTZoT:
			x, err := ToUTCOffsetVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding LastMod
			} else {
				out.OffsetTo = x
			}

		case prTZoF:
			x, err := ToUTCOffsetVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding LastMod
			} else {
				out.OffsetFrom = x
			}

		case prRRule:
			x := RecurrentRule(ToStringVal(prop))
			out.RRule = &x

		case prComm:
			x, err := ToTextVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding altrep
			}
			out.Comment = append(out.Comment, &x)

		case prRDate:
			x, err := cp.ToRecurrenceSetVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding ExceptionDate
			} else {
				out.RDate = append(out.RDate, &x)
			}

		case prTZName:
			x, err := ToTextVal(prop)
			if err != nil {
				//SHOULD NOT OCCUR, ALTREP IS NOT ALLOWED!
				//MAYBE return err
				// instead of silently discarding altrep
			}
			out.TZName = append(out.TZName, &x)

		default:
			out.OtherProperties = append(out.OtherProperties, prop)
		}

	}

	//MAYBE don't silently discard other Components
	// out.OtherComponents = comp.Comps

	//TODO Conformance Checking
	return
}
