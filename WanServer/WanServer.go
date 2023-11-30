package WanServer

import (
	"ProxyServer/Session/ForwardSession"
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"strings"
)

type Server struct {
	localOutAddr   string
	publicListener net.Listener

	localCtlAddr       string
	backAddr           string
	privateCtlListener net.Listener
	localInAddr        string
	privateListener    net.Listener
}

func CreateServer(localAddr string, localInAddr string, localCrlAddr string, back string) *Server {
	return &Server{
		localOutAddr: localAddr,
		localInAddr:  localInAddr,
		localCtlAddr: localCrlAddr,
		backAddr:     back,
	}
}

func (srv *Server) publicAcceptProcess(connCh chan<- net.Conn) {
	for {
		conn, err := srv.publicListener.Accept()
		if err != nil {
			log.Println("[wan] public accept ", err)
			break
		}

		log.Printf("[wan] new pub conn %s", conn.RemoteAddr().String())
		// todo:
		connCh <- conn
	}

}

func (srv *Server) privateAcceptCtlProcess(connCh <-chan net.Conn, connCh2 chan<- net.Conn) {
	var data [4]byte
	var count int32
	for {
		backconn, err := srv.privateCtlListener.Accept()
		if err != nil {
			log.Println("[wan] private accept to back end ", err)
			break
		}
		log.Printf("[wan] new pri ctl conn %s", backconn.RemoteAddr().String())

		if !strings.HasPrefix(backconn.RemoteAddr().String(), srv.backAddr) {
			backconn.Close()
			log.Printf("[wan] new pri ctl conn close")
			continue
		}

		for {
			conn := <-connCh
			if conn == nil {
				log.Println("[wan] err channel get nil net.Con")
				backconn.Close()
				return
			}

			count = 1
			bytesBuffer := bytes.NewBuffer(data[:])
			binary.Write(bytesBuffer, binary.BigEndian, count)

			_, err := backconn.Write(bytesBuffer.Bytes()[4:])
			if err != nil {
				log.Println("[wan] write to back end ", err)
				break
			}

			log.Printf("[wan] crl new session , ask backend")
			connCh2 <- conn
		}

		log.Printf("[wan] pri ctl conn close")
		backconn.Close()
	}

}

func (srv *Server) privateAcceptProcess(connCh <-chan net.Conn) {
	for {
		backconn, err := srv.privateListener.Accept()
		if err != nil {
			log.Println("[wan] private accept ", err)
			break
		}
		log.Printf("[wan] new backend conn %s", backconn.RemoteAddr().String())
		if !strings.HasPrefix(backconn.RemoteAddr().String(), srv.backAddr) {
			log.Println("[wan] private accept refuse")
			backconn.Close()
			continue
		}

		conn := <-connCh

		log.Printf("[wan] new forward sesion %s %s", backconn.RemoteAddr(), conn.RemoteAddr())
		sess := ForwardSession.GetForwardSession(conn, backconn)
		go sess.Start()
	}
}

func (srv *Server) Process() {
	var err error
	srv.publicListener, err = net.Listen("tcp", srv.localOutAddr)
	if err != nil {
		log.Println(err)
		return
	}

	srv.privateListener, err = net.Listen("tcp", srv.localInAddr)
	if err != nil {
		log.Println(err)
		return
	}

	srv.privateCtlListener, err = net.Listen("tcp", srv.localCtlAddr)
	if err != nil {
		log.Println(err)
		return
	}

	ch1 := make(chan net.Conn, 20)
	ch2 := make(chan net.Conn, 20)

	go srv.publicAcceptProcess(ch1)
	go srv.privateAcceptCtlProcess(ch1, ch2)
	srv.privateAcceptProcess(ch2)
}

func (srv *Server) Run() {
	go srv.Process()
}
