module github.com/zyxgad/kpnm_svr

go 1.15

require (
	github.com/kataras/iris/v12 v12.2.0-alpha2.0.20210717090056-b2cc3a287149
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/zyxgad/go-mc_util v0.0.0
)

replace "github.com/zyxgad/go-mc_util" => "../go-mc_util"
