
package kweb_util_json

type JsonObj map[string]interface{}

func (obj JsonObj)Get(key string)(v interface{}){
	return obj[key]
}

func (obj JsonObj)GetBool(key string)(v bool){
	return obj[key].(bool)
}

func (obj JsonObj)GetByte(key string)(v byte){
	return (byte)(obj[key].(float64))
}

func (obj JsonObj)GetInt(key string)(v int){
	return (int)(obj[key].(float64))
}

func (obj JsonObj)GetUInt(key string)(v uint){
	return (uint)(obj[key].(float64))
}

func (obj JsonObj)GetInt8(key string)(v int8){
	return (int8)(obj[key].(float64))
}

func (obj JsonObj)GetInt16(key string)(v int16){
	return (int16)(obj[key].(float64))
}

func (obj JsonObj)GetInt32(key string)(v int32){
	return (int32)(obj[key].(float64))
}

func (obj JsonObj)GetInt64(key string)(v int64){
	return (int64)(obj[key].(float64))
}

func (obj JsonObj)GetUInt8(key string)(v uint8){
	return (uint8)(obj[key].(float64))
}

func (obj JsonObj)GetUInt16(key string)(v uint16){
	return (uint16)(obj[key].(float64))
}

func (obj JsonObj)GetUInt32(key string)(v uint32){
	return (uint32)(obj[key].(float64))
}

func (obj JsonObj)GetUInt64(key string)(v uint64){
	return (uint64)(obj[key].(float64))
}

func (obj JsonObj)GetFloat32(key string)(v float32){
	return (float32)(obj[key].(float64))
}

func (obj JsonObj)GetFloat64(key string)(v float64){
	return obj[key].(float64)
}

func (obj JsonObj)GetString(key string)(v string){
	return obj[key].(string)
}

func (obj JsonObj)GetArray(key string)(v JsonArr){
	return (JsonArr)(obj[key].([]interface{}))
}

func (obj JsonObj)GetBytes(key string)(v []byte){
	arr := obj[key].([]interface{})
	v = make([]byte, len(arr))
	for i, _ := range arr {
		v[i] = (byte)(arr[i].(float64))
	}
	return
}

func (obj JsonObj)GetStrings(key string)(v []string){
	arr := obj[key].([]interface{})
	v = make([]string, len(arr))
	for i, _ := range arr {
		v[i] = arr[i].(string)
	}
	return
}

func (obj JsonObj)GetArrays(key string)(v []JsonArr){
	arr := obj[key].([]interface{})
	v = make([]JsonArr, len(arr))
	for i, _ := range arr {
		v[i] = (JsonArr)(arr[i].([]interface{}))
	}
	return
}

func (obj JsonObj)GetObjs(key string)(v []JsonObj){
	arr := obj[key].([]interface{})
	v = make([]JsonObj, len(arr))
	for i, _ := range arr {
		v[i] = (JsonObj)(arr[i].(map[string]interface{}))
	}
	return
}

func (obj JsonObj)GetObj(key string)(v JsonObj){
	return (JsonObj)(obj[key].(map[string]interface{}))
}

func (obj JsonObj)GetStringMap(key string)(v map[string]string){
	return obj[key].(map[string]string)
}

