package zero

type MessageHolder struct {
	uuid    string
	message *Message
}

const (
	SERVER_ST_UNKNOW = iota
	SERVER_ST_INITED
	SERVER_ST_RUNNING
	SERVER_ST_STOP
)

const (
	MSGID_HEARTBEAT = iota
)
