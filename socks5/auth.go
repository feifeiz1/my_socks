package socks5

import (
	"bufio"
	"fmt"
	"github.com/feifeiz1/my_socks/types"
	"io"
	"log"
	"net"
)

// Auth auth认证阶段
// +----+----------+----------+
// |VER | NMETHODS | METHODS  |
// +----+----------+----------+
// | 1  |    1     | 1 to 255 |
// +----+----------+----------+
// VER: 协议版本，socks5为0x05
// NMETHODS: 支持认证的方法数量
// METHODS: 对应NMETHODS，NMETHODS的值为多少，METHODS就有多少个字节。RFC预定义了一些值的含义，内容如下:
// X’00’ NO AUTHENTICATION REQUIRED
// X’02’ USERNAME/PASSWORD
func Auth(reader *bufio.Reader, c net.Conn) (err error) {
	ver, err := reader.ReadByte()
	if err != nil {
		return fmt.Errorf("[auth]read version failed,err:%v", err)
	}
	if ver != types.Socks5Ver {
		return fmt.Errorf("[auth]type is not match,ver:%v", ver)
	}

	nMethods, err := reader.ReadByte()
	if err != nil {
		return fmt.Errorf("[auth]read nMethods failed,err:%v", err)
	}

	methods := make([]byte, nMethods)
	_, err = io.ReadFull(reader, methods)
	if err != nil {
		return fmt.Errorf("[auth]read methods failed,err:%v", err)
	}
	log.Printf("[auth]ver:%v,methods:%v\n", ver, methods)

	_, err = c.Write([]byte{types.Socks5Ver, 0x00})
	if err != nil {
		return fmt.Errorf("[auth]conn.write failed,err:%v", err)
	}
	return nil
}
