package api

import "encoding/xml"

// For reasons unknown to me the API endpoint that returns devices
// only supports XML, even when specifying a "Accept: application/json"
// header.
type DeviceList struct {
	XMLName xml.Name `xml:"MediaContainer"`
	Devices []Device `xml:"Device"`
}

type Device struct {
	XMLName     xml.Name     `xml:"Device"`
	Name        string       `xml:"name,attr"`
	ClientID    string       `xml:"clientIdentifier,attr"`
	Roles       string       `xml:"provides,attr"`
	AccessToken string       `xml:"accessToken,attr"`
	Product     string       `xml:"product,attr"`
	Version     string       `xml:"productVersion,attr"`
	Platform    string       `xml:"platform,attr"`
	Connections []Connection `xml:"Connection"`
}

type Connection struct {
	XMLName  xml.Name `xml:"Connection"`
	Protocol string   `xml:"protocol,attr"`
	Address  string   `xml:"address,attr"`
	Port     int      `xml:"port,attr"`
	URI      string   `xml:"uri,attr"`
	Local    bool     `xml:"local,attr"`
}
