
package kweb_user

import (
	time "time"
	http "net/http"
	crand "crypto/rand"
	hex "encoding/hex"

	uuid "github.com/google/uuid"
	jwt "github.com/KpnmServer/go-util/jwt"
	iris "github.com/kataras/iris/v12"
	kses "github.com/KpnmServer/kpnm_web/src/session"
)

const (
	LOG_USER_KEY = "kweb.user"
)

const LOG_LIVE_MAX_TIME time.Duration = time.Hour * 24 * 23

func GetCtxLog(ctx iris.Context)(user *UserData){
	if ctx.Values().Exists(LOG_USER_KEY) {
		return ctx.Values().Get(LOG_USER_KEY).(*UserData)
	}
	uid := kses.GetCtxUuid(ctx)
	cmt := LOG_LIVE_MAX_TIME
	if logtk, err := kses.JWT_ENCODER.Decode(ctx.GetCookie("loguser", func(ctx iris.Context, c *http.Cookie, _ uint8){
		if c.MaxAge == 0 { cmt = 0 }
	})); err == nil {
		if userid, err := uuid.Parse(logtk["v"].(string)); err == nil {
			rtokenses := kses.GetSession(uid, "loguser")
			if rtokenses != nil && rtokenses.Value == logtk["a"].(string) {
				user = GetUserData(userid)
				if rtokenses.Overtime.Unix() <= time.Now().Unix() + 60 * 60 * 24 * 14 {
					user.SaveCtxLog(ctx, cmt)
				}
				return
			}
		}
	}
	return nil
}

func (user *UserData)SaveCtxLog(ctx iris.Context, live time.Duration)(err error){
	uid := kses.GetCtxUuid(ctx)
	rbts := make([]byte, 16)
	_, err = crand.Read(rbts)
	if err != nil { return }
	rtoken := hex.EncodeToString(rbts)
	err = kses.NewSession(uid, "loguser", rtoken, LOG_LIVE_MAX_TIME).Save()
	if err != nil { return }
	ctx.SetCookieKV("loguser", kses.JWT_ENCODER.Encode(jwt.SetOutdate(
		jwt.Json{"v": user.Id.String(), "a": rtoken}, LOG_LIVE_MAX_TIME)),
		func(_ iris.Context, c *http.Cookie, _ uint8){ c.MaxAge = int(live.Seconds())})
	return nil
}

func RemoveCtxLog(ctx iris.Context){
	ctx.RemoveCookie("loguser")
}

func LogOrDoHandler(handler iris.Handler)(iris.Handler){
	return func(ctx iris.Context){
		user := GetCtxLog(ctx)
		if user == nil {
			handler(ctx)
			return
		}
		ctx.Values().Set(LOG_USER_KEY, user)
		ctx.Next()
	}
}

func LogOrRedirectHandler(path string)(iris.Handler){
	return LogOrDoHandler(func(ctx iris.Context){
		ctx.Redirect(path, iris.StatusFound)
	})
}

var MustLogHandler iris.Handler = LogOrDoHandler(func(ctx iris.Context){
	ctx.StatusCode(iris.StatusUnauthorized)
})
