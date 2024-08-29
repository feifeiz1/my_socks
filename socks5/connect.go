package socks5

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/feifeiz1/my_socks/types"
	"io"
	"log"
	"net"
	"time"
)

// Connect 连接
// +----+-----+-------+------+----------+----------+
// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
// +----+-----+-------+------+----------+----------+
// | 1  |  1  | X'00' |  1   | Variable |    2     |
// +----+-----+-------+------+----------+----------+
// VER 版本号，socks5的值为0x05
// CMD 0x01表示CONNECT请求
// RSV 保留字段，值为0x00
// ATYP 目标地址类型，DST.ADDR的数据对应这个字段的类型。
//
//	0x01表示IPv4地址，DST.ADDR为4个字节
//	0x03表示域名，DST.ADDR是一个可变长度的域名
//
// DST.ADDR 一个可变长度的值
// DST.PORT 目标端口，固定2个字节
func Connect(reader *bufio.Reader, c net.Conn) (err error) {
	buf := make([]byte, 4)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return fmt.Errorf("[connect]read to buffer failed,err:%v", err)
	}
	ver, cmd, aType := buf[0], buf[1], buf[3]
	if ver != types.Socks5Ver {
		return fmt.Errorf("[connect]version not supposed,ver:%v", ver)
	}

	if cmd != types.CMDConnect {
		return fmt.Errorf("[connect]cmd is not supposed,cmd:%v", cmd)
	}

	addr := ""
	switch aType {
	case types.ATypeIPV4:
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return fmt.Errorf("[connect]read ipv4 addr failed,err:%v", err)
		}
		addr = fmt.Sprintf("%v.%v.%v.%v", buf[0], buf[1], buf[2], buf[3])
	case types.ATypeHost:
		hostSize, err := reader.ReadByte()
		if err != nil {
			return fmt.Errorf("[connect]read host size failed,err:%v", err)
		}
		host := make([]byte, hostSize)
		if _, err = io.ReadFull(reader, host); err != nil {
			return fmt.Errorf("[connect]read host failed,err:%v", err)
		}
		addr = string(host)
	case types.ATypeIPV6:
		return fmt.Errorf("[connect]unsupposed aType,aType:%v", aType)
	default:
		return fmt.Errorf("[connect]unknow aType,aType:%v", aType)
	}
	_, err = io.ReadFull(reader, buf[:2])
	if err != nil {
		return fmt.Errorf("[connect]read port failed,err:%v", err)
	}
	port := binary.BigEndian.Uint16(buf[:2])
	log.Println("addr:", addr, " port:", port)

	dest, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", addr, port), time.Second*5)
	if err != nil {
		return fmt.Errorf("[connect]dial dest failed,err:%v", err)
	}
	defer dest.Close()
	// +----+-----+-------+------+----------+----------+
	// |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	// VER socks版本，这里为0x05
	// REP Relay field,内容取值如下 X’00’ succeeded
	// RSV 保留字段
	// ATYPE 地址类型
	// BND.ADDR 服务绑定的地址
	// BND.PORT 服务绑定的端口DST.PORT

	_, err = c.Write([]byte{types.Socks5Ver, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		return fmt.Errorf("[connect]connection write failed,err:%v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_, cpErr := io.Copy(dest, reader)
		if cpErr != nil {
			log.Printf("[ERROR]transfer to dest failed,err:%v\n", cpErr)
		}
		cancel()
	}()

	go func() {
		_, cpErr := io.Copy(c, dest)
		if cpErr != nil {
			log.Printf("[ERROR]transfer to source failed,err:%v\n", cpErr)
		}
		cancel()
	}()
	<-ctx.Done()
	return nil
}
