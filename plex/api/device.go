package api

import (
	"encoding/xml"
)

type DeviceList struct {
	XMLName xml.Name `xml:"MediaContainer"`
	Devices []Device `xml:"Device"`
}

type Device struct {
	XMLName     xml.Name     `xml:"Device"`
	Roles       string       `xml:"provides,attr"`
	AccessToken string       `xml:"accessToken,attr"`
	Owned       bool         `xml:"owned,attr"`
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
