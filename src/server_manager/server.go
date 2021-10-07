
package kweb_server_mnr

import (
	crand "crypto/rand"
	base64 "encoding/base64"

	uuid "github.com/google/uuid"
	kpsql "github.com/KpnmServer/go-kpsql"
	sql_mnr "github.com/KpnmServer/kpnm_web/src/sql"
	data_mnr "github.com/KpnmServer/kpnm_web/src/data_manager"
)

var SERVER_DATA_FOLDER = data_mnr.GetDataFolder("server")

type ServerStatus uint8

const (
	SERVER_OFFLINE ServerStatus = iota
	SERVER_CONN_HTTP
	SERVER_CONN_WS
)

func (s ServerStatus)String()(string){
	switch s {
	case SERVER_OFFLINE:
		return "OFFLINE"
	case SERVER_CONN_HTTP:
		return "CONN_HTTP"
	case SERVER_CONN_WS:
		return "CONN_WS"
	}
	return "UNKNOWN"
}

type ServerData struct{
	Id string            `sql:"id" sql_primary:"true"`
	Name string          `sql:"name"`
	Version string       `sql:"version"`
	Description string   `sql:"description"`
	Addrstr string       `sql:"addrstr"`
	Group uuid.UUID      `sql:"group"`
	State ServerStatus   `sql:"status"`
	data_folder *data_mnr.DataFolder
}

var SERVER_SQL_TABLE kpsql.SqlTable = sql_mnr.SQLDB.GetTable("servers", &ServerData{})

func NewServer(id string, name string, version string, desc string, addrs []string, group uuid.UUID)(*ServerData){
	addrstr := ""
	if len(addrs) > 0 {
		for _, a := range addrs {
			addrstr += a
			addrstr += ";"
		}
	}
	addrstr = addrstr[0:len(addrstr) - 1]
	return &ServerData{
		Id: id,
		Name: name,
		Version: version,
		Description: desc,
		Addrstr: addrstr,
		Group: group,
	}
}

func GetServer(id string)(*ServerData){
	ins, err := SERVER_SQL_TABLE.SelectPrimary(ServerData{Id: id})
	if err != nil || ins == nil {
		return nil
	}
	return ins.(*ServerData)
}

func GetServerByGroup(gid uuid.UUID)(*ServerData){
	ins, err := SERVER_SQL_TABLE.Select(kpsql.OptWMapEq("group", gid), kpsql.OptLimit(1))
	if err != nil || len(ins) != 1 {
		return nil
	}
	return ins[0].(*ServerData)
}

func GetServerList()(svrlist []*ServerData){
	list, err := SERVER_SQL_TABLE.Select()
	if err != nil {
		return nil
	}
	svrlist = make([]*ServerData, len(list))
	for i, ins := range list {
		svrlist[i] = ins.(*ServerData)
	}
	return
}

func (svr *ServerData)UpdateData(taglist ...string)(err error){
	_, err = SERVER_SQL_TABLE.Update(svr, kpsql.OptTags(taglist...))
	return
}

func (svr *ServerData)InsertData()(err error){
	_, err = SERVER_SQL_TABLE.Insert(svr)
	if err == nil {
		svr.GetDataFolder().Create()
	}
	return
}

func (svr *ServerData)GetDataFolder()(*data_mnr.DataFolder){
	if svr.data_folder == nil {
		svr.data_folder = SERVER_DATA_FOLDER.Folder(svr.Id)
	}
	return svr.data_folder
}

const TOKEN_LENGTH = 1024

func randToken()(token []byte){
	token = make([]byte, TOKEN_LENGTH)
	_, err := crand.Read(token)
	if err != nil {
		panic(err)
	}
	return
}

func (svr *ServerData)GetToken()(token []byte){
	tkfile := svr.GetDataFolder().File("token.txt")
	if tkfile.IsExist() {
		if src, err := tkfile.ReadAll(); err == nil {
			token = make([]byte, base64.StdEncoding.DecodedLen(len(src)))
			_, err = base64.StdEncoding.Decode(token, src)
			if err == nil {
				return
			}
		}
	}
	token = randToken()
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(token)))
	base64.StdEncoding.Encode(buf, token)
	tkfile.Write(buf)
	return
}

func (svr *ServerData)VerifyToken(token []byte)(ok bool){
	right := svr.GetToken()
	leng := len(right)
	if len(token) < leng {
		ok = false
		leng = len(token)
	}
	for i := 0; i < leng; i++ {
		if right[i] != token[i] {
			ok = false
		}
	}
	return
}

func SearchServer(key string)(svrlist []*ServerData){
	list, err := SERVER_SQL_TABLE.Select(kpsql.OptTags("id"),
		kpsql.OptWMapOr("id", "LIKE", key, "name", "LIKE", key, "version", "LIKE", key, "description", "LIKE", key))
	if err != nil {
		return nil
	}
	svrlist = make([]*ServerData, len(list))
	for i, ins := range list {
		svrlist[i] = ins.(*ServerData)
	}
	return
}

