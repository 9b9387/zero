package zero

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// Encode from Message to []byte
func Encode(msg *Message) ([]byte, error) {
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.LittleEndian, msg.msgSize)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buffer, binary.LittleEndian, msg.msgID)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buffer, binary.LittleEndian, msg.data)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buffer, binary.LittleEndian, msg.checksum)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Decode from []byte to Message
func Decode(data []byte) (*Message, error) {
	bufReader := bytes.NewReader(data)

	dataSize := len(data)
	// 读取消息ID
	var msgID int32
	err := binary.Read(bufReader, binary.LittleEndian, &msgID)
	if err != nil {
		return nil, err
	}

	// 读取数据
	dataBufLength := dataSize - 4 - 4
	dataBuf := make([]byte, dataBufLength)
	err = binary.Read(bufReader, binary.LittleEndian, &dataBuf)
	if err != nil {
		return nil, err
	}

	// 检查checksum
	var checksum uint32
	err = binary.Read(bufReader, binary.LittleEndian, &checksum)
	if err != nil {
		return nil, err
	}

	message := &Message{}
	message.msgSize = int32(dataSize)
	message.msgID = msgID
	message.data = dataBuf
	message.checksum = checksum

	if message.Verify() {
		return message, nil
	}

	return nil, errors.New("checksum error")
}
