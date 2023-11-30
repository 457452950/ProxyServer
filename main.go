package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

var (
	backEndHost = "mylan.ubuntu01.com"
	backEndPort = 8004
	whiteList   = map[string]bool{"127.0.0.1": true}
)

func readAndSend(readConn net.Conn, sendConn net.Conn) {
	var buf [128]byte
	for {
		n, err := readConn.Read(buf[:])
		if err != nil {
			log.Println(err)
			return
		}

		n, err = sendConn.Write(buf[:n])
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func frontEndProcess(frontConn net.Conn) {
	defer frontConn.Close() // 关闭连接

	// check remote host
	remoteAddr := frontConn.RemoteAddr().String()
	remoteHost := strings.Split(remoteAddr, ":")[0]
	if whiteList[remoteHost] == false {
		log.Printf("banded host : %s \n", remoteAddr)
		return
	}

	backEndConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", backEndHost, backEndPort))
	if err != nil {
		fmt.Println("err : ", err)
		return
	}
	defer backEndConn.Close() // 关闭TCP连接

	go readAndSend(frontConn, backEndConn)
	readAndSend(backEndConn, frontConn)
}

func main() {
	listen, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println("Listen() failed, err: ", err)
		return
	}
	for {
		conn, err := listen.Accept() // 监听客户端的连接请求
		if err != nil {
			fmt.Println("Accept() failed, err: ", err)
			continue
		}
		go frontEndProcess(conn) // 启动一个goroutine来处理客户端的连接请求
	}
}
