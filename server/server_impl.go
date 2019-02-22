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
	mes      *multiEchoServer
	curWrite int
}

type multiEchoServer struct {
	// TODO: implement this!
	listener     net.Listener
	clients      chan map[int]*client
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
		clients:      make(chan map[int]*client, 1),
	}
	server.clients <- make(map[int]*client) //init clients
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
				//continue
				return
			}
			cli := &client{
				conn:     conn,
				id:       mes.curID,
				sendQuit: make(chan bool, 1),
				recQuit:  make(chan bool, 1),
				//readMsg:  make(chan []byte, 1),
				writeMsg: make(chan []byte, 30),
				mes:      mes,
				curWrite: 0,
			}
			clients := <-mes.clients
			clients[mes.curID] = cli
			mes.clients <- clients
			mes.curID++
			go cli.loopRead()
			go cli.loopWrite()

		}

	}
}

func (mes *multiEchoServer) bordcast() {
	for {
		data := <-mes.boradcastMsg
		clients := <-mes.clients
		for _, cli := range clients {
			// if cli.curWrite < 100 {
			// 	cli.curWrite++
			// 	cli.writeMsg <- data
			// }
			//cli.writeMsg <- data
			select {
			case cli.writeMsg <- data:
				//break
			default:
				//close(cli)
			}
		}
		mes.clients <- clients
	}
}

func (cli *client) loopRead() {
	reader := bufio.NewReader(cli.conn)
	for {
		msg, err := reader.ReadBytes('\n')
		if err != nil {
			cli.sendQuit <- true
			return
			//continue
		}
		cli.mes.boradcastMsg <- msg
	}
}

func (cli *client) loopWrite() {
	for {
		select {
		case msg := <-cli.writeMsg:
			//cli.conn.SetWriteDeadline(time.Now().Add(15000))
			_, err := cli.conn.Write(msg)
			//cli.curWrite--
			if err != nil {
				continue
			}
		case <-cli.sendQuit:
			deleteClient(cli)
			close(cli)
			return
		}
	}
}

func (mes *multiEchoServer) Close() {
	// TODO: implement this!
	mes.quit <- true
	mes.listener.Close()
	clients := <-mes.clients
	for _, client := range clients {
		client := client
		close(client)
	}
	mes.clients <- clients
}

func close(cli *client) {
	cli.conn.Close()
	//fmt.Println("connectin closed" + strconv.Itoa(cli.id))
}

func deleteClient(cli *client) {
	clients := <-cli.mes.clients
	if _, ok := clients[cli.id]; ok == true {
		delete(clients, cli.id)
	}
	cli.mes.clients <- clients
}

func (mes *multiEchoServer) Count() int {
	// TODO: implement this!
	clients := <-mes.clients
	defer func() {
		mes.clients <- clients
	}()
	return len(clients)
}

// TODO: add additional methods/functions below!
