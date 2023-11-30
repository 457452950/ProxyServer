package LanServer

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"

	"ProxyServer/Session/ForwardSession"
)

type Server struct {
	frontAddr    string
	frontCtlAddr string
	backAddr     string

	controlConn net.Conn
}

func CreateServer(frontCtl string, frontIn string, back string) *Server {
	return &Server{
		frontCtlAddr: frontCtl,
		frontAddr:    frontIn,
		backAddr:     back,
	}
}

func (srv *Server) NewForwardSession() {
	log.Println("[lan] server new forward session")
	conn, err := net.Dial("tcp", srv.frontAddr)
	if err != nil {
		log.Println("[lan] connect to front ", err)
		return
	}

	log.Println("[lan] new forward session ")
	session, err := ForwardSession.CreateForwardSession(conn, srv.backAddr)
	if err != nil {
		log.Println("[lan] connect to backend ", err)
		return
	}
	session.Start()
}

func (srv *Server) AcceptProcess() {
	var data [4]byte
	var count int32
	for {
		_, err := srv.controlConn.Read(data[:])
		if err != nil {
			log.Println("[lan] ctl read ", err)
			break
		}
		bytesBuffer := bytes.NewBuffer(data[:])
		binary.Read(bytesBuffer, binary.BigEndian, &count)

		log.Printf("[lan] ctl want %d\n", count)

		for i := 0; i < int(count); i++ {
			srv.NewForwardSession()
		}
	}

}

func (srv *Server) Process() {
	var err error
	srv.controlConn, err = net.Dial("tcp", srv.frontCtlAddr)

	if err != nil {
		log.Println(err)
		return
	}

	srv.AcceptProcess()
}

func (srv *Server) Run() {
	go srv.Process()
}
