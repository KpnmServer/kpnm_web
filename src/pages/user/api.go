
package page_user

import (
	regexp  "regexp"

	iris "github.com/kataras/iris/v12"
	ufile "github.com/KpnmServer/go-util/file"
	kses "github.com/KpnmServer/kpnm_web/src/session"
)

var (
	reg_name  *regexp.Regexp = regexp.MustCompile(`^[A-Za-z_-][0-9A-Za-z_-]{3,31}$`)
	reg_pwd   *regexp.Regexp = regexp.MustCompile(`^[A-Za-z][0-9A-Za-z_-]{7,127}$`)
	reg_email *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`)
)

func UserHeadApi(ctx iris.Context){
	userid := ctx.Params().Get("id")
	path := ufile.JoinPathWithoutAbs(USER_DATA_PATH, userid, "head.png")
	ctx.ServeFile(path)
}

func VerifyEmailApi(ctx iris.Context){
	uid := kses.GetCtxUuid(ctx)
	code := ctx.PostValue("code")
	emailtk, ok := verifyMailCode(uid, code)
	if !ok {
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
	email := ctx.PostValue("email")
	if !reg_email.MatchString(email) {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "IllegalDataError",
			"errorMessage": "Email value is illegal data",
		})
		return
	}
	err := sendVerifyMail(uid, email)
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
	/*uid :=*/ kses.GetCtxUuid(ctx)
	email, ok := checkUserEmailToken(ctx.PostValue("emailtk"))
	if !ok {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "RegisterError",
			"errorMessage": "Email not verify",
		})
		return
	}
	name := ctx.PostValue("name")
	password := ctx.PostValue("pwd")
	description := ctx.PostValue("desc")
	err := InsertUserData(NewUser(name, password, email, description))
	if err != nil {
		ctx.JSON(iris.Map{
			"status": "error",
			"error": "RegisterError",
			"errorMessage": err.Error(),
		})
		return
	}
	ctx.JSON(iris.Map{"status": "ok"})
}

func UserLoginApi(ctx iris.Context){
	uid := kses.GetCtxUuid(ctx)
}

func InitApi(group iris.Party){
	group.Get("/api/head/{id:uuid}", UserHeadApi)
	group.Post("/api/register", UserRegisterApi)
	group.Post("/api/verify/email", VerifyEmailApi)
	group.Post("/api/verify/email/send", SendVerifyEmailApi)
}
