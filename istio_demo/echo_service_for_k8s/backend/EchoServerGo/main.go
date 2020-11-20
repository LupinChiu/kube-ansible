package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	port := os.Args[1]
	host := fmt.Sprintf("0.0.0.0:%v", port)
	listen, err := net.Listen("tcp", host)
	if err != nil {
		fmt.Println(err)
	}
	defer listen.Close()
	fmt.Printf("TCP server start and listening on %s.\n", host)

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("Some connection error: %s\n", err)
		}

		go handleConnection(conn)
	}

	wait := make(chan bool)
	<-wait
}

func handleConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Client connected from: " + remoteAddr)

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	for {
		// Read the incoming connection into the buffer.
		reqLen, err := conn.Read(buf)
		if err != nil {

			if err.Error() == "EOF" {
				fmt.Println("Disconned from ", remoteAddr)
				break
			} else {
				fmt.Println("Error reading:", err.Error())
				break
			}
		}
		// Send a response back to person contacting us.
		conn.Write([]byte("server send:" + time.Now().String()))

		fmt.Printf("len: %d, recv: %s\n", reqLen, string(buf[:reqLen]))
	}
	// Close the connection when you're done with it.
	conn.Close()
}
