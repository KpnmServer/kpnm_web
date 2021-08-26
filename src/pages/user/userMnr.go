
package page_user

import (
	time "time"
	crand "crypto/rand"
	sha256 "crypto/sha256"
	strconv "strconv"
	errors "errors"

	uuid "github.com/google/uuid"
	jwt "github.com/KpnmServer/go-util/jwt"
	kpsql "github.com/KpnmServer/go-kpsql"
	kses "github.com/KpnmServer/kpnm_web/src/session"
	kmail "github.com/KpnmServer/kpnm_web/src/email"
	page_mnr "github.com/KpnmServer/kpnm_web/src/page_manager"
)

var JWT_ENCODER jwt.Encoder = jwt.NewAutoEncoder(60 * 60 * 24, 2048)

type UserData struct{
	Userid uuid.UUID   `sqlk:"id"`
	Username string    `sql:"username"`
	Email string       `sql:"email"`
	Password [32]byte  `sql:"password"`
	Frozen int8        `sql:"frozen"`
	Desc string        `sql:"description"`
}

var USER_SQL_TABLE kpsql.SqlTable = page_mnr.SQLDB.GetTable("users", &UserData{})

func NewUser(name string, password string, email string, desc string)(*UserData){
	return &UserData{
		Userid: uuid.New(),
		Username: name,
		Email: email,
		Password: sha256.Sum256(password),
		Frozen: 0,
		Desc: desc,
	}
}

func GetUserData(userid uuid.UUID)(user *UserData){
	obj, err := USER_SQL_TABLE.SelectPrimary(userid)
	if err != nil || obj == nil {
		return nil
	}
	return obj.(*UserData)
}

func GetUserDataByName(name string)(user *UserData){
	obj, err := USER_SQL_TABLE.Select(kpsql.WhereMap{{"username", "=", name}}, 1)
	if err != nil || obj == nil || len(obj) != 1 {
		return nil
	}
	return obj[0].(*UserData)
}

func UpdateUserData(user *UserData, taglist []string)(err error){
	_, err = USER_SQL_TABLE.Update(user, nil, taglist)
	return
}

func InsertUserData(user *UserData)(err error){
	_, err = USER_SQL_TABLE.Insert(user)
	return
}

var MAIL_CODE_LIVE_TIME int64 = 60 * 5

func sendVerifyMail(sesid uuid.UUID, email string)(err error){
	mailc := kses.GetSession(sesid, "mailc")
	if mailc != nil {
		vmail := kses.GetSession(sesid, "vmail")
		mailcot := kses.GetSession(sesid, "mailcot")
		if vmail == nil || vmail.Value != email || mailcot == nil{
			mailc = nil
		}else{
			var n int64
			n, err = strconv.Atoi(mailcot.Value)
			if err != nil || n < time.Now().Unix() {
				mailc = nil
			}
		}
	}
	if mailc == nil {
		var code []byte = make([]byte, 6)
		_, err = crand.Read(code)
		if err != nil { return }
		for i, _ := range code {
			code[i] = '0' + code[i] % 10
		}
		mailc = kses.NewSession(sesid, "mailc", (string)(code))
		err = kses.NewSession(sesid, "mailcot", strconv.Itoa(time.Now().Unix() + MAIL_CODE_LIVE_TIME)).Save()
		if err != nil { return }
		err = kses.NewSession(sesid, "vmail", email).Save()
		if err != nil { return }
	}
	err = kmail.SendHtml(addr, "Verify your email address", "verify-mail.html", mailc.Value)
	return nil
}

func verifyMailCode(sesid uuid.UUID, code string)(tk string, ok bool){
	vmail := kses.GetSession(sesid, "vmail")
	mailc := kses.GetSession(sesid, "mailc")
	mailsendtime := kses.GetSession(sesid, "mailst")
	if vmail == nil || mailc == nil {
		return "", false
	}
	if n, err := strconv.Atoi(mailsendtime.Value);
		err != nil || n + MAIL_CODE_LIVE_TIME < time.Now().Unix() {
		return "", false
	}
	if code != mailc.Value {
		return "", false
	}
	vmail.Delete()
	mailc.Delete()
	mailsendtime.Delete()
	return JWT_ENCODER.Encode(jwt.SetOutdate(jwt.Json{"email": vmail.Value}, time.Minute * 15)), true
}

func checkUserEmailToken(token string)(email string, ok bool){
	obj, err = JWT_ENCODER.Decode(token)
	if err != nil {
		return "", false
	}
	return obj["email"].(string), true
}

