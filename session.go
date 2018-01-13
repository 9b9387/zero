package zero

import (
	uuid "github.com/satori/go.uuid"
)

type Session struct {
	sID      string
	uID      string
	conn     *Conn
	settings map[string]interface{}
}

func NewSession(conn *Conn) *Session {
	id, _ := uuid.NewV4()
	session := &Session{
		sID:      id.String(),
		uID:      "",
		conn:     conn,
		settings: make(map[string]interface{}),
	}

	return session
}

func (s *Session) GetSessionID() string {
	return s.sID
}

func (s *Session) BindUserID(uid string) {
	s.uID = uid
}

func (s *Session) GetUserID() string {
	return s.uID
}

func (s *Session) GetConn() *Conn {
	return s.conn
}

func (s *Session) SetConn(conn *Conn) {
	s.conn = conn
}

func (s *Session) GetSetting(key string) interface{} {

	if v, ok := s.settings[key]; ok {
		return v
	} else {
		return nil
	}
}

func (s *Session) SetSetting(key string, value interface{}) {
	s.settings[key] = value
}
