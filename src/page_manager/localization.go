
package kweb_manager

import (
	os "os"
	strings "strings"

	iris "github.com/kataras/iris/v12"
	json "github.com/KpnmServer/go-util/json"
)

type LocalMap map[string]string
type AreaMap struct{
	areaMap map[string]LocalMap
	defaultArea string
	localMap LocalMap
	localArea string
}
type I18nMap struct{
	languages map[string]AreaMap
	defaultLang string
	localAreaMap AreaMap
	localLang string
}

var GLOBAL_I18N_MAP *I18nMap = &I18nMap{
	languages: make(map[string]AreaMap),
	defaultLang: "",
	localLang: "",
}

func GetGlobalI18nMapCopy()(I18nMap){
	return *GLOBAL_I18N_MAP
}

func splitLocalString(local string)(lang string, area string){
	arr := strings.SplitN(strings.ToLower(local), "-", 2)
	lang = arr[0]
	if len(arr) > 1 {
		area = arr[1]
	}else{
		area = ""
	}
	return
}

func (imap *I18nMap)LoadLanguage(local string, path string)(err error){
	var fd *os.File
	fd, err = os.Open(path)
	if err != nil {
		return
	}
	var obj = make(json.JsonObj)
	err = json.ReadJson(fd, &obj)
	if err != nil {
		return
	}
	lang, area := splitLocalString(local)
	APPLICATION.Logger().Debugf("Loading language: '%s/%s'", lang, area)
	areaMap, ok := imap.languages[lang]
	if !ok {
		areaMap = AreaMap{
			areaMap: make(map[string]LocalMap),
			defaultArea: "",
			localArea: "",
		}
	}
	if imap.defaultLang == "" {
		imap.defaultLang = lang
		imap.localAreaMap = areaMap
	}
	if areaMap.defaultArea == "" {
		areaMap.defaultArea = area
	}
	localMap, ok := areaMap.areaMap[area]
	if !ok {
		localMap = make(LocalMap)
	}
	for k, v := range obj {
		localMap[k], _ = v.(string)
	}
	areaMap.areaMap[area] = localMap
	if areaMap.localMap == nil {
		areaMap.localMap = localMap
	}
	imap.languages[lang] = areaMap
	return nil
}

func (imap *I18nMap)getLocalMap(local0 string)(areaMap AreaMap, localMap LocalMap, local string){
	lang, area := splitLocalString(local0)
	areaMap, ok := imap.languages[lang]
	if !ok {
		areaMap = imap.languages[imap.defaultLang]
		lang = imap.defaultLang
		area = ""
	}
	if area == "" {
		area = areaMap.defaultArea
	}
	localMap = areaMap.areaMap[area]
	local = lang
	if area != "" {
		local += "-" + area
	}
	return
}

func (imap *I18nMap)SetLocalLang(local string){
	var localMap LocalMap
	locals := strings.Split(strings.SplitN(local, ";", 2)[0], ",")
	local = locals[0]
	LOGGER.Debugf("User language=%s", local)
	imap.localAreaMap, localMap, imap.localLang = imap.getLocalMap(local)
	imap.localAreaMap.localMap = localMap
	LOGGER.Debugf("Set language=%s", imap.localLang)
}

func (imap *I18nMap)GetLocalLang()(string){
	return imap.localLang
}

func (imap *I18nMap)Localization(id string)(local string){
	id = strings.ToLower(id)
	var ok bool
	local, ok = imap.localAreaMap.localMap[id]
	if ok {
		return
	}
	local, ok = imap.localAreaMap.areaMap[imap.localAreaMap.defaultArea][id]
	if ok {
		return
	}
	areaMap := imap.languages[imap.defaultLang]
	local, ok = areaMap.areaMap[areaMap.defaultArea][id]
	if ok {
		return
	}
	return id
}

func LocalHandle(i18nmap *I18nMap)(iris.Handler){
	return func(ctx iris.Context){
		lang := ctx.Request().Header.Get("Accept-Language")
		i18nmap.SetLocalLang(lang)
		ctx.Next()
	}
}
