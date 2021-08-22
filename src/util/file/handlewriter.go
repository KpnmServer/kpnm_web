
package kweb_util_file

type HandleWriter func([]byte)(int, error)

func (writer HandleWriter)Write(bts []byte)(int, error){
	return writer(bts)
}
