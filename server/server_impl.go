// Implementation of a MultiEchoServer. Students should write their code in this file.

package server

import (
	"fmt"
	"net"
	"strconv"
)

type multiEchoServer struct {
	// TODO: implement this!
	listener net.Listener
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	// TODO: implement this!
	server := &multiEchoServer{nil}
	return server
}

func (mes *multiEchoServer) Start(port int) error {
	// TODO: implement this!

	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err == nil {
		mes.listener = l
	} else {
		fmt.Println("can't greate server error: " + err.Error())
	}
	return err
}

func accept(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		ch := make(chan [1]bool)
		if err == nil {
			go dealMsg(conn, ch)
		} else {
			break
		}
	}
}

func dealMsg(conn net.Conn, stop chan [1]bool) {
	//
}

func (mes *multiEchoServer) Close() {
	// TODO: implement this!
	mes.listener.Close()
}

func (mes *multiEchoServer) Count() int {
	// TODO: implement this!
	return -1
}

// TODO: add additional methods/functions below!
