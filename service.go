package zero

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"
)

type SocketService struct {
	onMessage    func(*Session, *Message)
	onConnect    func(*Session)
	onDisconnect func(*Session, error)
	sessions     *sync.Map
	hbInterval   time.Duration
	hbTimeout    time.Duration
	laddr        string
	status       int
	listener     net.Listener
	stopCh       chan error
}

func NewSocketService(laddr string) (*SocketService, error) {

	l, err := net.Listen("tcp", laddr)

	if err != nil {
		return nil, err
	}

	s := &SocketService{
		sessions:   &sync.Map{},
		stopCh:     make(chan error),
		hbInterval: 0 * time.Second,
		hbTimeout:  0 * time.Second,
		laddr:      laddr,
		status:     SERVER_ST_INITED,
		listener:   l,
	}

	return s, nil
}

func (s *SocketService) RegOnMessageHandler(handler func(*Session, *Message)) {
	s.onMessage = handler
}

func (s *SocketService) RegOnConnectHandler(handler func(*Session)) {
	s.onConnect = handler
}

func (s *SocketService) RegOnDisconnectHandler(handler func(*Session, error)) {
	s.onDisconnect = handler
}

func (s *SocketService) Serv() {

	s.status = SERVER_ST_RUNNING
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		s.status = SERVER_ST_STOP
		cancel()
		s.listener.Close()
	}()

	go s.acceptHandler(ctx)

	for {
		select {

		case <-s.stopCh:
			return
		}
	}
}

func (s *SocketService) acceptHandler(ctx context.Context) {
	for {
		c, err := s.listener.Accept()
		if err != nil {
			s.stopCh <- err
			return
		}

		go s.connectHandler(ctx, c)
	}
}

func (s *SocketService) connectHandler(ctx context.Context, c net.Conn) {
	conn := NewConn(c, s.hbInterval, s.hbTimeout)
	session := NewSession(conn)
	s.sessions.Store(session.GetSessionID(), session)

	connctx, cancel := context.WithCancel(ctx)

	defer func() {
		cancel()
		conn.Close()
		s.sessions.Delete(session.GetSessionID())
	}()

	go conn.readCoroutine(connctx)
	go conn.writeCoroutine(connctx)

	if s.onConnect != nil {
		s.onConnect(session)
	}

	for {
		select {
		case err := <-conn.done:

			if s.onDisconnect != nil {
				s.onDisconnect(session, err)
			}
			return

		case msg := <-conn.messageCh:
			if s.onMessage != nil {
				s.onMessage(session, msg)
			}
		}
	}
}

func (s *SocketService) GetStatus() int {
	return s.status
}

func (s *SocketService) Stop(reason string) {
	s.stopCh <- errors.New(reason)
}

func (s *SocketService) SetHeartBeat(hbInterval time.Duration, hbTimeout time.Duration) error {
	if s.status == SERVER_ST_RUNNING {
		return errors.New("Can't set heart beat on service running.")
	}

	s.hbInterval = hbInterval
	s.hbTimeout = hbTimeout

	return nil
}

func (s *SocketService) GetConnsCount() int {
	var count int
	s.sessions.Range(func(k, v interface{}) bool {
		count++
		return true
	})
	return count
}

// 发送消息
func (s *SocketService) Unicast(sid string, msg *Message) {
	v, ok := s.sessions.Load(sid)
	if ok {
		session := v.(*Session)
		err := session.GetConn().SendMessage(msg)
		if err != nil {
			return
		}
	}
}

// 广播消息
func (s *SocketService) Broadcast(msg *Message) {
	s.sessions.Range(func(k, v interface{}) bool {
		s := v.(*Session)
		if err := s.GetConn().SendMessage(msg); err != nil {
			// log.Println(err)
		}
		return true
	})
}
