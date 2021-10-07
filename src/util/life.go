
package kweb_util

var _CLOSE_HANDLES = make([]func(), 0)

func RegisterClose(call func()){
	_CLOSE_HANDLES = append(_CLOSE_HANDLES, call)
}

func OnClose(){
	for _, h := range _CLOSE_HANDLES {
		h()
	}
}

