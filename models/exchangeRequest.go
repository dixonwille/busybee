package models

import (
	"encoding/xml"
	"io"
	"math"
	"time"
)

const (
	xsi            = "http://www.w3.org/2001/XMLSchema-instance"
	xsd            = "http://www.w3.org/2001/XMLSchema"
	soap           = "http://schemas.xmlsoap.org/soap/envelope/"
	xmlnsTypes     = "http://schemas.microsoft.com/exchange/services/2006/types"
	xmlnsMessages  = "http://schemas.microsoft.com/exchange/services/2006/messages"
	dateTimeFormat = "2006-01-02T15:04:05"
	timeFormat     = "15:04:05"
)

//RequestEnvelope is the main structure of our exchange soap request.
type RequestEnvelope struct {
	XMLName                    xml.Name                   `xml:"soap:Envelope"`
	Xsi                        string                     `xml:"xmlns:xsi,attr"`
	Xsd                        string                     `xml:"xmlns:xsd,attr"`
	Soap                       string                     `xml:"xmlns:soap,attr"`
	T                          string                     `xml:"xmlns:t,attr"`
	GetUserAvailabilityRequest GetUserAvailabilityRequest `xml:"soap:Body>GetUserAvailabilityRequest"`
}

//NewRequestEnvelope creates the entire request for us.
//It is highly recommended to use this method first then update values as needed.
//There are a lot of default values that this creates which is enough for what we need it to do.
func NewRequestEnvelope(startTime, endTime time.Time, addresses ...string) *RequestEnvelope {
	return &RequestEnvelope{
		Xsi:  xsi,
		Xsd:  xsd,
		Soap: soap,
		T:    xmlnsTypes,
		GetUserAvailabilityRequest: newGetUserAvailabilityRequest(startTime, endTime, addresses...),
	}
}

//GetUserAvailabilityRequest is the soap request we want to make.
type GetUserAvailabilityRequest struct {
	Xmlns               string              `xml:"xmlns,attr"`
	T                   string              `xml:"xmlns:t,attr"`
	TimeZone            TimeZone            `xml:"t:TimeZone"`
	MailboxDataArray    MailboxDataArray    `xml:"MailboxDataArray"`
	FreeBusyViewOptions FreeBusyViewOptions `xml:"t:FreeBusyViewOptions"`
}

func newGetUserAvailabilityRequest(startTime, endTime time.Time, addresses ...string) GetUserAvailabilityRequest {
	return GetUserAvailabilityRequest{
		Xmlns:               xmlnsMessages,
		T:                   xmlnsTypes,
		TimeZone:            newTimeZone(),
		FreeBusyViewOptions: newFreeBusyViewOptions(startTime, endTime),
		MailboxDataArray:    newMailboxDataArray(addresses...),
	}
}

//TimeZone is used to tell the server information about our timezone.
type TimeZone struct {
	Xmlns        string       `xml:"xmlns,attr"`
	Bias         int          `xml:"Bias"`
	StandardTime TimeZoneTime `xml:"StandardTime"`
	DaylightTime TimeZoneTime `xml:"DaylightTime"`
}

func newTimeZone() TimeZone {
	now := time.Now()
	std := time.Date(now.Year(), time.November, nthSunday(now.Year(), time.November, 1), 2, 0, 0, 0, time.Local)
	day := time.Date(now.Year(), time.March, nthSunday(now.Year(), time.March, 2), 2, 0, 0, 0, time.Local)
	_, offsetSeconds := now.Zone()
	offset := time.Duration(offsetSeconds) * time.Second
	if now.After(day) && now.Before(std) {
		offset = offset - (time.Duration(60) * time.Minute)
	}
	return TimeZone{
		Xmlns:        xmlnsTypes,
		Bias:         -1 * int(math.Ceil(offset.Minutes())),
		StandardTime: newTimeZoneTime(0, 1, int(std.Month()), std, time.Sunday.String()),
		DaylightTime: newTimeZoneTime(-60, 2, int(day.Month()), day, time.Sunday.String()),
	}
}

//TimeZoneTime is shared structure for StandardTime and DaylightTime.
type TimeZoneTime struct {
	Bias      int    `xml:"Bias"`
	Time      string `xml:"Time"`
	DayOrder  int    `xml:"DayOrder"`
	Month     int    `xml:"Month"`
	DayOfWeek string `xml:"DayOfWeek"`
}

func newTimeZoneTime(offsetMinutes, dayOrder, month int, when time.Time, dayOfWeek string) TimeZoneTime {
	return TimeZoneTime{
		Bias:      offsetMinutes,
		Time:      when.Format(timeFormat),
		DayOrder:  dayOrder,
		Month:     month,
		DayOfWeek: dayOfWeek,
	}
}

//MailboxDataArray is the overall structure for each individual email address.
type MailboxDataArray struct {
	MailboxData []MailboxData `xml:"t:MailboxData"`
}

func newMailboxDataArray(addresses ...string) MailboxDataArray {
	mda := new(MailboxDataArray)
	for _, address := range addresses {
		mda.MailboxData = append(mda.MailboxData, newMailboxData(address))
	}
	return *mda
}

//MailboxData is a single entry in MailboxDataArray.
//It holds information like which user you want to look up.
type MailboxData struct {
	Address          string `xml:"t:Email>t:Address"`
	AttendeeType     string `xml:"t:AttendeeType"`
	ExcludeConflicts bool   `xml:"t:ExcludeConflicts"`
}

func newMailboxData(address string) MailboxData {
	return MailboxData{
		Address:          address,
		AttendeeType:     "Required",
		ExcludeConflicts: false,
	}
}

//FreeBusyViewOptions lets the server know how you want to view the data that will be returned.
//I only care if there is something going on now and if you are busy for it.
//So I return the minimum required.
type FreeBusyViewOptions struct {
	StartTime                       string `xml:"t:TimeWindow>t:StartTime"`
	EndTime                         string `xml:"t:TimeWindow>t:EndTime"`
	MergedFreeBusyIntervalInMinutes int    `xml:"t:MergedFreeBusyIntervalInMinutes"`
	RequestedView                   string `xml:"t:RequestedView"`
}

func newFreeBusyViewOptions(startTime, endTime time.Time) FreeBusyViewOptions {
	return FreeBusyViewOptions{
		StartTime: startTime.Format(dateTimeFormat),
		EndTime:   endTime.Format(dateTimeFormat),
		MergedFreeBusyIntervalInMinutes: 30,
		RequestedView:                   "FreeBusy",
	}
}

//Encode turns this request into a soap+xml that the server can understand.
func (req *RequestEnvelope) Encode(writer io.Writer) error {
	_, err := writer.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	return xml.NewEncoder(writer).Encode(req)
}

func nthSunday(year int, month time.Month, nth int) int {
	t := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	return (7-int(t.Weekday()))%7 + (1 + (nth-1)*7)
}
