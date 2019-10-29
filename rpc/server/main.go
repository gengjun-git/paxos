package main

import (
	"bufio"
	"net"
	"paxos/log"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Log("listen tcp on port 8080 failed, err[%v]", err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Log("tcp accept failed, err[%v]", err)
			break
		}

		go echo(conn)
	}

	log.Log("server exit")
}

func echo(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		bytes, err := reader.ReadBytes('\n')
		if err != nil {
			log.Log("read data from conn failed, err[%v]", err)
			break
		}
		log.Log("read data: %s", bytes)

		if _, err = conn.Write(bytes); err != nil {
			log.Log("write data to conn failed, err[%v]", err)
		}
	}
}
