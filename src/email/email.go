
package kweb_email

import (
	os "os"

	email "github.com/KpnmServer/go-util/email"
	ufile "github.com/KpnmServer/go-util/file"
	json "github.com/KpnmServer/go-util/json"
	kutil "github.com/KpnmServer/kpnm_web/src/util"
)

type Map = map[string]interface{}

var svrmail *email.Email

func SendMail(to []string, title string, text string)(err error){
	return svrmail.SendMail(to, title, text)
}

func SendHtml(to []string, title string, path string, _value ...interface{})(err error){
	return svrmail.SendHtml(to, title, path, _value...)
}

func init(){
	{ // read config file
		var fd *os.File
		var err error
		fd, err = os.Open(ufile.JoinPath("config", "email.json"))
		if err != nil {
			panic(err)
		}
		defer fd.Close()

		var obj = make(json.JsonObj)
		err = json.ReadJson(fd, &obj)
		if err != nil {
			panic(err)
		}

		host := obj.GetString("host")
		port := obj.GetInt("port")
		addr := obj.GetString("addr")
		pwd := obj.GetString("pwd")

		svrmail = email.NewEmail(host, port, addr, pwd)
		err = svrmail.Login()
		if err != nil {
			panic(err)
		}
		kutil.RegisterClose(func(){
			svrmail.Close()
		})
	}
	email.LoadTemplateDir("emails")
}

