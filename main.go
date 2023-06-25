package main

import (
	"myConnect/server"
	"myConnect/tlog"
)

func main() {
	tlog.Init()
	server.Init()
}
