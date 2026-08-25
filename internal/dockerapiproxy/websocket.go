package dockerapiproxy

import (
	"encoding/base64"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type WebSocketIo struct {
	Conn *websocket.Conn
}

func (w *WebSocketIo) Read() ([]byte, error) {
	_, bytes, err := w.Conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	bytes, err = base64.StdEncoding.DecodeString(string(bytes))
	logrus.Debugf("Websocket read: %d bytes", len(bytes))
	return bytes, err
}

func (w *WebSocketIo) Write(buf []byte) (int, error) {
	logrus.Debugf("Websocket write: %d bytes", len(buf))
	str := base64.StdEncoding.EncodeToString(buf)
	err := w.Conn.WriteMessage(websocket.TextMessage, []byte(str))
	return len(buf), err
}
