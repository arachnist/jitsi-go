module github.com/arachnist/jitsi-go

go 1.24.1

require gosrc.io/xmpp v0.5.1

require (
	github.com/google/uuid v1.1.1 // indirect
	golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7 // indirect
	nhooyr.io/websocket v1.6.5 // indirect
)

replace gosrc.io/xmpp => ../go-xmpp