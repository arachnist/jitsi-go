package main

import (
	"encoding/xml"
	"gosrc.io/xmpp/stanza"
)

type StatsID struct {
	XMLName xml.Name `xml:"stats-id"`
	Value   string   `xml:",innerxml"`
}

type Region struct {
	XMLName xml.Name `xml:"http://jitsi.org/jitsi-meet region"`
	Id      string   `xml:"id,attr"`
}

type JitsiParticipantRegion struct {
	XMLName xml.Name `xml:"jitsi_participant_region"`
	Value   string   `xml:",innerxml"`
}

type VideoMuted struct {
	XMLName xml.Name `xml:"videomuted"`
	Value   string   `xml:",innerxml"`
}

type AudioMuted struct {
	XMLName xml.Name `xml:"audiomuted"`
	Value   string   `xml:",innerxml"`
}

type JitsiParticipantCodecType struct {
	XMLName xml.Name `xml:"jitsi_participant_codecType"`
	Value   string   `xml:",innerxml"`
}

type Nick struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/nick nick"`
	Value   string   `xml:",innerxml"`
}

type Property struct {
	XMLName xml.Name `xml:"property"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}

type Conference struct {
	XMLName    xml.Name   `xml:"http://jitsi.org/protocol/focus conference"`
	MachineUID string     `xml:"machine-uid,attr"`
	Room       string     `xml:"room,attr"`
	Properties []Property `xml:",innerxml"`
	ResultSet  *stanza.ResultSet `xml:"set,omitempty"`
}

func (c *Conference) GetSet() *stanza.ResultSet {
	return c.ResultSet
}

func (c *Conference) Namespace() string {
	return c.XMLName.Space
}