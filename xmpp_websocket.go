package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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

func getPostConnectFunc(client *xmpp.Client) func() error {
	return func() error {
		jid, err := stanza.NewJid(client.Session.BindJid)
		if err != nil {
			return err
		}

		roomJid := fmt.Sprintf("%s@conference.%s/%s", room, server, jid.Resource)

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
