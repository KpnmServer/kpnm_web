
package kweb_captchaimg

import (
	bytes "bytes"
	base64 "encoding/base64"
	time "time"

	captcha "github.com/dchest/captcha"
	kpsql "github.com/KpnmServer/go-kpsql"
	sql_mnr "github.com/KpnmServer/kpnm_web/src/sql"
)

const (
	CAPT_LENGTH = 6
	CAPT_WIDTH = 240
	CAPT_HEIGHT = 100
)

type captchaData struct{
	Id string          `sql:"id" sql_primary:"true"`
	Value string       `sql:"value"`
	Overtime time.Time `sql:"overtime"`
}

var (
	sqlCaptTable kpsql.SqlTable = sql_mnr.SQLDB.GetTable("captchas", &captchaData{})
)

func NewCaptcha()(id string, imgdata string, err error){
	id = captcha.NewLen(CAPT_LENGTH)

	var imgbuf *bytes.Buffer = bytes.NewBuffer([]byte{})
	err = captcha.WriteImage(imgbuf, id, CAPT_WIDTH, CAPT_HEIGHT)
	if err != nil {
		return "", "", err
	}

	imgdata = "data:image/png;base64," + base64.StdEncoding.EncodeToString(imgbuf.Bytes())
	return id, imgdata, nil
}

func VerifyCaptcha(id string, value string)(ok bool){
	return captcha.VerifyString(id, value)
}

func RemoveCaptcha(id string)(ok bool){
	_, err := sqlCaptTable.Delete(kpsql.OptWMapEq("id", id))
	return err != nil
}

type sqlCaptStore struct{
	sqltb kpsql.SqlTable
}

func (cst *sqlCaptStore)Set(id string, digits []byte){
	var digstr []byte = make([]byte, len(digits))
	for i, d := range digits {
		digstr[i] = '0' + d
	}
	cst.sqltb.Insert(&captchaData{
		Id: id,
		Value: (string)(digstr),
		Overtime: time.Now().Add(time.Minute * 5),
	})
}

func (cst *sqlCaptStore)Get(id string, clear bool)(digits []byte){
	capt, err := cst.sqltb.SelectPrimary(captchaData{Id:id})
	if err != nil || capt == nil {
		return nil
	}
	if clear {
		cst.sqltb.Delete(kpsql.OptWMapEq("id", id))
	}
	value := capt.(*captchaData).Value
	digits = make([]byte, len(value))
	for i, d := range ([]byte)(value) {
		if d < '0' || '9' < d {
			return nil
		}
		digits[i] = d - '0'
	}
	return digits
}

func init(){
	sqlCaptTable.Delete()
	captcha.SetCustomStore(&sqlCaptStore{sqltb: sqlCaptTable})
}
