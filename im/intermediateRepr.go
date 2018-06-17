package im

import (
	"errors"

	"github.com/soh335/icalparser"
)

type Component struct {
	Name       string
	Comps      []*Component
	Properties []*icalparser.ContentLine
}

func ToIntermediate(obj *icalparser.Object) (out *Component, err error) {
	out = &Component{}
	out.Name = obj.HeaderLine.Value.C
	out.Properties = obj.PropertiyLines
	for _, comp := range obj.Components {
		x, e := toIntermediate(comp)
		if e != nil {
			return nil, e
		}
		out.Comps = append(out.Comps, x)
	}
	return
}

func toIntermediate(comp *icalparser.Component) (out *Component, err error) {
	out = &Component{}
	out.Name = comp.HeaderLine.Value.C
	props := comp.PropertiyLines
	for i := 0; i < len(props); i++ {
		if props[i].Name.C == "BEGIN" {
			newcomp := &Component{}
			newcomp.Name = props[i].Value.C
			x, e := parseComponent(props[i+1:], newcomp)
			if e != nil {
				return nil, e
			}
			i = i + x
			out.Comps = append(out.Comps, newcomp)
		} else {
			out.Properties = append(out.Properties, props[i])
		}
	}
	return
}

func parseComponent(props []*icalparser.ContentLine, comp *Component) (lines int, err error) {

	for i := 0; i < len(props); i++ {
		if props[i].Name.C == "BEGIN" {
			newcomp := &Component{}
			newcomp.Name = props[i].Value.C
			x, e := parseComponent(props[i+1:], newcomp)
			if e != nil {
				return 0, e
			}
			i = i + x
			comp.Comps = append(comp.Comps, newcomp)
		} else if props[i].Name.C == "END" {
			if props[i].Value.C != comp.Name {
				return 0, errors.New("Unexpected END:" + props[i].Value.C + " , expected END:" + comp.Name)
			} else {
				return i + 1, nil
			}
		} else { //Name = END
			comp.Properties = append(comp.Properties, props[i])
		}
	}
	return 0, errors.New("Expected END:" + comp.Name + " , found nothing")
}
