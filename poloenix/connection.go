package poloenix

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type wsConn struct {
	connection *websocket.Conn
	mx         *sync.Mutex
}

func getConnection() (c *wsConn, err error) {
	conn, _, err := websocket.DefaultDialer.Dial(poloenixWsUrl, nil)
	if err != nil {
		err = fmt.Errorf("connecting to ws: %v", err)
		return
	}
	c = &wsConn{
		connection: conn,
		mx:         &sync.Mutex{},
	}

	return c, nil
}

func (w *wsConn) WriteJSON(value interface{}) error {
	w.mx.Lock()
	defer w.mx.Unlock()
	return w.connection.WriteJSON(value)
}
func (w *wsConn) ReadJSON(value interface{}) error {
	w.mx.Lock()
	defer w.mx.Unlock()
	return w.connection.ReadJSON(value)
}
