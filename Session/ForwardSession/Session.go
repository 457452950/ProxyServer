package ForwardSession

import (
	"log"
	"net"
	"sync"
)

type Session struct {
	firstConnection  net.Conn
	secondConnection net.Conn
}

func (fs *Session) process() {
	defer fs.firstConnection.Close()
	defer fs.secondConnection.Close()

	ch := make(chan int)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go ReadAndSend(fs.firstConnection, fs.secondConnection, ch, wg)
	go ReadAndSend(fs.secondConnection, fs.firstConnection, ch, wg)

	wg.Wait()
	close(ch)
	log.Println("[forard] session close ", fs.firstConnection.RemoteAddr().String(), fs.secondConnection.RemoteAddr().String())
}

func (fs *Session) Start() {
	go fs.process()
}

func CreateForwardSession(conn net.Conn, addr string) (*Session, error) {
	con, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	log.Println("[forward] ", con.RemoteAddr().String(), conn.RemoteAddr().String())

	return &Session{
		firstConnection:  conn,
		secondConnection: con,
	}, nil
}

func GetForwardSession(conn net.Conn, con net.Conn) *Session {
	log.Println("[forward] ", con.RemoteAddr().String(), conn.RemoteAddr().String())
	return &Session{
		firstConnection:  conn,
		secondConnection: con,
	}
}
