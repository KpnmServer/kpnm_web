
package kweb_user

import (
	time "time"
	crand "crypto/rand"
	sha256 "crypto/sha256"

	uuid "github.com/google/uuid"
	jwt "github.com/KpnmServer/go-util/jwt"
	kpsql "github.com/KpnmServer/go-kpsql"
	kses "github.com/KpnmServer/kpnm_web/src/session"
	kmail "github.com/KpnmServer/kpnm_web/src/email"
	sql_mnr "github.com/KpnmServer/kpnm_web/src/sql"
)

var JWT_ENCODER jwt.Encoder = jwt.NewAutoEncoder(jwt.NewEncoder(nil), 2048, 60 * 60 * 24)

type UserData struct{
	Id uuid.UUID       `sqlk:"id"`
	Username string    `sql:"username"`
	Email string       `sql:"email"`
	Password [32]byte  `sql:"password"`
	Frozen int8        `sql:"frozen"`
	Desc string        `sql:"description"`
}

var USER_SQL_TABLE kpsql.SqlTable = sql_mnr.SQLDB.GetTable("users", &UserData{})

func NewUser(name string, password string, email string, desc string)(*UserData){
	return &UserData{
		Id: uuid.New(),
		Username: name,
		Email: email,
		Password: sha256.Sum256(([]byte)(password)),
		Frozen: 0,
		Desc: desc,
	}
}

func (user *UserData)CheckPassword(password string)(bool){
	shapwd := sha256.Sum256(([]byte)(password))
	for i := 0; i < len(user.Password) ;i++ {
		if user.Password[i] != shapwd[i] {
			return false
		}
	}
	return true
}

func (user *UserData)SetPassword(password string)(err error){
	user.Password = sha256.Sum256(([]byte)(password))
	return user.UpdateData("password")
}

func GetUserData(userid uuid.UUID)(user *UserData){
	obj, err := USER_SQL_TABLE.SelectPrimary(userid)
	if err != nil || obj == nil {
		return nil
	}
	return obj.(*UserData)
}

func GetUserDataByName(name string)(user *UserData){
	obj, err := USER_SQL_TABLE.Select(kpsql.WhereMap{{"username", "=", name, ""}}, 1)
	if err != nil || obj == nil || len(obj) != 1 {
		return nil
	}
	return obj[0].(*UserData)
}

func GetUserDataByEmail(email string)(user *UserData){
	obj, err := USER_SQL_TABLE.Select(kpsql.WhereMap{{"email", "=", email, ""}}, 1)
	if err != nil || obj == nil || len(obj) != 1 {
		return nil
	}
	return obj[0].(*UserData)
}

func (user *UserData)UpdateData(taglist ...string)(err error){
	_, err = USER_SQL_TABLE.Update(user, nil, taglist)
	return
}

func (user *UserData)InsertData()(err error){
	_, err = USER_SQL_TABLE.Insert(user)
	return
}

var MAIL_CODE_LIVE_TIME = time.Second * 60 * 5

var CODE_BASE []byte = ([]byte)("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
var CODE_LEN int = 6

func SendVerifyMail(sesid uuid.UUID, email string)(err error){
	mailc := kses.GetSession(sesid, "mailc")
	if mailc != nil {
		vmail := kses.GetSession(sesid, "vmail")
		if vmail == nil || vmail.Value != email{
			mailc = nil
		}
	}
	if mailc == nil {
		var code []byte = make([]byte, CODE_LEN)
		_, err = crand.Read(code)
		if err != nil { return }
		for i, _ := range code {
			code[i] = CODE_BASE[code[i] % (byte)(len(CODE_BASE))]
		}
		mailc = kses.NewSession(sesid, "mailc", (string)(code), MAIL_CODE_LIVE_TIME)
		err = mailc.Save()
		if err != nil { return }
		err = kses.NewSession(sesid, "vmail", email, MAIL_CODE_LIVE_TIME).Save()
		if err != nil { return }
	}
	err = kmail.SendHtml(email, "Verify your email address", "register-verify.html", kmail.Map{
		"addr": email,
		"code": mailc.Value,
	})
	return
}

func VerifyMailCode(sesid uuid.UUID, code string)(tk string, ok bool){
	vmail := kses.GetSession(sesid, "vmail")
	mailc := kses.GetSession(sesid, "mailc")
	if vmail == nil || mailc == nil {
		return "", false
	}
	if code != mailc.Value {
		return "", false
	}
	vmail.Delete()
	mailc.Delete()
	return JWT_ENCODER.Encode(jwt.SetOutdate(jwt.Json{"email": vmail.Value}, time.Minute * 15)), true
}

func CheckUserEmailToken(token string)(email string, ok bool){
	obj, err := JWT_ENCODER.Decode(token)
	if err != nil {
		return "", false
	}
	return obj["email"].(string), true
}

