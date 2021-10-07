
package kweb_manager

import (
	time "time"
	fmt "fmt"
	iris "github.com/kataras/iris/v12"
)

func BindSiteMap(app *iris.Application, prefix string, url_ ...string){
	url := "/site-map.xml"
	if len(url_) > 0 {
		url = url_[0]
	}
	var sitemapstr string
	sitemapstr = `<?xml version="1.0" encoding="utf-8" standalone="yes"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">`
	{
		for _, r := range app.GetRoutes() {
			if !r.IsStatic() || r.Subdomain != "" || !r.IsOnline() || r.NoSitemap || r.Method != "GET" {
				continue
			}
			path := r.StaticPath()
			cont := fmt.Sprintf(`<url>
	<loc>%s</loc>
	<lastmod>%s</lastmod>
</url>`, prefix + path, time.Now().Format("2006-01-02T15:04:05-07:00"))
			sitemapstr += cont
		}
	}
	sitemapstr += `</urlset>`
	app.Logger().Debugf("========sitemap========\n%s\n========sitemap========", sitemapstr)
	app.Get(url, func(ctx iris.Context){
		ctx.WriteString(sitemapstr)
	})
}
