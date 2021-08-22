
package kpnmmail

import (
	bytes "bytes"
	// texttemp "text/template"
	htmltemp "html/template"
)

var (
	// texttp *texttemp.Template
	htmltp *htmltemp.Template = htmltemp.New("template_html")
)

func LoadHtmlFiles(paths ...string)(err error){
	_, err = htmltp.ParseFiles(paths...)
	return err
}

func ExeHtmlTemp(path string, value interface{})(text string, err error){
	buf := bytes.NewBuffer([]byte{})
	err = htmltp.ExecuteTemplate(buf, path, value)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func init(){
	htmltp.Funcs(htmltemp.FuncMap{
		"odd": func(num int)(bool){ return num % 2 == 0 },
	})
}

