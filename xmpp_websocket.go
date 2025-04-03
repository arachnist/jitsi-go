package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
)

func init() {
	flag.StringVar(&server, "server", "meet.ffmuc.net", "jitsi server address")
	flag.StringVar(&room, "room", "hackcenter-test", "jitsi room to join")
	flag.StringVar(&nickname, "nickname", "notbot", "nick to use")
}

func main() {
	flag.Parse()

	wsAddr := fmt.Sprintf("wss://%s/xmpp-websocket", server)
	jid := fmt.Sprintf("%s@%s", nickname, server)
	config := &xmpp.Config{
		TransportConfiguration: xmpp.TransportConfiguration{
			Address: wsAddr, // ?room=hackcenter-0fa0asdhgkjds",
		},
		Jid:          jid,
		Credential:   xmpp.Anonymous(),
		StreamLogger: os.Stdout,
	}

	router := xmpp.NewRouter()
	router.HandleFunc("message", handleMessage)

	client, err := xmpp.NewClient(config, router, errHandler)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	client.PostConnectHook = getPostConnectFunc(client)

	cm := xmpp.NewStreamManager(client, nil)
	log.Fatal(cm.Run())
}
/* generated:
 * <iq xmlns="jabber:client" type="set" id="b03e99be-5e64-4dc2-ab48-11841e7d50de" to="focus.meet.ffmuc.net">
 * 	<conference xmlns="http://jitsi.org/protocol/focus" machine-uid="cc5874acc959cc1728cc0e0a69725589" room="hackcenter-test@conference.meet.ffmuc.net/wnyU0M8nJj4l">
 * 		<property name="startAudioMuted" value="25"></property>
 * 		<property name="startVideoMuted" value="25"></property>
 * 		<property name="rtcstatsEnabled" value="false"></property>
 * 		<property name="visitors-version" value="1"></property>
 * 	</conference>
 * </iq>
 */
/* actual client: ->
 * <iq id="1cb0086a-7087-4d46-820e-65ab6cf4b50a:sendIQ" to="focus.meet.ffmuc.net" type="set" xmlns="jabber:client">
 * 	<conference machine-uid="cc5874acc959cc1728cc0e0a69725589" room="hackcenter-test@conference.meet.ffmuc.net" xmlns="http://jitsi.org/protocol/focus">
 * 		<property name="startAudioMuted" value="25"/>
 * 		<property name="startVideoMuted" value="25"/>
 * 		<property name="rtcstatsEnabled" value="false"/>
 * 		<property name="visitors-version" value="1"/>
 * 	</conference>
 * </iq>
 */
/* <-
 * <iq xmlns='jabber:client' id='1cb0086a-7087-4d46-820e-65ab6cf4b50a:sendIQ' xml:lang='en-US' to='9c93d8e3-53db-49e5-8030-9678448cb193@meet.ffmuc.net/bIFlCbOe8e7B' type='result' from='focus.meet.ffmuc.net'>
 * 	<conference focusjid='focus@auth.meet.ffmuc.net' room='hackcenter-test@conference.meet.ffmuc.net' xmlns='http://jitsi.org/protocol/focus' ready='true'>
 * 		<property name='authentication' value='false'/>
 * 	</conference>
 * </iq>
 */
/* ->
 * <presence to="hackcenter-test@conference.meet.ffmuc.net/9c93d8e3" xmlns="jabber:client">
 * 	<x xmlns="http://jabber.org/protocol/muc"/>
 * 	<stats-id>Russ-xho</stats-id>
 * 	<c hash="sha-1" node="https://jitsi.org/jitsi-meet" ver="W+eXiSPXEn8io7DtaLMtt3J13E4=" xmlns="http://jabber.org/protocol/caps"/>
 * 	<SourceInfo>{}</SourceInfo>
 * 	<jitsi_participant_region>ffmuc-de1</jitsi_participant_region>
 * 	<jitsi_participant_codecList>vp8,vp9</jitsi_participant_codecList>
 * 	<nick xmlns="http://jabber.org/protocol/nick">aritest</nick>
 * </presence>
 */
func getPostConnectFunc(client *xmpp.Client) func() error {
	return func() error {
		jid, err := stanza.NewJid(client.Session.BindJid)
		if err != nil {
			return err
		}
		roomJid := fmt.Sprintf("%s@conference.%s/%s", room, server, jid.Resource)

		conferenceReq, err := stanza.NewIQ(stanza.Attrs{
			To:   fmt.Sprintf("focus.%s", server),
			Type: stanza.IQTypeSet,
		})
		if err != nil {
			return err
		}
		conferenceReq.Payload = &Conference{
			MachineUID: "cc5874acc959cc1728cc0e0a69725589",
			Room: roomJid,
			Properties: []Property{
				{Name: "startAudioMuted", Value: "25"},
				{Name: "startVideoMuted", Value: "25"},
				{Name: "rtcstatsEnabled", Value: "false"},
				{Name: "visitors-version", Value: "1"},
			},
		}

		fmt.Printf("%#v\n", conferenceReq)

		ctx, _ := context.WithTimeout(context.Background(), 30 * time.Second)
		stanzaCh, err := client.SendIQ(ctx, conferenceReq)
		if err != nil {
			return err
		}

		<-stanzaCh

		presence := stanza.Presence{
			Attrs: stanza.Attrs{
				To: roomJid,
			},
			Extensions: []stanza.PresExtension{
				stanza.MucPresence{},
				StatsID{Value: "Joy-4gA"},
				Region{Id: "ffmuc-de1"},
				stanza.Caps{
					Hash: "sha-1",
					Node: "https://jitsi.org/jitsi-meet",
					Ver:  "Bb/FVrXF6cseTgvTP/PiwlsETz4=",
				},
				JitsiParticipantRegion{
					Value: "ffmuc-de1",
				},
				VideoMuted{Value: "true"},
				AudioMuted{Value: "true"},
				JitsiParticipantCodecType{Value: ""},
				Nick{Value: nickname},
			},
		}

		err = client.Send(presence)
		if err != nil {
			return err
		}

		return nil
	}
}

func errHandler(e error) {
	log.Println("xmpp router error:", e)
}

func handleMessage(s xmpp.Sender, p stanza.Packet) {
	msg, ok := p.(stanza.Message)
	if !ok {
		_, _ = fmt.Fprintf(os.Stdout, "Ignoring packet: %T\n", p)
		return
	}

	_, _ = fmt.Fprintf(os.Stdout, "Body = %s - from = %s\n", msg.Body, msg.From)
}

type PostConnectHook func() error

var (
	server, room, nickname string
)
