package zero

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	host := "127.0.0.1:18787"

	ss, err := NewSocketService(host)
	if err != nil {
		return
	}

	// ss.SetHeartBeat(5*time.Second, 30*time.Second)

	ss.RegMessageHandler(HandleMessage)
	ss.RegConnectHandler(HandleConnect)
	ss.RegDisconnectHandler(HandleDisconnect)

	go NewClientConnect()

	timer := time.NewTimer(time.Second * 1)
	go func() {
		<-timer.C
		ss.Stop("stop service")
		t.Log("service stoped")
	}()

	t.Log("service running on " + host)
	ss.Serv()
}

func HandleMessage(s *Session, msg *Message) {
	fmt.Println("receive msgID:", msg)
	fmt.Println("receive data:", string(msg.GetData()))
}

func HandleDisconnect(s *Session, err error) {
	fmt.Println(s.GetConn().GetName() + " lost.")
}

func HandleConnect(s *Session) {
	fmt.Println(s.GetConn().GetName() + " connected.")
}

func NewClientConnect() {
	host := "127.0.0.1:18787"
	tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return
	}

	msg := NewMessage(1, []byte("Hello Zero!"))
	data, err := Encode(msg)
	if err != nil {
		return
	}
	conn.Write(data)
}
