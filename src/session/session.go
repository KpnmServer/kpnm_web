
package kweb_session

import (
	time "time"

	uuid "github.com/google/uuid"
	kpsql "github.com/KpnmServer/go-kpsql"
	page_mnr "github.com/KpnmServer/kpnm_web/src/page_manager"
)

type Session struct{
	Id uuid.UUID `sql:"uuid"`
	Key string   `sql:"key"`
	Value string `sql:"value"`
	Overtime time.Time `sql:"overtime"`
}

var SESSION_SQL_TABLE kpsql.SqlTable = page_mnr.SQLDB.GetTable("sessions", &Session{})

var (
	last_clean_time int64 = 0
	clean_interval int64 = 60 * 60 * 24
	DEFAULT_LIVETIME time.Duration = time.Hour * 24 * 30
)

func NewSession(id uuid.UUID, key string, value string, livetime_ ...time.Duration)(*Session){
	var livetime time.Duration = DEFAULT_LIVETIME
	if len(livetime_) > 0 {
		livetime = livetime_[0]
	}
	return &Session{
		Id: id,
		Key: key,
		Value: value,
		Overtime: time.Now().Add(livetime),
	}
}

func GetSession(id uuid.UUID, key string)(*Session){
	CheckSessions()
	lines, err := SESSION_SQL_TABLE.Select(
		kpsql.WhereMap{{"uuid", "=", id, "AND"}, {"key", "=", key, "AND"}, {"overtime", ">", time.Now(), ""}}, 1)
	if err != nil || len(lines) != 1 {
		return nil
	}
	return lines[0].(*Session)
}

func GetSessionStr(id uuid.UUID, key string)(string){
	ses := GetSession(id, key)
	if ses == nil {
		return ""
	}
	return ses.Value
}

func (session *Session)Save()(err error){
	CheckSessions()
	if session.Overtime.Before(time.Now()) {
		return nil
	}
	if n, _ := SESSION_SQL_TABLE.Count(
		kpsql.WhereMap{{"uuid", "=", session.Id, "AND"}, {"key", "=", session.Key, ""}}, 1); n > 0 {
		_, err = SESSION_SQL_TABLE.Update(session,
			kpsql.WhereMap{{"uuid", "=", session.Id, "AND"}, {"key", "=", session.Key, ""}}, nil, 1)
	}else{
		_, err = SESSION_SQL_TABLE.Insert(session)
	}
	return
}

func (session *Session)Delete()(err error){
	CheckSessions()
	_, err = SESSION_SQL_TABLE.Delete(kpsql.WhereMap{{"uuid", "=", session.Id, "AND"}, {"key", "=", session.Key, ""}}, 1)
	return
}

func (session *Session)Flush(livetime_ ...time.Duration)(*Session){
	var livetime time.Duration = DEFAULT_LIVETIME
	if len(livetime_) > 0 {
		livetime = livetime_[0]
	}
	session.Overtime = time.Now().Add(livetime)
	return session
}

func DelUserSessions(id uuid.UUID)(n int64, err error){
	CheckSessions()
	return SESSION_SQL_TABLE.Delete(kpsql.WhereMap{{"uuid", "=", id, ""}})
}

func CheckSessions()(n int64, err error){
	if last_clean_time + clean_interval < time.Now().Unix() {
		return
	}
	return CleanSessions()
}

func CleanSessions()(n int64, err error){
	last_clean_time = time.Now().Unix()
	return SESSION_SQL_TABLE.Delete(kpsql.WhereMap{{"overtime", "<=", time.Now(), ""}})
}

func ChangeUUID(last uuid.UUID, newo uuid.UUID)(n int64, err error){
	return SESSION_SQL_TABLE.Update(&Session{Id: newo}, kpsql.WhereMap{{"uuid", "=", last, ""}}, []string{"uuid"})
}
