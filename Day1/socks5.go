package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

//版本方法
const socks5Ver = 0x05
const cmdBind = 0x01
const atypIPV4 = 0x01
const atypeHOST = 0x03
const atypeIPV6 = 0x04

func auth(reader *bufio.Reader, conn net.Conn) (err error) {
	//从io流中读一个字节
	ver, err := reader.ReadByte()
	if err != nil {
		return fmt.Errorf("read ver failed:%w", err)
	}
	//协议版本号 固定是5
	if ver != socks5Ver {
		return fmt.Errorf("not supported ver:%v", ver)
	}
	//再读下一个字节 拿到method size
	methodSize, err := reader.ReadByte()
	if err != nil {
		return fmt.Errorf("read methodSize failed:%w", err)
	}
	//make 对应大小的slice
	method := make([]byte, methodSize)
	//填充进去
	_, err = io.ReadFull(reader, method)
	if err != nil {
		return fmt.Errorf("read method failed:%w", err)
	}
	log.Println("ver", ver, "method", method)
	//version 和 method 作为响应写到conn中
	_, err = conn.Write([]byte{socks5Ver, 0x00})
	if err != nil {
		return fmt.Errorf("write failed:%w", err)
	}
	return nil
}
func connect(reader *bufio.Reader, conn net.Conn) (err error) {
	//包的长度一共四个字节 0 1 3分别是 version cmd atype
	buf := make([]byte, 4)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return fmt.Errorf("read header failed:%w", err)
	}
	ver, cmd, atyp := buf[0], buf[1], buf[3]
	//判断ver是不是socks5
	if ver != socks5Ver {
		return fmt.Errorf("not supported ver:%v", ver)
	}
	//判断cmd是不是1
	if cmd != cmdBind {
		return fmt.Errorf("not supported cmd:%v", ver)
	}
	//填充地址
	addr := ""
	switch atyp {
	//IPV4类型 直接填满
	case atypIPV4:
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return fmt.Errorf("read atyp failed:%w", err)
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
	//HOST类型
	case atypeHOST:
		//拿到hostSize之后 Readfull填充
		hostSize, err := reader.ReadByte()
		if err != nil {
			return fmt.Errorf("read hostSize failed:%w", err)
		}
		host := make([]byte, hostSize)
		_, err = io.ReadFull(reader, host)
		if err != nil {
			return fmt.Errorf("read host failed,%w", err)
		}
		addr = string(host)
	case atypeIPV6:
		return errors.New("IPv6 not supported yet")
	default:
		return errors.New("invalid atyp")
	}
	_, err = io.ReadFull(reader, buf[:2])
	if err != nil {
		return fmt.Errorf("read port failed:%w", err)
	}
	//从协议规定的大端字节序读取到port（2 byte)
	port := binary.BigEndian.Uint16(buf[:2])
	//dial 该服务器
	dest, err := net.Dial("tcp", fmt.Sprintf("%v:%v", addr, port))
	if err != nil {
		return fmt.Errorf("dial dst failed:%w", err)
	}
	defer dest.Close()
	log.Printf("dial", addr, port)
	//向conn写东西 后面的0都是填充字段
	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		return fmt.Errorf("write failed:%w", err)
	}
	//建立一个有cancel的context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_, _ = io.Copy(dest, reader)
		cancel()
	}()
	go func() {
		_, _ = io.Copy(conn, dest)
		cancel()
	}()
	//一旦cancel执行 立马done
	<-ctx.Done()
	return nil
}
func process(conn net.Conn) {
	//延迟关闭
	defer conn.Close()
	//读conn中的数据
	reader := bufio.NewReader(conn)
	//调用auth 和 connect函数
	err := auth(reader, conn)
	if err != nil {
		log.Printf("client %v auth failed:%v", conn.RemoteAddr(), err)
		return
	}
	err = connect(reader, conn)
	if err != nil {
		if err != nil {
			log.Printf("client %v connect failed:%v", conn.RemoteAddr(), err)
			return
		}
	}
}
func main() {
	//监听8080端口 这里一定要写全IP地址 不然可能会出现半连接的问题
	server, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	//开循环不停接受请求
	for {
		client, err := server.Accept()
		if err != nil {
			log.Printf("Accept failed %v", err)
			continue
		}
		go process(client)
	}
}
