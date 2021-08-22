
package kweb_util_json


type JsonArr []interface{}

func (arr JsonArr)Get(ind int)(v interface{}){
	return arr[ind]
}

func (arr JsonArr)GetBool(ind int)(v bool){
	return arr[ind].(bool)
}

func (arr JsonArr)GetByte(ind int)(v byte){
	return (byte)(arr[ind].(float64))
}

func (arr JsonArr)GetInt(ind int)(v int){
	return (int)(arr[ind].(float64))
}

func (arr JsonArr)GetUInt(ind int)(v uint){
	return (uint)(arr[ind].(float64))
}

func (arr JsonArr)GetInt8(ind int)(v int8){
	return (int8)(arr[ind].(float64))
}

func (arr JsonArr)GetInt16(ind int)(v int16){
	return (int16)(arr[ind].(float64))
}

func (arr JsonArr)GetInt32(ind int)(v int32){
	return (int32)(arr[ind].(float64))
}

func (arr JsonArr)GetInt64(ind int)(v int64){
	return (int64)(arr[ind].(float64))
}

func (arr JsonArr)GetUInt8(ind int)(v uint8){
	return (uint8)(arr[ind].(float64))
}

func (arr JsonArr)GetUInt16(ind int)(v uint16){
	return (uint16)(arr[ind].(float64))
}

func (arr JsonArr)GetUInt32(ind int)(v uint32){
	return (uint32)(arr[ind].(float64))
}

func (arr JsonArr)GetUInt64(ind int)(v uint64){
	return (uint64)(arr[ind].(float64))
}

func (arr JsonArr)GetFloat32(ind int)(v float32){
	return (float32)(arr[ind].(float64))
}

func (arr JsonArr)GetFloat64(ind int)(v float64){
	return arr[ind].(float64)
}

func (arr JsonArr)GetString(ind int)(v string){
	return arr[ind].(string)
}

func (arr JsonArr)GetArray(ind int)(v JsonArr){
	return (JsonArr)(arr[ind].([]interface{}))
}

func (arr JsonArr)GetBytes(ind int)(v []byte){
	a := arr[ind].([]interface{})
	v = make([]byte, len(a))
	for i, _ := range a {
		v[i] = (byte)(a[i].(float64))
	}
	return
}

func (arr JsonArr)GetStrings(ind int)(v []string){
	a := arr[ind].([]interface{})
	v = make([]string, len(a))
	for i, _ := range a {
		v[i] = a[i].(string)
	}
	return
}

func (arr JsonArr)GetArrays(ind int)(v []JsonArr){
	a := arr[ind].([]interface{})
	v = make([]JsonArr, len(a))
	for i, _ := range a {
		v[i] = (JsonArr)(a[i].([]interface{}))
	}
	return
}

func (arr JsonArr)GetObjs(ind int)(v []JsonObj){
	a := arr[ind].([]interface{})
	v = make([]JsonObj, len(a))
	for i, _ := range a {
		v[i] = (JsonObj)(a[i].(map[string]interface{}))
	}
	return
}

func (arr JsonArr)GetObj(ind int)(v JsonObj){
	return (JsonObj)(arr[ind].(map[string]interface{}))
}

func (arr JsonArr)GetStringMap(ind int)(v map[string]string){
	return arr[ind].(map[string]string)
}

