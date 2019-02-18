// Implementation of a MultiEchoServer. Students should write their code in this file.

package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

type client struct {
	conn     net.Conn
	id       int
	sendQuit chan bool
	recQuit  chan bool
	readMsg  chan []byte
	writeMsg chan []byte
}

type multiEchoServer struct {
	// TODO: implement this!
	listener     net.Listener
	clients      map[int]*client
	quit         chan bool
	curID        int
	boradcastMsg chan []byte
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	// TODO: implement this!
	server := &multiEchoServer{
		curID:        0,
		boradcastMsg: make(chan []byte, 1),
		quit:         make(chan bool, 1),
		clients:      make(map[int]*client),
	}
	return server
}

func (mes *multiEchoServer) Start(port int) error {
	// TODO: implement this!

	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err == nil {
		mes.listener = l
		go mes.bordcast()
		go mes.serve()
	} else {
		fmt.Println("can't greate server error: " + err.Error())
	}
	return err
}

func (mes *multiEchoServer) serve() {
	for {

		select {
		case <-mes.quit:
			//fmt.Println("get quit semo")
			return
		default:
			//mes.listener.SetDeadline(time.Now().Add(1e9))
			conn, err := mes.listener.Accept()
			if err != nil {
				continue
			}
			cli := &client{
				conn:     conn,
				id:       mes.curID,
				sendQuit: make(chan bool, 1),
				recQuit:  make(chan bool, 1),
				readMsg:  make(chan []byte, 1),
				writeMsg: make(chan []byte, 1),
			}
			mes.clients[mes.curID] = cli
			mes.curID++
			go cli.loopRead(mes.boradcastMsg)
			go cli.loopWrite()

		}

	}
}

func (mes *multiEchoServer) bordcast() {
	for {
		data := <-mes.boradcastMsg
		clients := mes.clients
		for _, cli := range clients {
			cli.writeMsg <- data
		}
	}
}

func (cli *client) loopRead(boradcastMsg chan []byte) {
	for {
		select {
		case <-cli.recQuit:
			cli.sendQuit <- true
			return
		default:
			msg, err := bufio.NewReader(cli.conn).ReadBytes('\n')
			if err != nil {
				cli.sendQuit <- true
				return
			}
			boradcastMsg <- msg
		}
	}
}

func (cli *client) loopWrite() {
	for {
		select {
		case msg := <-cli.writeMsg:
			cli.conn.Write(msg)
		case <-cli.sendQuit:
			cli.conn.Close()
			return

		}
	}
}

func (mes *multiEchoServer) Close() {
	// TODO: implement this!
	mes.quit <- true
	mes.listener.Close()
	clients := mes.clients
	for id := range clients {
		mes.close(clients[id])
	}
}

func (mes *multiEchoServer) close(cli *client) {
	cli.conn.Close()
	cli.recQuit <- true
	delete(mes.clients, cli.id)
}

func (mes *multiEchoServer) Count() int {
	// TODO: implement this!
	return len(mes.clients)
}

// TODO: add additional methods/functions below!
