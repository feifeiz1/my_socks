package main

import (
	"bufio"
	"github.com/feifeiz1/my_socks/socks5"
	"log"
	"net"
	"time"
)

func main() {
	l, err := net.Listen("tcp4", ":13255")
	if err != nil {
		log.Fatalf("net.Listen failed,err:%v\n", err)
	}
	log.Println("Server Listen on:127.0.0.1:13255")
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("[ERROR]accept connection failed,err:%v\n", err)
			continue
		}
		go handlerConn(conn)
	}
}

func handlerConn(c net.Conn) {
	defer c.Close()
	reader := bufio.NewReader(c)

	timeOut := time.Now().Add(time.Millisecond * 500)
	c.SetReadDeadline(timeOut)
	//c.SetWriteDeadline(timeOut)

	if err := socks5.Auth(reader, c); err != nil {
		log.Printf("auth failed,err:%v", err)
		return
	}

	if err := socks5.Connect(reader, c); err != nil {
		log.Printf("connect failed,err:%v", err)
		return
	}

	//for {
	//	b, err := reader.ReadByte()
	//	if err != nil {
	//		log.Printf("[ERROR]connection.Read failed,err:%v\n", err)
	//		return
	//	}
	//	_, err = c.Write([]byte{b})
	//	if err != nil {
	//		log.Printf("[ERROR]connection.Write failed,err:%v\n", err)
	//		break
	//	}
	//}
}
