package zero

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"net"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Conn struct {
	uuid       string
	rawConn    net.Conn
	sendCh     chan []byte
	done       chan error
	hbTimer    *time.Timer
	name       string
	messageCh  chan MessageHolder
	hbInterval time.Duration
	hbTimeout  time.Duration
}

func (c *Conn) GetName() string {
	return c.name
}

func (c *Conn) GetUUID() string {
	return c.uuid
}

func NewConn(c net.Conn, hbInterval time.Duration, hbTimeout time.Duration) *Conn {
	conn := &Conn{
		rawConn:    c,
		sendCh:     make(chan []byte, 100),
		done:       make(chan error),
		messageCh:  make(chan MessageHolder, 100),
		hbInterval: hbInterval,
		hbTimeout:  hbTimeout,
	}

	conn.name = c.RemoteAddr().String()
	conn.uuid = uuid.NewV4().String()
	conn.hbTimer = time.NewTimer(conn.hbInterval)

	if conn.hbInterval == 0 {
		conn.hbTimer.Stop()
	}

	return conn
}

func (c *Conn) Close() {
	c.hbTimer.Stop()
	c.rawConn.Close()
}

func (c *Conn) SendMessage(msg *Message) error {
	pkg, err := Encode(msg)
	if err != nil {
		return err
	}

	c.sendCh <- pkg
	return nil
}

func (c *Conn) writeCoroutine(ctx context.Context) {
	hbData := make([]byte, 0)

	for {
		select {
		case <-ctx.Done():
			return

		case pkt := <-c.sendCh:

			if pkt == nil {
				continue
			}

			if _, err := c.rawConn.Write(pkt); err != nil {
				c.done <- err
			}

		case <-c.hbTimer.C:
			hbMessage := NewMessage(MSGID_HEARTBEAT, hbData)
			c.SendMessage(hbMessage)
		}
	}
}

func (c *Conn) readCoroutine(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return

		default:
			// 设置超时
			if c.hbInterval > 0 {
				err := c.rawConn.SetReadDeadline(time.Now().Add(c.hbTimeout))
				if err != nil {
					c.done <- err
					continue
				}
			}
			// 读取长度
			buf := make([]byte, 4)
			_, err := io.ReadFull(c.rawConn, buf)
			if err != nil {
				c.done <- err
				continue
			}

			bufReader := bytes.NewReader(buf)

			var dataSize int32
			err = binary.Read(bufReader, binary.LittleEndian, &dataSize)
			if err != nil {
				c.done <- err
				continue
			}

			// 读取数据
			databuf := make([]byte, dataSize)
			_, err = io.ReadFull(c.rawConn, databuf)
			if err != nil {
				c.done <- err
				continue
			}

			// 解码
			msg, err := Decode(databuf)
			if err != nil {
				c.done <- err
				continue
			}

			// 设置心跳timer
			if c.hbInterval > 0 {
				c.hbTimer.Reset(c.hbInterval)
			}

			if msg.GetID() == MSGID_HEARTBEAT {
				continue
			}

			c.messageCh <- MessageHolder{
				message: msg,
				uuid:    c.uuid,
			}
		}
	}
}
