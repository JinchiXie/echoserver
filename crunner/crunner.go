package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"

	"github.com/JinchiXie/echoserver/server"
)

const (
	defaultHost = "localhost"
	defaultPort = 9999
)

// To test your server implementation, you might find it helpful to implement a
// simple 'client runner' program. The program could be very simple, as long as
// it is able to connect with and send messages to your server and is able to
// read and print out the server's echoed response to standard output. Whether or
// not you add any code to this file will not affect your grade.
func main() {
	s := server.New()
	s.Start(defaultPort)
	conn, _ := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(defaultPort))
	send := "test fggg fgfg\n"
	conn.Write([]byte(send))
	mes, error := bufio.NewReader(conn).ReadBytes('\n')
	fmt.Println("from server: " + string(mes))
	//fmt.Println("send: " + send)
	//fmt.Println("same : %t", string(mes) == send)
	if error != nil {
		fmt.Println(error.Error())
	}
	fmt.Println("before connection closed")
	s.Close()
}
