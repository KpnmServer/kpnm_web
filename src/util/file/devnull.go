
package kweb_util_file

type devNull struct{}

func (devNull)Write(bt []byte)(n int, err error){
	return len(bt), nil
}

func (devNull)Read([]byte)(n int, err error){
	return 0, nil
}

var DevNull devNull = devNull{}
