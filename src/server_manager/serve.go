
package kweb_server_mnr

import (
	time "time"
	fmt "fmt"
	strings "strings"

	json "github.com/KpnmServer/go-util/json"
)

type Event struct{
	Space string
	Id string
	Data string
}

func ParseEvent(msg string)(*Event){
	msgs := strings.SplitN(msg, ";", 3)
	if len(msgs) < 3 {
		return nil
	}
	return &Event{
		Space: msgs[0],
		Id: msgs[1],
		Data: msgs[2],
	}
}

func (e *Event)String()(string){
	if e == nil {
		return ""
	}
	return e.Space + ";" + e.Id + ";" + e.Data
}

func (e *Event)Bytes()([]byte){
	if e == nil {
		return nil
	}
	return ([]byte)(e.String())
}

type Serve interface{
	OnRegister(conn *ServerConn)
	OnUnRegister(conn *ServerConn)
	OnMessage(conn *ServerConn, event *Event)(re *Event)
}

var SERVE_MAP = make(map[string]Serve)

func RegisterServe(space string, serve Serve){
	if _, ok := SERVE_MAP[space]; ok {
		panic("serve '" + space + "' already exists")
	}
	SERVE_MAP[space] = serve
}

type builtInServe struct{}

const (
	_BUILDIN_ISINIT_ID  = "buildin.isinit"
	_BUILDIN_ISCLOSE_ID = "buildin.isclose"
	_BUILDIN_EXIT_CH_ID = "buildin.exitch"
	_BUILDIN_PING_CH_ID = "buildin.pingch"
)

func (builtInServe)OnRegister(conn *ServerConn){}

func (builtInServe)OnUnRegister(conn *ServerConn){}

func (builtInServe)OnMessage(conn *ServerConn, event *Event)(re *Event){
	data := event.Data
	switch event.Id {
	case "init":
		conn.Init()
		return nil
	case "close":
		conn.Close()
	case "ping":
		conn.Storage[_BUILDIN_PING_CH_ID].(chan<- struct{}) <- struct{}{}
		conn.msglock.Lock()
		msglist := *conn.msgcache
		conn.msgcache = &[]string{}
		conn.msglock.Unlock()
		return &Event{
			Space: event.Space,
			Id: "pong",
			Data: fmt.Sprintf("%d;%s", time.Now().Unix(), json.EncodeJsonStr(msglist)),
		}
	case "register":
		if data != "" {
			if x ,ok := conn.serves[data]; ok && x {
				return nil
			}
			if s, ok := SERVE_MAP[data]; ok {
				conn.serves[data] = true
				s.OnRegister(conn)
			}
		}
	case "unregister":
		if data != "" {
			if x ,ok := conn.serves[data]; ok && x {
				if s, ok := SERVE_MAP[data]; ok {
					conn.serves[data] = false
					s.OnUnRegister(conn)
				}
			}
		}
	}
	return nil
}


func init(){
	RegisterServe("", builtInServe{})
}