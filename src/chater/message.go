
package kweb_chater

import (
	io "io"
	time "time"

	uuid "github.com/google/uuid"
	json "github.com/KpnmServer/go-util/json"
	kpsql "github.com/KpnmServer/go-kpsql"
	data_mnr "github.com/KpnmServer/kpnm_web/src/data_manager"
)

type MsgType int8

const (
	MSG_KNOWN MsgType = iota
	MSG_TEXT
	MSG_SHORT_TEXT
	MSG_AT
	MSG_FILE
	MSG_IMG
)

type MessageData struct{
	Id uuid.UUID     `sql:"id" sqlword:"CHAR(32)" sql_primary:"true"`
	Date time.Time   `sql:"date" sqlword:"DATETIME"`
	Owner uuid.UUID  `sql:"owner" sqlword:"CHAR(32)"`
	Type MsgType     `sql:"type" sqlword:"INTEGER(1) UNSIGNED"`
	Data string      `sql:"data" sqlword:"TEXT"`
	SData string     `sql:"sdata" sqlword:"VARCHAR(255)"`
	NextId uuid.UUID `sql:"nextid" sqlword:"CHAR(32)" sql_foreign:"messages:id"`
	Isend bool       `sql:"isend" sqlword:"BOOLEAN"`
}

func NewMessageData(mtype MsgType, data string, owner uuid.UUID)(msg *MessageData){
	id := uuid.New()
	msg = &MessageData{
		Id: id,
		Owner: owner,
		Type: mtype,
		NextId: id,
	}
	switch msg.Type {
	case MSG_AT, MSG_FILE, MSG_IMG:
		if len(data) >= 255 {
			panic("len(data) >= 255")
		}
		msg.SData = data
	case MSG_SHORT_TEXT:
		if len(data) >= 255 {
			msg.Type = MSG_TEXT
			msg.Data = data
		}else{
			msg.SData = data
		}
	case MSG_TEXT:
		if len(data) < 255 {
			msg.Type = MSG_SHORT_TEXT
			msg.SData = data
		}else{
			msg.Data = data
		}
	}
	return
}

func (msg *MessageData)GetData()(string){
	switch msg.Type{
	case MSG_SHORT_TEXT, MSG_AT, MSG_FILE, MSG_IMG:
		return msg.SData
	case MSG_TEXT:
		return msg.Data
	}
	return ""
}

var MESSAGE_SQL_TYPE *kpsql.SqlType = kpsql.NewSqlType(&MessageData{})

type Message interface{
	MData()(*MessageData)
	Group()(*Group)
	setNext(Message)
	getNext()(Message)
	GetEnd()(Message)
	SetDate(time.Time)
	SetOwner(uuid.UUID)
	setGroup(*Group)
	Now()(Message)
	Append(Message)(Message)
	AppendMsg(MsgType, string)(Message)
	LoadAll()(n int64, err error)
	InsertAll()(n int64, err error)
	String()(string)
	json()(json.JsonObj)
	Json()(json.JsonObj)
}

type message struct{
	*MessageData
	group *Group
	Next Message
}

func newMessageByData(msg0 *MessageData, group *Group)(*message){
	return &message{
		MessageData: msg0,
		group: group,
		Next: nil,
	}
}

func newMessage(mtype MsgType, data string, owner uuid.UUID, group *Group)(*message){
	return &message{
		MessageData: NewMessageData(mtype, data, owner),
		group: group,
		Next: nil,
	}
}

func NewMsg(mtype MsgType, data string)(*message){
	return newMessage(mtype, data, uuid.Nil, nil)
}

func (msg *message)MData()(*MessageData){
	return msg.MessageData
}

func (msg *message)Group()(*Group){
	return msg.group
}


func (msg *message)Now()(Message){
	msg.SetDate(time.Now())
	return msg
}

func (msg *message)setNext(next Message){
	msg.Next = next
}

func (msg *message)getNext()(Message){
	return msg.Next
}

func (msg *message)GetEnd()(end Message){
	end = msg
	for end.getNext() != nil && end.getNext() != end {
		end = end.getNext()
	}
	return
}

func (msg *message)SetDate(t time.Time){
	msg.MData().Date = t
	var tg Message = msg
	for tg.getNext() != nil && tg.getNext() != tg {
		tg = tg.getNext()
		tg.MData().Date = t
	}
}

func (msg *message)SetOwner(id uuid.UUID){
	msg.Owner = id
	var tg Message = msg
	for tg.getNext() != nil && tg.getNext() != tg {
		tg = tg.getNext()
		tg.MData().Owner = id
	}
}

func (msg *message)setGroup(group *Group){
	msg.group = group
}

func (msg *message)Append(other Message)(Message){
	other.SetDate(msg.Date)
	other.SetOwner(msg.Owner)
	other.setGroup(msg.group)
	end := msg.GetEnd()
	end.MData().NextId = other.MData().Id
	end.setNext(other)
	return msg
}

func (msg *message)AppendMsg(mtype MsgType, data string)(Message){
	return msg.Append(NewMsg(mtype, data))
}

func (msg *message)LoadAll()(n int64, err error){
	n, err = 0, nil
	msg.Next = nil
	if msg.NextId == msg.Id {
		return
	}
	var v interface{}
	v, err = msg.group.getMsgTable().SelectPrimary(MessageData{Id: msg.NextId})
	if err != nil { return }
	if v != nil {
		msg.Next = &message{
			MessageData: v.(*MessageData), group: msg.group, Next: nil,
		}
		n, err = msg.Next.LoadAll()
		n += 1
		return
	}
	return
}

