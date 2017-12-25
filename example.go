package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"./zero"
)

func main() {
	host := "127.0.0.1:18787"

	ss, err := zero.NewSocketService(host)
	if err != nil {
		return
	}

	ss.SetHeartBeat(5*time.Second, 30*time.Second)

	ss.RegOnMessageHandler(HandleMessage)
	ss.RegOnConnectHandler(HandleConnect)
	ss.RegOnDisconnectHandler(HandleDisconnect)

	log.Println("server running on " + host)

	go NewClientConnect()
	ss.Serv()
}

func HandleMessage(s *zero.Session, msg *zero.Message) {
	var msgID int32 = msg.GetID()
	log.Println("receive msgID=", msgID)
}

func HandleDisconnect(s *zero.Session, err error) {
	log.Println(s.GetConn().GetName() + " lost.")
}

func HandleConnect(s *zero.Session) {
	log.Println(s.GetConn().GetName() + " connected.")
}

func NewClientConnect() {
	host := "127.0.0.1:18787"
	tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return
	}

	_, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return
	}
}
