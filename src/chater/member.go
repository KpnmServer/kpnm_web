
package kweb_chater

import (
	time "time"

	uuid "github.com/google/uuid"
	kpsql "github.com/KpnmServer/go-kpsql"
	sql_mnr "github.com/KpnmServer/kpnm_web/src/sql"
)


type MemberType int8

const (
	MEMBER_NORMAL MemberType = iota
	MEMBER_OWNER
	MEMBER_OPERATOR
)

type Member struct{
	GroupId uuid.UUID  `sql:"group_id"`
	UserId uuid.UUID   `sql:"user_id"`
	Type MemberType    `sql:"type"`
	LastRead time.Time `sql:"last_read"`
	group *Group
}

var MEMBER_SQL_TABLE kpsql.SqlTable = sql_mnr.SQLDB.GetTable("members", &Member{})

func (mem *Member)getGroup()(*Group){
	if mem.group == nil {
		gp, _ := GROUP_SQL_TABLE.SelectPrimary(Group{Id: mem.GroupId})
		if gp != nil {
			mem.group = gp.(*Group)
		}
	}
	return mem.group
}

func (mem *Member)Insert()(err error){
	_, err = MEMBER_SQL_TABLE.Insert(mem)
	return
}

func (mem *Member)Delete()(err error){
	_, err = MEMBER_SQL_TABLE.Delete(kpsql.OptWMapEqAnd("group_id", mem.GroupId, "user_id", mem.UserId), kpsql.OptLimit(1))
	return
}

func (mem *Member)SetType(mtype MemberType)(err error){
	mem.Type = mtype
	_, err = MEMBER_SQL_TABLE.Update(mem, kpsql.OptTags("type"), kpsql.OptLimit(1))
	return
}

func (mem *Member)IsOwner()(bool){
	return mem.Type == MEMBER_OWNER
}

func (mem *Member)IsOperator()(bool){
	return mem.Type == MEMBER_OWNER || mem.Type == MEMBER_OPERATOR
}

func (mem *Member)NewMessage(mtype MsgType, data string)(Message){
	return newMessage(mtype, data, mem.UserId, mem.getGroup())
}

func (mem *Member)GetUnreadMsgs()(msgs []Message, err error){
	tmnow := time.Now()
	msgs, err = mem.group.GetMessages(mem.LastRead)
	if err == nil {
		mem.LastRead = tmnow
		MEMBER_SQL_TABLE.Update(mem, kpsql.OptTags("last_read"), kpsql.OptLimit(1))
	}
	return
}
