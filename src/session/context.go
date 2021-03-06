
package kweb_session

import (
	time "time"

	uuid "github.com/google/uuid"
	jwt "github.com/KpnmServer/go-util/jwt"
	iris "github.com/kataras/iris/v12"
)

var JWT_ENCODER jwt.Encoder = jwt.NewAutoEncoder(
	jwt.NewFileEncoder(jwt.NewEncoder(nil), "keys/session_jwt.key"), 2048, 60 * 60 * 24 * 23)

var cookie_max_time time.Duration = time.Hour * 24 * 23

func GetCtxUuid(ctx iris.Context)(id uuid.UUID){
	if uidtk, err := JWT_ENCODER.Decode(ctx.GetCookie("sesid")); err == nil {
		if id, err = uuid.Parse(uidtk["v"].(string)); err == nil {
			if (int64)(uidtk["iat"].(float64)) <= time.Now().Unix() + 60 * 60 * 24 * 7 {
				nid := uuid.New()
				if _, err = ChangeUUID(id, nid); err == nil {
					ctx.SetCookieKV("sesid", JWT_ENCODER.Encode(jwt.SetOutdate(
						jwt.Json{"v": id.String()}, cookie_max_time)), iris.CookieExpires(cookie_max_time))
					id = nid
				}
			}
			return
		}
	}
	id = uuid.New()
	ctx.SetCookieKV("sesid", JWT_ENCODER.Encode(jwt.SetOutdate(
		jwt.Json{"v": id.String()}, cookie_max_time)), iris.CookieExpires(cookie_max_time))
	return
}

func DelCtxUuid(ctx iris.Context)(err error){
	defer ctx.RemoveCookie("sesid")
	var (
		uidtk jwt.Json
		id uuid.UUID
	)
	if uidtk, err = JWT_ENCODER.Decode(ctx.GetCookie("sesid")); err == nil {
		if id, err = uuid.Parse(uidtk["v"].(string)); err == nil {
			_, err = DelUserSessions(id)
			return
		}
	}
	return nil
}



