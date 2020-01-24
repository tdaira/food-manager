module github.com/tdaira/food-manager

go 1.13

require (
	cloud.google.com/go/storage v1.4.0
	github.com/creack/goselect v0.1.1 // indirect
	github.com/mattn/go-isatty v0.0.11 // indirect
	github.com/nats-io/jwt v0.3.2 // indirect
	github.com/nats-io/nats.go v1.9.1 // indirect
	github.com/urfave/cli v1.22.2 // indirect
	go.bug.st/serial.v1 v0.0.0-20191202182710-24a6610f0541 // indirect
	gobot.io/x/gobot v1.14.0
	golang.org/x/crypto v0.0.0-20191219195013-becbf705a915 // indirect
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553 // indirect
	golang.org/x/sys v0.0.0-20191220220014-0732a990476f // indirect
)

replace gobot.io/x/gobot v1.14.0 => github.com/tdaira/gobot v1.14.1-0.20191222055855-c4ed06d3e197
