
package kweb_util_json

import (
	io "io"
	ioutil "io/ioutil"
	errors "errors"
	json "encoding/json"
)

func EncodeJson(obj interface{})(data []byte){
	data, err := json.Marshal(obj)
	if err != nil {
		return nil
	}
	return data
}

func DecodeJson(data []byte, obj_p interface{})(err error){
	return json.Unmarshal(data, obj_p)
}

func EncodeJsonStr(obj interface{})(data string){
	bts := EncodeJson(obj)
	if bts == nil {
		return ""
	}
	return string(bts)
}

func DecodeJsonStr(data string, obj_p interface{})(err error){
	return DecodeJson(([]byte)(data), obj_p)
}

func ReadJson(r io.Reader, obj_p interface{})(err error){
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return DecodeJson(data, obj_p)
}

func WriteJson(w io.Writer, obj interface{})(err error){
	data := EncodeJson(obj)
	if data != nil {
		return errors.New("Encode Error")
	}
	return nil
}

