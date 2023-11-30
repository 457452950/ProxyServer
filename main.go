package main

import (
	"ProxyServer/WanServer"
	"time"

	"ProxyServer/LanServer"
)

var (
	frontEndAddr    = "127.0.0.1:8500"
	frontCrlEndAddr = "127.0.0.1:8501"
	frontInEndAddr  = "127.0.0.1:8502"
	backAddr        = "127.0.0.1"

	backEndAddr = "192.168.0.196:8004"
	whiteList   = map[string]bool{"127.0.0.1": true}
)

func main() {
	frontSrv := WanServer.CreateServer(frontEndAddr, frontInEndAddr, frontCrlEndAddr, backAddr)
	frontSrv.Run()

	time.Sleep(100 * time.Millisecond)

	srv := LanServer.CreateServer(frontCrlEndAddr, frontInEndAddr, backEndAddr)
	srv.Process()
}
