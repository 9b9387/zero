package zero

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/adler32"
)

type Message struct {
	msgSize  int32
	msgID    int32
	data     []byte
	checksum uint32
}

func NewMessage(msgID int32, data []byte) *Message {
	msg := &Message{
		msgSize: int32(len(data)) + 4 + 4,
		msgID:   msgID,
		data:    data,
	}

	msg.checksum = msg.calcChecksum()
	return msg
}

func (msg *Message) GetData() []byte {
	return msg.data
}

func (msg *Message) GetID() int32 {
	return msg.msgID
}

func (msg *Message) Verify() bool {
	return msg.checksum == msg.calcChecksum()
}

func (msg *Message) calcChecksum() uint32 {
	if msg == nil {
		return 0
	}

	data := new(bytes.Buffer)

	err := binary.Write(data, binary.LittleEndian, msg.msgID)
	if err != nil {
		return 0
	}
	err = binary.Write(data, binary.LittleEndian, msg.data)
	if err != nil {
		return 0
	}

	checksum := adler32.Checksum(data.Bytes())
	return checksum
}

func (msg *Message) String() string {
	return fmt.Sprintf("Size=%d ID=%d DataLen=%d Checksum=%d", msg.msgSize, msg.GetID(), len(msg.GetData()), msg.checksum)
}
