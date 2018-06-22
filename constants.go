package go_ical

const (
	//Components:
	vCal      = "VCALENDAR"
	vEv       = "VEVENT"
	vTODO     = "VTODO"
	vJnl      = "VJOURNAL"
	vFB       = "VFREEBUSY"
	vTZ       = "VTIMEZONE"
	vAl       = "VALARM"
	vDaylight = "DAYLIGHT"
	vStandard = "STANDARD"

	//Properties
	// calendar
	prProdID   = "PRODID"
	prCalScale = "CALSCALE"
	prMethod   = "METHOD"
	prVersion  = "VERSION"

	// descriptive
	prAttach   = "ATTACH"
	prCat      = "CATEGORIES"
	prClass    = "CLASS"
	prComm     = "COMMENT"
	prDesc     = "DESCRIPTION"
	prGeo      = "GEO"
	prLoc      = "LOCATION"
	prPctCompl = "PERCENT-COMPLETE"
	prPrio     = "PRIORITY"
	prRes      = "RESOURCES"
	prStat     = "STATUS"
	prSumm     = "SUMMARY"
	// date+time
	prCompl   = "COMPLETED"
	prDTEnd   = "DTEND"
	prDue     = "DUE"
	prDTStart = "DTSTART"
	prDur     = "DURATION"
	prFB      = "FREEBUSY"
	prTTransp = "TRANSP"
	// time zone
	prTZid   = "TZID"
	prTZName = "TZNAME"
	prTZoF   = "TZOFFSETFROM"
	prTZoT   = "TZOFFSETTO"
	prTZURL  = "TZURL"
	// relations
	prAttendee  = "ATTENDEE"
	prContact   = "CONTACT"
	prOrganizer = "ORGANIZER"
	prRid       = "RECURRENCE-ID"
	prRelTo     = "RELATED-TO"
	prURL       = "URL"
	prUid       = "UID"
	// recurrence
	prExcDate = "EXDATE"
	prRDate   = "RDATE"
	prRRule   = "RRULE"
	// alarm
	prAction  = "ACTION"
	prRepeat  = "REPEAT"
	prTrigger = "TRIGGER"
	// change-management
	prCreated = "CREATED"
	prDTStamp = "DTSTAMP"
	prLastMod = "LAST-MODIFIED"
	prSeq     = "SEQUENCE"
	// other
	prReqStat = "REQUEST-STATUS"

	//Parameters
	paAltRep   = "ALTREP"
	paCN       = "CN"
	paCUType   = "CUTYPE"
	paDelTo    = "DELEGATED-TO"
	paDelFrom  = "DELEGATED-FROM"
	paDir      = "DIR"
	paEnc      = "ENCODING"
	paFmtType  = "FMTTYPE"
	paFBType   = "FBTYPE"
	paLang     = "LANGUAGE"
	paMbr      = "MEMBER"
	paPartStat = "PARTSTAT"
	paRange    = "RANGE"
	paTrRel    = "RELATED"
	paRelType  = "RELTYPE"
	paRole     = "ROLE"
	paRSVP     = "RSVP"
	paSentBy   = "SENT-BY"
	paTZid     = "TZID"
	paValue    = "VALUE"

	//Value Types
	tBinary   = "BINARY"
	tBool     = "BOOLEAN"
	tCalAddr  = "CAL-ADRESS"
	tDate     = "DATE"
	tDateTime = "DATE-TIME"
	tTime     = "TIME"
	tDur      = "DURATION"
	tFloat    = "FLOAT"
	tInteger  = "INTEGER"
	tPeriod   = "PERIOD"
	tRecur    = "RECUR"
	tText     = "TEXT"
	tURI      = "URI"
	tUTCOffs  = "UTC-OFFSET"

	//CalendarUser Types
	cuIndiv = "individual"
	cuGroup = "group"
	cuRes   = "resource"
	cuRoom  = "room"
	cuUnkn  = "unknown"

	//Free/Busy Types
	fbFree    = "free"
	fbBusy    = "busy"
	fbBusUnav = "busy-unavailable"
	fbBusTent = "busy-tentative"

	//Participation Status Types
	psNAction = "needs-action"
	psAcc     = "accepted"
	psDecl    = "declined"
	psTent    = "tentative"
	psDeleg   = "delegated"
	psCompl   = "completed"
	psInProc  = "in-process"

	//Relationship types
	rtChild   = "child"
	rtParent  = "parent"
	rtSibling = "sibling"

	//Participation Role Types
	ptChair = "chair"
	ptRPart = "req-participant"
	ptOPart = "opt-participant"
	ptNPart = "non-participant"

	//Actions
	aAudio   = "audio"
	aDisplay = "display"
	aEMail   = "email"

	//Classifications
	cPublic  = "public"
	cPrivate = "private"
	cConf    = "confidential"

	//from RFC7986:
	prName    = "name"
	prRefresh = "refresh-interval"
	prSource  = "source"
	prColor   = "color"
	prImage   = "image"
	prConf    = "conference"

	paDisp    = "display"
	paEMail   = "email"
	paFeature = "feature"
	paLabel   = "label"

	//Display Types
	diBadge    = "badge"
	diGraphic  = "graphic"
	diFullsize = "fullsize"
	diThumb    = "thumbnail"

	//Feature Types
	feAudio  = "audio"
	feChat   = "chat"
	feFeed   = "feed"
	feMod    = "moderator"
	fePhone  = "phone"
	feScreen = "screen"
	feVid    = "video"
)
