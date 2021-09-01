module github.com/KpnmServer/kpnm_web

go 1.15

require (
	github.com/KpnmServer/go-kpsql v0.0.0-00010101000000-000000000000
	github.com/KpnmServer/go-mc_util v0.0.0
	github.com/KpnmServer/go-util v0.0.0-00010101000000-000000000000
	github.com/dchest/captcha v0.0.0-20200903113550-03f5f0333e1f
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/uuid v1.3.0
	github.com/kataras/golog v0.1.7
	github.com/kataras/iris/v12 v12.2.0-alpha2.0.20210717090056-b2cc3a287149
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/net v0.0.0-20210825183410-e898025ed96a // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // indirect
)

replace github.com/KpnmServer/go-util => ../go-util

replace github.com/KpnmServer/go-mc_util => ../go-mc_util

replace github.com/KpnmServer/go-kpsql => ../go-kpsql
