
package kpnmmail

import (
	bufio "bufio"
	tls   "crypto/tls"
	os    "os"

	gomail "gopkg.in/gomail.v2"
	kfutil "github.com/zyxgad/kpnm_svr/src/util/file"
	json "github.com/zyxgad/kpnm_svr/src/util/json"
)


type Email struct{
	host string
	port int
	address string
	password string
}

func NewEmail(host string, port int, address string, password string)(mail *Email){
	mail = new(Email)
	mail.host = host
	mail.port = port
	mail.address = address
	mail.password = password
	return mail
}

func (mail *Email)SendMail(to string, title string, text string)(err error){
	msg := gomail.NewMessage()
	msg.SetHeader("From", mail.address)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", title)
	msg.SetBody("text/html", text)
	dial := gomail.NewDialer(mail.host, mail.port, mail.address, mail.password)
	dial.TLSConfig = &tls.Config{ InsecureSkipVerify: true }
	err = dial.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}

func (mail *Email)SendHtml(to string, title string, path string, value interface{})(err error){
	text, err := ExeHtmlTemp(path, value)
	if err != nil {
		return err
	}
	return mail.SendMail(to, title, text)
}

var svrmail *Email

func SendMail(to string, title string, text string)(err error){
	return svrmail.SendMail(to, title, text)
}

func SendHtml(to string, title string, path string, value interface{})(err error){
	return svrmail.SendHtml(to, title, path, value)
}

func init(){
	{ // read config file
		var fd *os.File
		var err error
		fd, err = os.Open(kfutil.JoinPath("config", "email.json"))
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
		port := obj.GetInt("host")
		addr := obj.GetString("addr")
		pwd := obj.GetString("pwd")

		svrmail = NewEmail(host, port, addr, pwd)
	}
}

