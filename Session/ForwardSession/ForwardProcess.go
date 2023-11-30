package ForwardSession

import (
	"log"
	"net"
	"sync"
)

func ReadAndSend(readConn net.Conn, sendConn net.Conn, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	var buf [128]byte
	for {
		select {
		case <-ch:
			return

		default:
			n, err := readConn.Read(buf[:])
			if err != nil {
				log.Println("[forward] read ", readConn.RemoteAddr().String(), err)
				return
			}

			//log.Println(string(buf[:n]))

			n, err = sendConn.Write(buf[:n])
			if err != nil {
				log.Println("[forward] write ", sendConn.RemoteAddr().String(), err)
				return
			}
		}

	}
}
