package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"paxos/log"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Log("connect to server failed, err[%v]", err)
		return
	}

	stdReader := bufio.NewReader(os.Stdin)
	tcpReader := bufio.NewReader(conn)
	for {
		bytes, err := stdReader.ReadBytes('\n')
		if err != nil {
			log.Log("read line from stdin failed, err[%v]", err)
			break
		}

		if err = conn.Close(); err != nil {
			log.Log("conn close failed, err[%v]")
			break
		}

		if _, err = conn.Write(bytes); err != nil {
			log.Log("write bytes to conn failed, err[%v]", err)
			break
		}

		back, err := tcpReader.ReadString('\n')
		if err != nil {
			log.Log("read line from tcp conn failed, err[%v]", err)
			break
		}

		fmt.Print(back)
	}
}
