# go-ical

a library for parsing icalendar formatted data

This library will try to conform with RFC 5545 and 7986.

This is a work in Progress:
### TODO:
- [ ] parsing of subcomponents
  - [x] VEVENT
    - [x] Properties
    - [x] VALARM
  - [x] VJOURNAL
  - [x] VTODO
  - [x] VFREEBUSY
  - [ ] VTIMEZONE
    - [x] implement VTIMEZONE
    - [ ] implement timezoning in other components
- [ ] replace structural parser
  - see RFC6868 (additional Characters and encoding)
  - see Section
- [ ] Implement tests
- [ ] Recurrence matching helper functions
- [ ] struct -> ical stream encoder
- [ ] Conformance checking
  - [ ] strictness settings
- (Un)escaping chars for TEXT vars (StringVar, TextVar,?)
### MAYBE:

- [ ] implement Venue Draft (https://tools.ietf.org/id/draft-norris-ical-venue-00.txt)
- [ ] implement jCal and xCal
- [ ] implement missing RFCs
## Relevant RFCs:

- 5545: current ICal Spec
  - [Errata](https://www.rfc-editor.org/errata_search.php?rfc=5545)
  - obsoletes RFC2445
- 7986: New Attributes for ICal (Image/Conference/Color)
- 6868: Additional Characters and encoding
- 7953: VAVAILABILITY
  - will not be implemented for now (open an Issue with your use case)
- 7529: iCalendar RScale extension
 - will not be implemented for now (open an Issue with your use case)