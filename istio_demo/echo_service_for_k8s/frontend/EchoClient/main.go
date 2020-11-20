package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	host := os.Args[1]
	port := os.Args[2]
	addr := fmt.Sprintf("%s:%s", host, port)
	res, err := startTCP(addr)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}

	wait := make(chan bool)
	<-wait
}

func startTCP(addr string) (string, error) {
	fmt.Printf("1.1.0, Connect to: %v\n",addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	for {
		conn.Write([]byte(time.Now().String()))
		time.Sleep(5 * time.Second)

		bs := make([]byte, 1024)
		len, err := conn.Read(bs)
		if err != nil {
			fmt.Println(err)

		} else {
			fmt.Println(string(bs[:len]))
		}
	}

	return "Heartbeat Start", nil
}
