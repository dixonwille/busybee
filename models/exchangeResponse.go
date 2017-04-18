package models

// <?xml version="1.0" encoding="UTF-8"?>
// <s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
//    <s:Header>
//       <h:ServerVersionInfo xmlns:h="http://schemas.microsoft.com/exchange/services/2006/types" xmlns="http://schemas.microsoft.com/exchange/services/2006/types" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" MajorVersion="14" MinorVersion="3" MajorBuildNumber="319" MinorBuildNumber="2" />
//    </s:Header>
//    <s:Body xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
//       <GetUserAvailabilityResponse xmlns="http://schemas.microsoft.com/exchange/services/2006/messages">
//          <FreeBusyResponseArray>
//             <FreeBusyResponse>
//                <ResponseMessage ResponseClass="Success">
//                   <ResponseCode>NoError</ResponseCode>
//                </ResponseMessage>
//                <FreeBusyView>
//                   <FreeBusyViewType xmlns="http://schemas.microsoft.com/exchange/services/2006/types">FreeBusy</FreeBusyViewType>
//                   <CalendarEventArray xmlns="http://schemas.microsoft.com/exchange/services/2006/types">
//                      <CalendarEvent>
//                         <StartTime>2017-04-18T09:30:00</StartTime>
//                         <EndTime>2017-04-18T09:45:00</EndTime>
//                         <BusyType>Busy</BusyType>
//                      </CalendarEvent>
//                      <CalendarEvent>
//                         <StartTime>2017-04-18T10:30:00</StartTime>
//                         <EndTime>2017-04-18T11:00:00</EndTime>
//                         <BusyType>Busy</BusyType>
//                      </CalendarEvent>
//                      <CalendarEvent>
//                         <StartTime>2017-04-18T16:00:00</StartTime>
//                         <EndTime>2017-04-18T16:30:00</EndTime>
//                         <BusyType>Busy</BusyType>
//                      </CalendarEvent>
//                   </CalendarEventArray>
//                   <WorkingHours xmlns="http://schemas.microsoft.com/exchange/services/2006/types">
//                      <TimeZone>
//                         <Bias>300</Bias>
//                         <StandardTime>
//                            <Bias>0</Bias>
//                            <Time>02:00:00</Time>
//                            <DayOrder>1</DayOrder>
//                            <Month>11</Month>
//                            <DayOfWeek>Sunday</DayOfWeek>
//                         </StandardTime>
//                         <DaylightTime>
//                            <Bias>-60</Bias>
//                            <Time>02:00:00</Time>
//                            <DayOrder>2</DayOrder>
//                            <Month>3</Month>
//                            <DayOfWeek>Sunday</DayOfWeek>
//                         </DaylightTime>
//                      </TimeZone>
//                      <WorkingPeriodArray>
//                         <WorkingPeriod>
//                            <DayOfWeek>Monday Tuesday Wednesday Thursday Friday</DayOfWeek>
//                            <StartTimeInMinutes>480</StartTimeInMinutes>
//                            <EndTimeInMinutes>1020</EndTimeInMinutes>
//                         </WorkingPeriod>
//                      </WorkingPeriodArray>
//                   </WorkingHours>
//                </FreeBusyView>
//             </FreeBusyResponse>
//          </FreeBusyResponseArray>
//       </GetUserAvailabilityResponse>
//    </s:Body>
// </s:Envelope>
