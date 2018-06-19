# go-ical

a library for parsing icalendar formatted data

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
- [ ] Implement tests
- [ ] Recurrence matching helper functions
- [ ] struct -> ical stream encoder
- [ ] Conformance checking
  - [ ] strictness settings

### MAYBE:
- [ ] implement Venue Draft (https://tools.ietf.org/id/draft-norris-ical-venue-00.txt)
