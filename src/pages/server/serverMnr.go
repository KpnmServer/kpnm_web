
package page_server

import (
	os "os"
	ioutil "io/ioutil"

	kfutil "github.com/KpnmServer/go-util/file"
	json "github.com/KpnmServer/go-util/json"
)

var SERVER_DATA_PATH string = "./data/server"

type ServerInfo struct{
	Name string
	Version string
	Description string
	Addrs []json.JsonArr
}

type serverCache struct{
	mtime int64
	info *ServerInfo
}

var SERVER_CACHE = make(map[string]*serverCache)

func GetServerInfo(name string)(svr *ServerInfo, err error){
	path := kfutil.JoinPathWithoutAbs(SERVER_DATA_PATH, name, "info.json")

	file_stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	var cache = &serverCache{
		mtime: file_stat.ModTime().Unix(),
		info: nil,
	}

	if che, ok := SERVER_CACHE[name]; ok && che.mtime != 0 && che.mtime == cache.mtime{
		return che.info, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var obj = make(json.JsonObj)
	err = json.ReadJson(file, &obj)
	if err != nil {
		return nil, err
	}
	cache.info = &ServerInfo{
		Name: obj.GetString("name"),
		Version: obj.GetString("version"),
		Description: obj.GetString("desc"),
		Addrs: obj.GetArrays("addrs"),
	}
	SERVER_CACHE[name] = cache
	return cache.info, nil
}

func SetServerInfo(svr *ServerInfo)(err error){
	delete(SERVER_CACHE, svr.Name)
	kfutil.CreateDir(kfutil.JoinPathWithoutAbs(SERVER_DATA_PATH, svr.Name))
	path := kfutil.JoinPathWithoutAbs(SERVER_DATA_PATH, svr.Name, "info.json")
	file, err := os.OpenFile(path, os.O_CREATE | os.O_WRONLY | os.O_SYNC, 0600)
	if err != nil {
		return err
	}
	err = json.WriteJson(file, json.JsonObj{
		"name": svr.Name,
		"desc": svr.Description,
		"addrs": svr.Addrs,
	})
	return
}

func GetServerReadme(name string)(data []byte, err error){
	file, err := os.Open(kfutil.JoinPathWithoutAbs(SERVER_DATA_PATH, name, "README.MD"))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err = ioutil.ReadAll(file)
	return
}

