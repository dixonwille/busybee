package models

import "encoding/xml"
import "io"

//ResponseEnvelope is the highest level element.
//There is a single FreeBusyResponse for each user specified in the request.
type ResponseEnvelope struct {
	XMLName           xml.Name           `xml:"Envelope"`
	FreeBusyResponses []FreeBusyResponse `xml:"Body>GetUserAvailabilityResponse>FreeBusyResponseArray>FreeBusyResponse"`
}

//FreeBusyResponse holds all the Calendar Events for a single user.
type FreeBusyResponse struct {
	CalendarEvents []CalendarEvent `xml:"FreeBusyView>CalendarEventArray>CalendarEvent"`
}

//CalendarEvent is the start and end time in the requested time zone.
//Also returns whether they are busy or not at those times.
type CalendarEvent struct {
	StartTime string `xml:"StartTime"`
	EndTime   string `xml:"EndTime"`
	BusyType  string `xml:"BusyType"`
}

//Decode reads from reader and writes it into the model.
func (res *ResponseEnvelope) Decode(reader io.Reader) error {
	return xml.NewDecoder(reader).Decode(res)
}