type sqlinsertable interface{
	insertAll(tx *kpsql.SqlTx)(n int64, err error)
	InsertAll()(n int64, err error)
}

func (msg *message)insertAll(tx *kpsql.SqlTx)(n int64, err error){
	if msg.NextId != msg.Id && msg.Next != nil {
		next := msg.Next.(sqlinsertable)
		msg.Next.MData().Isend = false
		n, err = next.insertAll(tx)
		if err != nil { return }
	}
	_, err = tx.Insert(msg.MessageData)
	if err != nil { return }
	n += 1
	return
}

func (msg *message)InsertAll()(n int64, err error){
	msgtable := msg.group.getMsgTable()
	tx, err := msgtable.Begin()
	if err != nil { return }
	defer func(){
		if err == nil {
			tx.Commit()
		}else{
			tx.Rollback()
		}
	}()
	msg.Isend = true
	n, err = msg.insertAll(tx)
	return
}

func (msg *message)String()(str string){
	switch msg.Type {
	case MSG_TEXT, MSG_SHORT_TEXT:
		str = msg.GetData()
	case MSG_AT:
		str = "{at:" + msg.GetData() + "}"
	case MSG_FILE:
		str = "{file:" + msg.Id.String() + "/" + msg.GetData() + "}"
	case MSG_IMG:
		str = "{image:" + msg.Id.String() + "}"
	}
	if msg.Next != nil && msg.Next != msg {
		str += msg.Next.String()
	}
	return
}

func (msg *message)json()(obj json.JsonObj){
	obj = make(json.JsonObj)
	obj["id"] = msg.Id.String()
	obj["type"] = msg.Type
	obj["data"] = msg.GetData()
	return
}

func (msg *message)Json()(obj json.JsonObj){
	obj = make(json.JsonObj)
	obj["date"] = formatTime(msg.Date)
	obj["owner"] = msg.Owner.String()
	if msg.Next == nil || msg.Next == msg {
		obj["id"] = msg.Id.String()
		obj["type"] = msg.Type
		obj["data"] = msg.GetData()
		return
	}
	obj["list"] = true
	data := make([]json.JsonObj, 1, 3)
	data[0] = msg.json()
	var tg Message = msg
	for tg.MData().NextId != tg.MData().Id && tg.getNext() != nil {
		tg = tg.getNext()
		data = append(data, tg.json())
	}
	obj["data"] = data
	return
}


type fileMessage struct{
	*message
	file *data_mnr.DataFile
}

func newFileMessage(mtype MsgType, name string, owner uuid.UUID, group *Group)(msg *fileMessage){
	if mtype != MSG_FILE || mtype != MSG_IMG {
		panic("mtype != MSG_FILE || mtype != MSG_IMG")
	}
	msg0 := newMessage(mtype, name, owner, group)
	return &fileMessage{
		message: msg0,
		file: group.data_folder.File(msg0.Id.String()),
	}
}

func newFileMsgByData(msg0 *MessageData, group *Group)(msg *fileMessage){
	if msg0.Type != MSG_FILE || msg0.Type != MSG_IMG {
		panic("msg0.Type != MSG_FILE || msg0.Type != MSG_IMG")
	}
	return &fileMessage{
		message: newMessageByData(msg0, group),
		file: group.data_folder.File(msg0.Id.String()),
	}
}

func newFileMsgByMsg(old Message)(msg *fileMessage){
	var ok bool
	if msg, ok = old.(*fileMessage); ok{
		return
	}
	if old.MData().Type != MSG_FILE || old.MData().Type != MSG_IMG {
		panic("old.Type != MSG_FILE || old.Type != MSG_IMG")
	}
	msg = &fileMessage{
		message: old.(*message),
		file: old.Group().data_folder.File(old.MData().Id.String()),
	}
	return
}

func (msg *fileMessage)HasData()(bool){
	return msg.file.IsExist()
}

func (msg *fileMessage)WriteData(data []byte)(err error){
	_, err = msg.file.Write(data)
	if err != nil {
		msg.file.Remove()
	}
	return
}

func (msg *fileMessage)SetData(data []byte)(*fileMessage){
	msg.WriteData(data)
	return msg
}

func (msg *fileMessage)ReadFrom(r io.Reader)(n int64, err error){
	n, err = msg.file.ReadFrom(r)
	if err != nil {
		msg.file.Remove()
	}
	return
}

func (msg *fileMessage)WriteTo(w io.Writer)(n int64, err error){
	return msg.file.WriteTo(w)
}

func (msg *fileMessage)ReadData()(data []byte, err error){
	return msg.file.ReadAll()
}

func newAutoMessageByData(msg0 *MessageData, group *Group)(Message){
	switch msg0.Type {
	case MSG_FILE, MSG_IMG:
		return newFileMsgByData(msg0, group)
	case MSG_TEXT, MSG_SHORT_TEXT, MSG_AT:
		return newMessageByData(msg0, group)
	default:
		panic("Unknown message type")
	}
}

func formatTime(t time.Time)(string){
	return t.Format("2006-01-02T15:04:05.000Z-0700")
}
