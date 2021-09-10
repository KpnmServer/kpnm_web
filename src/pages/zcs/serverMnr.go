
package page_zcs

type serverStatus uint

const (
	SERVER_UNKNOWN serverStatus = iota
	SERVER_STARTING
	SERVER_STARTED
	SERVER_STOPPING
	SERVER_STOPPED
)

func (s serverStatus)String()(string){
	switch s {
	case SERVER_STARTING: return "STARTING"
	case SERVER_STARTED: return "RUNNING"
	case SERVER_STOPPING: return "STOPPING"
	case SERVER_STOPPED: return "STOPPED"
	case SERVER_UNKNOWN: fallthrough
	default: return "UNKNOWN"
	}
}

type ServerInfo struct{
	status serverStatus
	ticks int
	cpu_num uint
	java_version string
	os string
	max_mem uint64
	total_mem uint64
	used_mem uint64
	cpu_load float64
	cpu_time float64
	errstr string
	id string
	interval uint
}

var ZCS_SVR_INFOS = map[string]*ServerInfo{
	"main": &ServerInfo{
		status: SERVER_UNKNOWN,
		interval: 1,
	},
	"mirror": &ServerInfo{
		status: SERVER_UNKNOWN,
		interval: 1,
	},
}
