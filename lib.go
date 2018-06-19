package go_ical

import "github.com/soh335/icalparser"
import (
	"io"

	"strings"

	"github.com/mqus/go-ical/im"
)

type CalParser struct {
	inner           *icalparser.Parser
	currentCalendar *Calendar
	ComponentErrors chan<- error
	PropertyErrors  chan<- error
	ParameterErrors chan<- error
}

type ICalOutputter interface {
	OutputICal() string
}

func NewParser(reader io.Reader, componentErrors, propertyErrors, parameterErrors chan<- error) CalParser {
	return CalParser{
		inner:           icalparser.NewParser(reader),
		currentCalendar: nil,
		ComponentErrors: componentErrors,
		PropertyErrors:  propertyErrors,
		ParameterErrors: parameterErrors,
	}
}

func (cp *CalParser) ParseCalendar() (out *Calendar, err error) {
	cobj, err := cp.inner.Parse()
	if err != nil {
		return nil, err
	}
	calcomp, err2 := im.ToIntermediate(cobj)
	if err2 != nil {
		return nil, err2
	}

	out = &Calendar{Version: Version{"2.0", nil}}
	cp.currentCalendar = out
	for _, prop := range calcomp.Properties {
		switch strings.ToLower(prop.Name.C) {
		case prProdID:
			out.ProdID = ToStringVal(prop)

		case prVersion:
			out.Version = ToVersion(prop)

		case prCalScale:
			x := ToStringVal(prop)
			out.CalScale = &x

		case prMethod:
			x := ToStringVal(prop)
			out.Method = &x

		case prUid:
			x := ToStringVal(prop)
			out.Uid = &x

		case prLastMod:
			x, err := ToDateTimeVal(prop, true)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding LastMod
			} else {
				out.LastMod = &x
			}

		case prURL:
			x, err := ToURIVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding URL
			} else {
				out.URL = &x
			}
		case prRefresh:
			x, err := ToDurVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding the value
				//also: VALUE=DURATION is currently not required
			} else {
				out.Refresh = &x
			}

		case prSource:
			x, err := ToURIVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding SOURCE
			} else {
				out.Source = &x
			}
		case prColor:
			x := ToColorVal(prop)
			out.Color = &x
		case prName:
			x, err := ToTextVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding altrep
			}
			out.Name = append(out.Name, x)

		case prDesc:
			x, err := ToTextVal(prop)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding altrep
			}
			out.Desc = append(out.Desc, x)

		case prCat:
			out.Categories = append(out.Categories, ToCategVal(prop))
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
				out.Images = append(out.Images, x)
			}

		default:
			out.OtherProperties = append(out.OtherProperties, prop)

		}

	}

	//first parse all TimeZone components, to enable assigning accurate times for all components.
	for _, comp := range calcomp.Comps {
		if strings.ToLower(comp.Name) == vTZ {
			x, err, err2 := cp.parseVTIMEZONE(comp)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding TIMEZONE component
			} else {
				if err2 != nil {
					//MAYBE return err
					// instead of silently discarding subcomponents/properties
				}
				out.TimeZones[x.ID.Value] = x
			}
		}
	}

	//then all other components
	for _, comp := range calcomp.Comps {
		switch strings.ToLower(comp.Name) {
		case vEv:
			x, err, err2 := cp.parseVEVENT(comp)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding EVENT
			} else {
				if err2 != nil {
					//MAYBE return err
					// instead of silently discarding subcomponents/properties
				}
				out.Events = append(out.Events, x)
			}

		case vTODO:
			x, err, err2 := cp.parseVTODO(comp)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding EVENT
			} else {
				if err2 != nil {
					//MAYBE return err
					// instead of silently discarding subcomponents/properties
				}
				out.ToDos = append(out.ToDos, x)
			}

		case vJnl:
			x, err, err2 := cp.parseVJOURNAL(comp)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding JOURNAL
			} else {
				if err2 != nil {
					//MAYBE return err
					// instead of silently discarding subcomponents/properties
				}
				out.Journals = append(out.Journals, x)
			}

		case vFB:
			x, err, err2 := cp.parseVFREEBUSY(comp)
			if err != nil {
				//MAYBE return err
				// instead of silently discarding FREEBUSY
			} else {
				if err2 != nil {
					//MAYBE return err
					// instead of silently discarding subcomponents/properties
				}
				out.FreeBusy = append(out.FreeBusy, x)
			}

		case vTZ:
			//already parsed
		default:
			out.OtherComponents = append(out.OtherComponents, comp)
		}

	}
	return
}

func ToICal(in Calendar, out io.Writer) {
	//TODO
}
