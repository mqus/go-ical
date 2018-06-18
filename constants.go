package go_ical

const (
	//Components:
	vCal  = "vcalendar"
	vEv   = "vevent"
	vTODO = "vtodo"
	vJnl  = "vjournal"
	vFB   = "vfreebusy"
	vTZ   = "vtimezone"
	vAl   = "valarm"

	//Properties
	// calendar
	prProdID   = "prodid"
	prCalScale = "calscale"
	prMethod   = "method"
	prVersion  = "version"

	// descriptive
	prAttach   = "attach"
	prCat      = "categories"
	prClass    = "class"
	prComm     = "comment"
	prDesc     = "description"
	prGeo      = "geo"
	prLoc      = "location"
	prPctCompl = "percent-complete"
	prPrio     = "priority"
	prRes      = "resources"
	prStat     = "status"
	prSumm     = "summary"
	// date+time
	prCompl   = "completed"
	prDTEnd   = "dtend"
	prDue     = "due"
	prDTStart = "dtstart"
	prDur     = "duration"
	prFB      = "freebusy"
	prTTransp = "transp"
	// time zone
	prTZid   = "tzid"
	prTZName = "tzname"
	prTZoF   = "tzoffsetfrom"
	prTZoT   = "tzoffsetto"
	prTZURL  = "tzurl"
	// relations
	prAttendee  = "attendee"
	prContact   = "contact"
	prOrganizer = "organizer"
	prRid       = "recurrence-id"
	prRelTo     = "related-to"
	prURL       = "url"
	prUid       = "uid"
	// recurrence
	prExcDate = "exdate"
	prRDate   = "rdate"
	prRRule   = "rrule"
	// alarm
	prAction  = "action"
	prRepeat  = "repeat"
	prTrigger = "trigger"
	// change-management
	prCreated = "created"
	prDTStamp = "dtstamp"
	prLastMod = "last-modified"
	prSeq     = "sequence"
	// other
	prReqStat = "request-status"

	//Parameters
	paAltRep   = "altrep"
	paCN       = "cn"
	paCUType   = "cutype"
	paDelTo    = "delegated-to"
	paDelFrom  = "delegated-from"
	paDir      = "dir"
	paEnc      = "encoding"
	paFmtType  = "fmttype"
	paFBType   = "fbtype"
	paLang     = "language"
	paMbr      = "member"
	paPartStat = "partstat"
	paRange    = "range"
	paRel      = "related"
	paRelType  = "reltype"
	paRole     = "role"
	paRSVP     = "rsvp"
	paSentBy   = "sent-by"
	paTZid     = "tzid"
	paValue    = "value"

	//Value Types
	tBinary   = "binary"
	tBool     = "boolean"
	tCalAddr  = "cal-adress"
	tDate     = "date"
	tDateTime = "date-time"
	tTime     = "time"
	tDur      = "duration"
	tFloat    = "float"
	tInteger  = "integer"
	tPeriod   = "period"
	tRecur    = "recur"
	tText     = "text"
	tURI      = "uri"
	tUTCOffs  = "utc-offset"

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
