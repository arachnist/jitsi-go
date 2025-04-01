package main

import (
	"fmt"
	"log"
	"os"

	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
)

func main() {
	config := &xmpp.Config{
		TransportConfiguration: xmpp.TransportConfiguration{
			Address: "wss://meet.ffmuc.net/xmpp-websocket?room=hackcenter-0fa0asdhgkjds",
		},
		Jid:          "test@meet.ffmuc.net",
		Credential:   xmpp.Anonymous(),
		StreamLogger: os.Stdout,
		Insecure:     true,
	}

	router := xmpp.NewRouter()
	router.HandleFunc("message", handleMessage)

	client, err := xmpp.NewClient(config, router, errHandler)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	// If you pass the client to a connection manager, it will handle the reconnect policy
	// for you automatically.
	cm := xmpp.NewStreamManager(client, nil)
	log.Fatal(cm.Run())
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
