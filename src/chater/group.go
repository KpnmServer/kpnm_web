
package kweb_chater

import (
	time "time"

	uuid "github.com/google/uuid"
	kpsql "github.com/KpnmServer/go-kpsql"
	sql_mnr "github.com/KpnmServer/kpnm_web/src/sql"
	data_mnr "github.com/KpnmServer/kpnm_web/src/data_manager"
)

var CHAR_DATA_FOLDER = data_mnr.GetDataFolder("chat")

type GroupType = int8

const (
	GROUP_UNKNOWN GroupType = iota
	GROUP_NORMAL
	GROUP_PEOPLE
	GROUP_BOT
)

type Group struct{
	Id uuid.UUID     `sql:"id" sql_primary:"true"`
	Name string      `sql:"name"`
	Owner uuid.UUID  `sql:"owner"`
	Type GroupType   `sql:"type"`
	Desc string      `sql:"description"`
	msg_table kpsql.SqlTable
	data_folder *data_mnr.DataFolder
}

var GROUP_SQL_TABLE kpsql.SqlTable = sql_mnr.SQLDB.GetTable("groups", &Group{})

func NewGroup(name string, owner uuid.UUID, gtype GroupType, desc string)(*Group){
	id := uuid.New()
	return &Group{
		Id: id,
		Name: name,
		Owner: owner,
		Type: gtype,
		Desc: desc,
		data_folder: CHAR_DATA_FOLDER.Folder(id.String()),
	}
}

func GetGroup(id uuid.UUID)(*Group){
	ins, err := GROUP_SQL_TABLE.SelectPrimary(Group{Id: id})
	if err != nil {
		return nil
	}
	return ins.(*Group)
}

func GetGroupsByOwner(owner uuid.UUID)(groups []*Group){
	lines, err := GROUP_SQL_TABLE.Select(kpsql.OptWMapEq("owner", owner))
	if err != nil {
		return nil
	}
	groups = make([]*Group, len(lines))
	for i, ins := range lines {
		groups[i] = ins.(*Group)
	}
	return
}

func (gp *Group)getMsgTable()(kpsql.SqlTable){
	if gp.msg_table == nil {
		gp.msg_table = sql_mnr.SQLDB.GetTableBySqltype("messages_group_" + gp.Id.String(), MESSAGE_SQL_TYPE)
	}
	return gp.msg_table
}

func (gp *Group)GetMessages(last time.Time)(msgs []Message, err error){
	wmap := kpsql.OptWMapAnd("date", ">", last, "isend", "=", true)
	var (
		lines []interface{}
		msgtb = gp.getMsgTable()
		n int
	)
	lines, err = msgtb.Select(wmap, kpsql.OptLimit(100), kpsql.OptOrder("date", 1))
	if err != nil { return }
	n = len(lines)
	msgs = make([]Message, 0, n)
	var (
		msg Message
	)
	for i := 0; i < n ;i++ {
		mdata := lines[i].(*MessageData)
		msg = newAutoMessageByData(mdata, gp)
		_, err = msg.LoadAll()
		if err != nil { return }
		msgs = append(msgs, msg)
		last = msg.MData().Date
	}
	return
}

func (gp *Group)GetFileMsg(id uuid.UUID)(*fileMessage){
	md0, err := gp.getMsgTable().SelectPrimary(MessageData{Id: id})
	if err != nil { return nil }
	md := md0.(*MessageData)
	if md.Type != MSG_FILE && md.Type != MSG_IMG {
		return nil
	}
	return newFileMsgByData(md, gp)
}

func (gp *Group)Insert()(err error){
	_, err = GROUP_SQL_TABLE.Insert(gp)
	if err != nil { return }
	err = gp.getMsgTable().Create()
	if err != nil { return }
	return nil
}

func (gp *Group)Delete()(err error){
	_, err = GROUP_SQL_TABLE.Delete(kpsql.OptWMapEq("id", gp.Id), kpsql.OptLimit(1))
	if err != nil { return }
	err = gp.getMsgTable().Drop()
	if err != nil { return }
	return nil
}

func (gp *Group)SetOwner(id uuid.UUID)(err error){
	gp.Owner = id
	_, err = GROUP_SQL_TABLE.Update(gp, kpsql.OptTags("owner"), kpsql.OptLimit(1))
	return
}

func (gp *Group)SetType(gtype GroupType)(err error){
	gp.Type = gtype
	_, err = GROUP_SQL_TABLE.Update(gp, kpsql.OptTags("type"), kpsql.OptLimit(1))
	return
}

func (gp *Group)GetMembers()(members []*Member, err error){
	var lines []interface{}
	lines, err = MEMBER_SQL_TABLE.Select(kpsql.OptWMapEq("group_id", gp.Id))
	if err != nil {
		return nil, err
	}
	members = make([]*Member, 0, len(lines))
	for _, r := range lines {
		members = append(members, r.(*Member))
	}
	return members, nil
}

func (gp *Group)HasMember(id uuid.UUID)(ok bool){
	n, err := MEMBER_SQL_TABLE.Count(MEMBER_SQL_TABLE.SqlType().PriWhereOpt(Member{GroupId: gp.Id, UserId: id}), kpsql.OptLimit(1))
	if err != nil || n != 1 {
		return false
	}
	return true
}

func (gp *Group)CreateMember(id uuid.UUID, _type ...MemberType)(*Member){
	mtype := MEMBER_NORMAL
	if len(_type) > 0 {
		mtype = _type[0]
	}
	return &Member{
		GroupId: gp.Id,
		UserId: id,
		Type: mtype,
		LastRead: time.Time{},
		group: gp,
	}
}

func (gp *Group)GetMember(id uuid.UUID)(mem *Member){
	ins, err := MEMBER_SQL_TABLE.SelectPrimary(Member{GroupId: gp.Id, UserId: id})
	if err != nil {
		return nil
	}
	if ins == nil {
		return nil
	}
	mem = ins.(*Member)
	mem.group = gp
	return
}

