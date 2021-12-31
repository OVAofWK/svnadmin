package main

import (
	"fmt"
	"net/http"
	"regexp"
	"svn/ini"
	"svn/web"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
)

var CONFIG = ini.ReadConfYaml("conf/svnconf.yaml")
var BasicUser = map[string]string{CONFIG.Admin.User: CONFIG.Admin.Passwd}

type H map[string]string

func NewWeb() *gin.Engine {
	r := gin.Default()
	fsCss := assetfs.AssetFS{Asset: web.Asset, AssetDir: web.AssetDir, AssetInfo: web.AssetInfo, Prefix: "web/css"}
	fsJs := assetfs.AssetFS{Asset: web.Asset, AssetDir: web.AssetDir, AssetInfo: web.AssetInfo, Prefix: "web/js"}
	fs := assetfs.AssetFS{Asset: web.Asset, AssetDir: web.AssetDir, AssetInfo: web.AssetInfo, Prefix: "web"}
	r.StaticFS("/css", &fsCss)
	r.StaticFS("/js", &fsJs)
	r.StaticFS("/favicon.ico", &fs)
	return r
}

func Getconf(c *gin.Context) {
	svnParam := c.Param("svn_param")
	info := ini.ReadConf(CONFIG, svnParam)
	path := "web/" + svnParam + ".html"

	c.Writer.WriteHeader(200)
	adminHtml, _ := web.Asset(path)
	// 将html中的{{}}内的信息替换为变量内容
	html := Renderer(adminHtml, H{"title": CONFIG.Web.Title, "info": info})
	c.Writer.Write(html)
}

func Postconf(c *gin.Context) {
	svnParam := c.Param("svn_param")
	info := c.PostForm("info")
	err := ini.WriteConf(CONFIG, svnParam, info)
	if err == nil {
		path := "/admin/" + svnParam
		// c.HTML(http.StatusOK, path, gin.H{"title": CONFIG.Web.Title, "info": info})
		//重定向至提交的当前页
		c.Redirect(http.StatusMovedPermanently, path)
	} else {
		c.String(http.StatusOK, fmt.Sprint(err))
	}

}

// 自定义渲染HTML模板
func Renderer(html []byte, p map[string]string) []byte {
	shtml := string(html)

	for k, v := range p {
		param := "{{\\s*." + k + "\\s*}}"
		reg := regexp.MustCompile(param)
		shtml = reg.ReplaceAllString(shtml, v)
	}
	return []byte(shtml)
}

func Route(r *gin.Engine) *gin.Engine {
	r.GET("/", func(c *gin.Context) {
		c.Writer.WriteHeader(200)
		indexHtml, _ := web.Asset("web/index.html")
		html := Renderer(indexHtml, H{"title": CONFIG.Web.Title})
		c.Writer.Write(html)
		// c.Writer.Header().Add("Accept", "text/html")
		// c.Writer.Flush()
	})
	r.POST("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/admin")
	})
	r.GET("/admin", func(c *gin.Context) {
		c.Writer.WriteHeader(200)
		adminHtml, _ := web.Asset("web/admin.html")
		html := Renderer(adminHtml, H{"title": CONFIG.Web.Title})
		c.Writer.Write(html)
	})
	admin := r.Group("/admin", gin.BasicAuth(BasicUser))
	{
		admin.GET("/:svn_param/", Getconf)
		admin.POST("/:svn_param/", Postconf)
	}
	return r
}

func main() {
	r := NewWeb()
	r = Route(r)
	r.Run(CONFIG.Server.Listen)
}
