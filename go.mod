module github.com/Pavel7004/GraphPlot

go 1.15

replace (
	go.uber.org/atomic => go.uber.org/atomic v1.9.0
	golang.org/x/net => golang.org/x/net v0.0.0-20210428185706-aea814203247
	golang.org/x/sys => golang.org/x/sys v0.0.0-20200501145240-bc7a7d42d5c3
	gonum.org/v1/plot => gonum.org/v1/plot v0.9.0
)

require (
	github.com/Pavel7004/Common v0.0.0-20221027173841-9c6de7a5f6d7
	github.com/gin-gonic/gin v1.9.1
	github.com/gorilla/websocket v1.5.0
	github.com/kr/pretty v0.3.0 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/spf13/cobra v1.6.1
	github.com/spf13/viper v1.8.1
	github.com/tdewolff/canvas v0.0.0-20221118165558-64bc21122b8a
	gonum.org/v1/plot v0.11.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
