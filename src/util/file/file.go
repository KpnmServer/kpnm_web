
package kweb_util_file

import (
	os    "os"
)

func IsExist(path string)(bool){
	s, err := os.Stat(path)
	return (s != nil) || (err != nil && os.IsExist(err))
}

func IsNotExist(path string)(bool){
	_, err := os.Stat(path)
	return err != nil && os.IsNotExist(err)
}

func IsFile(path string)(bool){
	s, _ := os.Stat(path)
	return s != nil && !s.IsDir()
}

func IsDir(path string)(bool){
	s, _ := os.Stat(path)
	return s != nil && s.IsDir()
}

func RemoveFile(path string)(bool){
	err := os.Remove(path)
	return err == nil
}

func CreateDir(folderPath string){
	if IsNotExist(folderPath){
		os.Mkdir(folderPath, os.ModePerm)
		os.Chmod(folderPath, os.ModePerm)
	}
}

