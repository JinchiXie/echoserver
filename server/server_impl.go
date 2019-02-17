// Implementation of a MultiEchoServer. Students should write their code in this file.

package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

type set map[net.Conn]bool
type multiEchoServer struct {
	// TODO: implement this!
	listener net.Listener
	counter  int
	close    chan int
	quit     chan int
	cons     set
	mux      sync.Mutex
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	// TODO: implement this!
	server := &multiEchoServer{
		listener: nil,
		counter:  0,
		close:    make(chan int, 1),
		quit:     make(chan int, 1),
		cons:     make(set),
	}
	return server
}

func (mes *multiEchoServer) Start(port int) error {
	// TODO: implement this!

	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err == nil {
		mes.listener = l
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
			if err == nil {
				mes.atomInc()
				mes.cons[conn] = true
				go mes.handle(conn)
			} else {
				//fmt.Println("meet error")
				return
			}
		}

	}
}

func (mes *multiEchoServer) handle(conn net.Conn) {
	//
	defer conn.Close()
	//defer fmt.Println(mes.atomLoad())
	defer mes.atomDec()
	//fmt.Println(mes.atomLoad())

	//fmt.Println(mes.atomLoad())
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	//for {
	select {
	case <-mes.close:
		//fmt.Println("accept close")
		return
	default:
		//fmt.Println("before read")
		message, err := bufio.NewReader(conn).ReadBytes('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			//fmt.Println("server error: " + err.Error())
		}
		conn.Write([]byte(string(message)))

	}
	//}

}

func (mes *multiEchoServer) Close() {
	// TODO: implement this!
	fmt.Println("send semo to mes.quit")
	mes.quit <- 1
	//fmt.Println("send semo to mes.quit done")
	mes.listener.Close()
	for mes.atomLoad() != 0 {
		select {
		case mes.close <- 1:
			//fmt.Println("send close semo")
		default:
			//fmt.Println("wait")

		}
		//fmt.Println(mes.atomLoad())
	}
	//fmt.Println("all close")
}

func (mes *multiEchoServer) Count() int {
	// TODO: implement this!
	return mes.atomLoad()
}

// TODO: add additional methods/functions below!
func (mes *multiEchoServer) atomInc() {
	defer mes.mux.Unlock()
	mes.mux.Lock()
	mes.counter++

}

func (mes *multiEchoServer) atomLoad() int {
	defer mes.mux.Unlock()
	mes.mux.Lock()
	return mes.counter
}

func (mes *multiEchoServer) atomDec() {
	defer mes.mux.Unlock()
	mes.mux.Lock()
	mes.counter--
}
