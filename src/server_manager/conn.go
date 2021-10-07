
package kweb_server_mnr

import (
	time "time"
	sync "sync"

	uuid "github.com/google/uuid"
	websocket "github.com/kataras/iris/v12/websocket"
)


type ServerConn struct{
	id string
	wsconn *websocket.Conn
	Data *ServerData
	serves map[string]bool
	msgcache *[]string
	Storage map[string]interface{}

	msglock sync.Mutex

	currentspace string
}

func NewServerConn(data *ServerData, wsconn *websocket.Conn)(*ServerConn){
	return &ServerConn{
		id: uuid.NewString(),
		wsconn: wsconn,
		Data: data,
		serves: map[string]bool{"": true},
		msgcache: &[]string{},
		Storage: make(map[string]interface{}),
	}
}

func (conn *ServerConn)Clone()(*ServerConn){
	return &ServerConn{
		id: conn.id,
		wsconn: conn.wsconn,
		Data: conn.Data,
		serves: conn.serves,
		msgcache: conn.msgcache,
		Storage: conn.Storage,
		msglock: conn.msglock,
		currentspace: conn.currentspace,
	}
}

func (conn *ServerConn)ID()(string){
	return conn.id
}

func (conn *ServerConn)Init(){
	if _, ok := conn.Storage[_BUILDIN_ISINIT_ID]; ok {
		return
	}
	conn.Storage[_BUILDIN_ISINIT_ID] = struct{}{}
	exitch, pingch := make(chan struct{}, 1), make(chan struct{}, 2)
	conn.Storage[_BUILDIN_EXIT_CH_ID], conn.Storage[_BUILDIN_PING_CH_ID] = (chan<- struct{})(exitch), (chan<- struct{})(pingch)
	setServerConnect(conn.Data.Group, conn)
	go func(){
		for{
			select{
			case <-exitch:
				return
			case <-pingch:
			case <-time.After(60 * time.Second):
				conn.Close()
				return
			}
		}
	}()
}

func (conn *ServerConn)Close(){
	if _, ok := conn.Storage[_BUILDIN_ISCLOSE_ID]; ok {
		return
	}
	if conn.wsconn != nil {
		conn.wsconn.Close()
	}
	delServerConnect(conn.Data.Group)
	conn.Storage[_BUILDIN_ISCLOSE_ID] = struct{}{}
	conn.Storage[_BUILDIN_EXIT_CH_ID].(chan<- struct{}) <- struct{}{}
}

func (conn *ServerConn)HasServe(space string)(ok bool){
	x, ok := conn.serves[space]
	if !(ok && x) {
		return false
	}
	_, ok = SERVE_MAP[space]
	return
}

func (conn *ServerConn)GetServe(space string)(s Serve){
	x, ok := conn.serves[space]
	if !(ok && x) {
		return nil
	}
	s, ok = SERVE_MAP[space]
	if ok {
		return s
	}
	return nil
}

func (conn *ServerConn)SendMessage(id string, data string)(err error){
	return conn.SendMessageEvent(&Event{conn.currentspace, id, data})
}

func (conn *ServerConn)SendMessageEvent(event *Event)(err error){
	conn.msglock.Lock()
	defer conn.msglock.Unlock()
	if conn.wsconn != nil {
		err = conn.wsconn.Socket().WriteText(event.Bytes(), time.Second * 3)
		return
	}
	*conn.msgcache = append(*conn.msgcache, event.String())
	return nil
}

func (conn *ServerConn)OnMessage(msg string)(re *Event){
	event := ParseEvent(msg)
	if event == nil {
		return nil
	}
	if serve := conn.GetServe(event.Space); serve != nil {
		conn.currentspace = event.Space
		re = serve.OnMessage(conn, event)
		conn.currentspace = ""
		return
	}
	return nil
}
