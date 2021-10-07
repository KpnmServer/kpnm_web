
package kweb_server_mnr


type EventHandler func(conn *ServerConn, event *Event)(re *Event)
type EventMap map[string]EventHandler

const (
	EVENT_ON_REGISTER = "!register"
	EVENT_ON_UNREGISTER = "!unregister"
	EVENT_DEFAULT = "!default"
)

type EventServe struct{
	emap EventMap
}

func NewEventServe(events EventMap)(*EventServe){
	return &EventServe{
		emap: events,
	}
}

func RegisterEventServe(space string, events EventMap){
	RegisterServe(space, NewEventServe(events))
}

func (es *EventServe)OnRegister(conn *ServerConn){
	es.OnMessage(conn, &Event{"", EVENT_ON_REGISTER, ""})
}

func (es *EventServe)OnUnRegister(conn *ServerConn){
	es.OnMessage(conn, &Event{"", EVENT_ON_UNREGISTER, ""})
}

func (es *EventServe)OnMessage(conn *ServerConn, event *Event)(re *Event){
	handler, ok := es.emap[event.Id]
	if !ok {
		if handler, ok = es.emap[EVENT_DEFAULT]; !ok{
			return nil
		}
	}
	return handler(conn, event)
}