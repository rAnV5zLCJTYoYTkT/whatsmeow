module go.mau.fi/whatsmeow

go 1.21

require (
	github.com/gorilla/websocket v1.5.0
	go.mau.fi/libsignal v0.1.0
	go.mau.fi/util v0.4.1
	google.golang.org/protobuf v1.33.0
	philippgille.com/chromem-go v0.6.0
)

require (
	github.com/mattn/go-sqlite3 v1.14.22
	golang.org/x/crypto v0.25.0
	golang.org/x/net v0.27.0
)

// Personal fork - keeping dependencies up to date
// upstream: https://github.com/tulir/whatsmeow
// Note: bumped x/crypto and x/net to latest patch versions (security updates)
// TODO: look into replacing chromem-go with a lighter embedding solution
// TODO: upgrade to go 1.22 once libsignal and util deps are confirmed compatible
