
package page_user

import (
	regexp  "regexp"

	iris "github.com/kataras/iris/v12"
	ufile "github.com/KpnmServer/go-util/file"
	kpsql "github.com/KpnmServer/go-kpsql"
	kses "github.com/KpnmServer/kpnm_web/src/session"
	kcapt "github.com/KpnmServer/kpnm_web/src/captchaimg"
	kuser "github.com/KpnmServer/kpnm_web/src/user"
)

var (
	reg_name  *regexp.Regexp = regexp.MustCompile(`^[A-Za-z_-][0-9A-Za-z_-]{1,31}$`)
	reg_pwd   *regexp.Regexp = regexp.MustCompile(`^[A-Za-z][0-9A-Za-z_+\-*/!@#$%^&()~\[\]{}|=,.<>;:'"]{7,127}$`)
	reg_email *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`)
)

func UserHeadApi(ctx iris.Context){
	userid := ctx.Params().Get("id")
	path := ufile.JoinPathWithoutAbs(USER_DATA_PATH, userid, "head.png")
	ctx.ServeFile(path)
}

func UserLoginApi(ctx iris.Context){
	uid := kses.GetCtxUuid(ctx)
	cid := kses.GetSession(uid, "captid")
	if cid != nil {
		cid.Delete()
	}
	captcode := ctx.PostValue("capt")
	if cid == nil || !kcapt.VerifyCaptcha(cid.Value, captcode) {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "CaptcodeError",
		})
		return
	}
	username := ctx.PostValue("user")
	var user *kuser.UserData
	switch {
	case reg_email.MatchString(username):
		user = kuser.GetUserDataByEmail(username)
	case reg_name.MatchString(username):
		user = kuser.GetUserDataByName(username)
	default:
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "IllegalDataError",
			"errorMessage": "Username is illegal data",
		})
		return
	}
	if user == nil {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "UserNotExistError",
			"errorMessage": "User not exists",
		})
		return
	}

	password := ctx.PostValue("pwd")
	if !user.CheckPassword(password) {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "PasswordError",
			"errorMessage": "Password is wrong",
		})
		return
	}

	live_time := kuser.LOG_LIVE_MAX_TIME
	if ctx.PostValue("live") != "T" {
		// live_time = 0
	}
	err := user.SaveCtxLog(ctx, live_time)
	if err != nil {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "SaveCtxLogError",
			"errorMessage": err.Error(),
		})
		return
	}
	ctx.JSON(iris.Map{"status": "ok"})
}

func UserLogoutApi(ctx iris.Context){
	kses.DelCtxUuid(ctx)
	ctx.StatusCode(iris.StatusNoContent)
}

func VerifyEmailApi(ctx iris.Context){
	uid := kses.GetCtxUuid(ctx)
	if vf := kses.GetSession(uid, "verify_flag"); vf != nil {
		cid := kses.GetSession(uid, "captid")
		if cid != nil {
			cid.Delete()
		}
		captcode := ctx.PostValue("capt")
		if cid == nil || !kcapt.VerifyCaptcha(cid.Value, captcode) {
			ctx.JSON(iris.Map{
				"status": "error",
				"error": "CaptcodeError",
			})
			return
		}
		vf.Delete()
	}
	code := ctx.PostValue("code")
	emailtk, ok := kuser.VerifyMailCode(uid, code)
	if !ok {
		err := kses.NewSession(uid, "verify_flag", "true").Save()
		if err != nil {
			ctx.JSON(iris.Map{
				"status": "error",
				"error": "SaveSessionError",
				"errorMessage": err.Error(),
			})
			return
		}
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "VerifyError",
			"errorMessage": "code error",
		})
		return
	}
	ctx.JSON(iris.Map{
		"status": "ok",
		"token": emailtk,
	})
}

func SendVerifyEmailApi(ctx iris.Context){
	uid := kses.GetCtxUuid(ctx)
	cid := kses.GetSession(uid, "captid")
	if cid != nil {
		cid.Delete()
	}
	captcode := ctx.PostValue("capt")
	if cid == nil || !kcapt.VerifyCaptcha(cid.Value, captcode) {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "CaptcodeError",
		})
		return
	}
	email := ctx.PostValue("email")
	if !reg_email.MatchString(email) {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "IllegalDataError",
			"errorMessage": "Email is illegal data",
		})
		return
	}
	err := kuser.SendVerifyMail(uid, email)
	if err != nil {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "SendMailError",
			"errorMessage": err.Error(),
		})
		return
	}
	ctx.JSON(iris.Map{"status": "ok"})
}

func UserRegisterApi(ctx iris.Context){
	email, ok := kuser.CheckUserEmailToken(ctx.PostValue("emailtk"))
	if !ok {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "RegisterError",
			"errorMessage": "Email not verify",
		})
		return
	}
	name := ctx.PostValue("name")
	if !reg_name.MatchString(name) {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "IllegalDataError",
			"errorMessage": "Username is illegal data",
		})
		return
	}
	password := ctx.PostValue("pwd")
	if !reg_pwd.MatchString(password) {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "IllegalDataError",
			"errorMessage": "Password is illegal data",
		})
		return
	}
	user := kuser.NewUser(name, password, email, "")
	err := user.InsertData()
	if err != nil {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "InsertUserError",
			"errorMessage": err.Error(),
		})
		return
	}
	ufile.CreateDir(ufile.JoinPathWithoutAbs(USER_DATA_PATH, user.Id.String()))
	ctx.JSON(iris.Map{"status": "ok"})
}

func UserRegcheckApi(ctx iris.Context){
	email := ctx.PostValue("email")
	if !reg_email.MatchString(email) {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "IllegalDataError",
			"errorMessage": "Email is illegal data",
		})
		return
	}
	if n, _ := kuser.USER_SQL_TABLE.Count(kpsql.WhereMap{{"email", "=", email, ""}}, 1); n == 1 {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "UserExistError",
			"errorMessage": "Email has been used",
		})
		return
	}
	name := ctx.PostValue("name")
	if !reg_name.MatchString(name) {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "IllegalDataError",
			"errorMessage": "Username is illegal data",
		})
		return
	}
	if n, _ := kuser.USER_SQL_TABLE.Count(kpsql.WhereMap{{"username", "=", name, ""}}, 1); n == 1 {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "UserExistError",
			"errorMessage": "Username has been used",
		})
		return
	}
	password := ctx.PostValue("pwd")
	if !reg_pwd.MatchString(password) {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "IllegalDataError",
			"errorMessage": "Password is illegal data",
		})
		return
	}
	ctx.JSON(iris.Map{"status": "ok"})
}

func GetCaptImgApi(ctx iris.Context){
	uid := kses.GetCtxUuid(ctx)
	var (
		captid string
		imgdt string
		err error
	)
	if cid := kses.GetSessionStr(uid, "captid"); cid != "" {
		kcapt.RemoveCaptcha(cid)
	}
	captid, imgdt, err = kcapt.NewCaptcha()
	if err != nil {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "MakeCaptError",
			"errorMessage": err.Error(),
		})
		return
	}
	err = kses.NewSession(uid, "captid", captid).Save()
	if err != nil {
		kcapt.RemoveCaptcha(captid)
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "SqlInsertError",
			"errorMessage": err.Error(),
		})
		return
	}
	ctx.JSON(iris.Map{
		"status": "ok",
		"data": imgdt,
	})
}

func InitApi(group iris.Party){
	apigp := group.Party("/api")
	apigp.Get("/head/{id:uuid}", UserHeadApi).ExcludeSitemap()
	apigp.Post("/login", UserLoginApi).ExcludeSitemap()
	apigp.Post("/logout", UserLogoutApi).ExcludeSitemap()
	apigp.Post("/register", UserRegisterApi).ExcludeSitemap()
	apigp.Post("/regcheck", UserRegcheckApi).ExcludeSitemap()
	apigp.Post("/verify/email", VerifyEmailApi).ExcludeSitemap()
	apigp.Post("/verify/email/send", SendVerifyEmailApi).ExcludeSitemap()
	apigp.Get("/captcha/image", GetCaptImgApi).ExcludeSitemap()
}
