
package kweb_server_mnr

import (
	sync "sync"

	uuid "github.com/google/uuid"
)

var _SERVER_CONNECTS = make(map[uuid.UUID]*ServerConn)
var _SERVER_CONNMAP_LOCK = sync.RWMutex{}

func HasServerConnect(id uuid.UUID)(ok bool){
	_SERVER_CONNMAP_LOCK.RLock()
	defer _SERVER_CONNMAP_LOCK.RUnlock()
	_, ok = _SERVER_CONNECTS[id]
	return
}

func GetServerConnect(id uuid.UUID)(*ServerConn){
	_SERVER_CONNMAP_LOCK.RLock()
	defer _SERVER_CONNMAP_LOCK.RUnlock()
	if conn, ok := _SERVER_CONNECTS[id]; ok {
		return conn
	}
	return nil
}

func setServerConnect(id uuid.UUID, conn *ServerConn){
	if HasServerConnect(id) {
		panic("'" + conn.Data.Id + "' is already connected")
	}
	_SERVER_CONNMAP_LOCK.Lock()
	defer _SERVER_CONNMAP_LOCK.Unlock()
	_SERVER_CONNECTS[id] = conn
}

func delServerConnect(id uuid.UUID){
	if HasServerConnect(id) {
		_SERVER_CONNMAP_LOCK.Lock()
		defer _SERVER_CONNMAP_LOCK.Unlock()
		delete(_SERVER_CONNECTS, id)
	}
}
