package main

import (
	"fmt"
	"github.com/feifeiz1/my_socks/types"
	"io"
	"net"
)

func main() {
	c, err := net.Dial("tcp4", "127.0.0.1:13255")
	if err != nil {
		panic(err)
	}

	_, err = c.Write([]byte{types.Socks5Ver, 0x02, 0x00})
	if err != nil {
		panic("write err" + err.Error())
	}
	b, err := io.ReadAll(c)
	if err != nil {
		panic("read err:" + err.Error())
	}
	fmt.Println(b)
}
