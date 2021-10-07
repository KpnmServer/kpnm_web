
package kweb_data_mnr

import (
	os "os"
	io "io"
	ioutil "io/ioutil"

	ufile "github.com/KpnmServer/go-util/file"
	json "github.com/KpnmServer/go-util/json"
)

var (
	DATA_BASE_PATH string = "data"
)

type DataFolder struct{
	base string
	path string
}

func GetDataFolder(name string)(*DataFolder){
	path := ufile.JoinPathWithoutAbs(DATA_BASE_PATH, name)
	if ufile.IsNotExist(path) {
		ufile.CreateDir(path)
	}
	return &DataFolder{
		base: ".",
		path: path,
	}
}

func (df *DataFolder)Path()(string){
	return df.path
}

func (df *DataFolder)IsExist()(bool){
	return ufile.IsExist(df.path)
}

func (df *DataFolder)List()(names []string, err error){
	files, err := df.ListInfo()
	if err != nil { return }
	names = make([]string, 0, len(files))
	for _, f := range files {
		names = append(names, f.Name())
	}
	return
}

func (df *DataFolder)ListInfo()(infos []os.FileInfo, err error){
	return ioutil.ReadDir(df.path)
}

func (df *DataFolder)Folder(name string)(*DataFolder){
	path := ufile.JoinPathWithoutAbs(df.path, name)
	if ufile.IsNotExist(df.base) {
		ufile.CreateDir(df.base)
	}
	return &DataFolder{
		base: df.path,
		path: path,
	}
}

func (df *DataFolder)Create()(err error){
	return ufile.CreateDir(df.path)
}

func (df *DataFolder)Remove()(err error){
	var infos []os.FileInfo
	infos, err = df.ListInfo()
	if err != nil { return }
	for _, f := range infos {
		path := ufile.JoinPathWithoutAbs(df.path, f.Name())
		if f.IsDir() {
			err = (&DataFolder{path: path}).Remove()
		}else{
			err = (&DataFile{path: path}).Remove()
		}
		if err != nil { return }
	}
	return nil
}

type DataFile struct{
	base string
	path string
}

func (df *DataFolder)File(name string)(*DataFile){
	path := ufile.JoinPathWithoutAbs(df.path, name)
	return &DataFile{
		base: df.path,
		path: path,
	}
}

func (df *DataFile)Path()(string){
	return df.path
}

func (df *DataFile)IsExist()(bool){
	return ufile.IsExist(df.path)
}

func (df *DataFile)Stat()(os.FileInfo, error){
	return os.Stat(df.path)
}

func (df *DataFile)WriteFunc(handle func(*os.File)(error))(err error){
	if ufile.IsNotExist(df.base) {
		ufile.CreateDir(df.base)
	}
	var fd *os.File
	fd, err = os.Create(df.path)
	if err != nil { return }
	defer fd.Close()
	return handle(fd)
}

func (df *DataFile)ReadFunc(handle func(*os.File)(error))(err error){
	var fd *os.File
	fd, err = os.Open(df.path)
	if err != nil { return }
	defer fd.Close()
	return handle(fd)
}

func (df *DataFile)Write(data []byte)(n int, err error){
	err = df.WriteFunc(func(fd *os.File)(err error){
		n, err = fd.Write(data)
		return
	})
	return
}

func (df *DataFile)Read(buf []byte)(n int, err error){
	err = df.ReadFunc(func(fd *os.File)(err error){
		n, err = fd.Read(buf)
		return
	})
	return
}

func (df *DataFile)ReadFrom(r io.Reader)(n int64, err error){
	err = df.WriteFunc(func(fd *os.File)(err error){
		n, err = io.Copy(fd, r)
		return
	})
	return
}

func (df *DataFile)WriteTo(w io.Writer)(n int64, err error){
	err = df.ReadFunc(func(fd *os.File)(err error){
		n, err = io.Copy(w, fd)
		return
	})
	return
}

func (df *DataFile)ReadAll()(data []byte, err error){
	err = df.ReadFunc(func(fd *os.File)(err error){
		data, err = ioutil.ReadAll(fd)
		return
	})
	return
}

func (df *DataFile)WriteJson(obj interface{})(n int, err error){
	err = df.WriteFunc(func(fd *os.File)(err error){
		n, err = json.WriteJson(fd, obj)
		return
	})
	return
}

func (df *DataFile)ReadJson(obj_p interface{})(err error){
	return df.ReadFunc(func(fd *os.File)(error){
		return json.ReadJson(fd, obj_p)
	})
}

func (df *DataFile)Remove()(err error){
	if ufile.IsExist(df.path) {
		return os.Remove(df.path)
	}
	return nil
}
