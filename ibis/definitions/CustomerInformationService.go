package definitions

import (
	"encoding/xml"
)

// definitions for the CustomerInformationService interface v2

type GetAllDataResponse struct {
	XMLName xml.Name `xml:"CustomerInformationService.GetAllDataResponse"`
	AllData AllData  `xml:"AllData"`
}

type AllData struct {
	TimeStamp              Value           `xml:"TimeStamp"`
	VehicleRef             Value           `xml:"VehicleRef"`
	DefaultLanguage        Value           `xml:"DefaultLanguage"`
	TripInformation        TripInformation `xml:"TripInformation"`
	CurrentStopIndex       IntValue        `xml:"CurrentStopIndex"`
	RouteDeviation         string          `xml:"RouteDeviation"`
	DoorState              string          `xml:"DoorState"`
	VehicleStopRequested   BoolValue       `xml:"VehicleStopRequested"`
	ExitSide               string          `xml:"ExitSide"`
	MovingDirectionForward BoolValue       `xml:"MovingDirectionForward"`
	VehicleMode            string          `xml:"VehicleMode"`
	SpeakerActive          BoolValue       `xml:"SpeakerActive"`
	StopInformationActive  BoolValue       `xml:"StopInformationActive"`
}

type TripInformation struct {
	TripRef        Value        `xml:"TripRef"`
	StopSequence   StopSequence `xml:"StopSequence"`
	LocationState  string       `xml:"LocationState"`
	TimetableDelay Value        `xml:"TimetableDelay"`
}

type StopSequence struct {
	StopPoints []StopPoint `xml:"StopPoint"`
}

type StopPoint struct {
	StopIndex       Value            `xml:"StopIndex"`
	StopRef         Value            `xml:"StopRef"`
	StopName        Value            `xml:"StopName"`
	DisplayContents []DisplayContent `xml:"DisplayContent"`
	Connections     []Connection     `xml:"Connection"`
}

type DisplayContent struct {
	DisplayContentRef Value           `xml:"DisplayContentRef"`
	LineInformation   LineInformation `xml:"LineInformation"`
	Destination       Destination     `xml:"Destination"`
}

type LineInformation struct {
	LineRef       Value `xml:"LineRef"`
	LineName      Value `xml:"LineName"`
	LineShortName Value `xml:"LineShortName"`
	LineNumber    Value `xml:"LineNumber"`
}

type Destination struct {
	DestinationRef  Value `xml:"DestinationRef"`
	DestinationName Value `xml:"DestinationName,omitempty"`
}

type Connection struct {
	StopRef        Value          `xml:"StopRef"`
	ConnectionRef  Value          `xml:"ConnectionRef"`
	ConnectionType string         `xml:"ConnectionType"`
	DisplayContent DisplayContent `xml:"DisplayContent"`
}

type Value struct {
	Value    string `xml:"Value"`
	Language string `xml:"Language,omitempty"`
}

type IntValue struct {
	Value int `xml:"Value"`
}

type BoolValue struct {
	Value bool `xml:"Value"`
}
