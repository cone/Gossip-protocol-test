package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
	"time"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func main() {
	var clientID = flag.String("client", "Client1", "Client ID (string)")
	var local = flag.String("local", "10001", "Local port")
	var server = flag.String("server", "10002", "Server port")

	flag.Parse()

	//ServerAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+*server)
	//CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+*local)
	CheckError(err)

	fmt.Println("Ping to: " + *server)
	////////////server
	/* Now listen at selected port */
	go func() {
		Listener, err := net.ListenUDP("udp", LocalAddr)
		CheckError(err)
		defer Listener.Close()

		buf := make([]byte, 1024)

		for {
			n, addr, err := Listener.ReadFromUDP(buf)
			if err != nil {
				fmt.Println("Error: ", err)
			}

			strMsg := string(buf[0:n])

			strArr := strings.Split(strMsg, ":")

			fmt.Println("Received ", strArr[0], " from ", addr)

			if len(strArr) < 2 {
				continue
			}

			Conn, err := net.Dial("udp", "127.0.0.1:"+strArr[1])
			CheckError(err)

			_, err = Conn.Write([]byte("You sent me " + strArr[0]))
			CheckError(err)

			Conn.Close()
		}
	}()

	fmt.Println("Our port:" + *local)
	Conn, err := net.Dial("udp", "127.0.0.1:"+*server)
	//Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)

	defer Conn.Close()

	for {
		msg := *clientID + " says Hello:" + *local
		buf := []byte(msg)
		_, err := Conn.Write(buf)
		if err != nil {
			fmt.Println(msg, err)
		}
		time.Sleep(time.Second * 1)
	}

}
