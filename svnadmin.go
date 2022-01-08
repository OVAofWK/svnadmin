package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"svn/ini"
	"svn/web"
	"time"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
)

// 读svnadmin配置文件
var CONFIG = ini.ReadConfYaml("conf/svnconf.yaml")
var BasicUser = map[string]string{CONFIG.Admin.User: CONFIG.Admin.Passwd}

type H map[string]string

func NewWeb() *gin.Engine {
	r := gin.Default()

	// 日志写盘
	if CONFIG.Admin.UseLog {
		if isExist(CONFIG.Admin.LogPath) {
			currentTime := time.Now().Format("2006-01-02-150405")
			oldLog := CONFIG.Admin.LogPath + "-" + currentTime
			os.Rename(CONFIG.Admin.LogPath, oldLog)
		}
		f, _ := os.Create(CONFIG.Admin.LogPath)
		gin.DefaultWriter = io.MultiWriter(f)

		r.Use(gin.Recovery())
		r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		}))
	}

	// 加载静态资源
	fsCss := assetfs.AssetFS{Asset: web.Asset, AssetDir: web.AssetDir, AssetInfo: web.AssetInfo, Prefix: "web/css"}
	fsJs := assetfs.AssetFS{Asset: web.Asset, AssetDir: web.AssetDir, AssetInfo: web.AssetInfo, Prefix: "web/js"}
	fs := assetfs.AssetFS{Asset: web.Asset, AssetDir: web.AssetDir, AssetInfo: web.AssetInfo, Prefix: "web"}
	r.StaticFS("/css", &fsCss)
	r.StaticFS("/js", &fsJs)
	r.StaticFS("/favicon.ico", &fs)
	return r
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

// 响应get请求
func Getconf(c *gin.Context) {
	svnParam := c.Param("param_1")
	info := ini.ReadConf(CONFIG, svnParam)
	path := "web/" + svnParam + ".html"
	c.Writer.WriteHeader(200)
	adminHtml, _ := web.Asset(path)
	// 将html中的{{}}内的信息替换为变量内容
	backupFileList := ini.GetBackupsFileList()
	html := Renderer(adminHtml, H{"title": CONFIG.Web.Title, "info": info, "backupsFileList": backupFileList})
	c.Writer.Write(html)
}

// 响应post请求
func Postconf(c *gin.Context) {
	svnParam := c.Param("param_1")
	info := c.PostForm("info")
	err := ini.WriteConf(CONFIG, svnParam, info)
	if err == nil {
		path := "/admin/" + svnParam
		//重定向至提交的当前页
		c.Redirect(http.StatusMovedPermanently, path)
	} else {
		c.String(http.StatusOK, fmt.Sprint(err))
	}

}

func GetBackups(c *gin.Context) {
	param2 := c.Param("param_2")
	info := ini.ReadFile("backups/" + param2)
	path := "web/backups.html"
	c.Writer.WriteHeader(200)
	adminHtml, _ := web.Asset(path)
	backupFileList := ini.GetBackupsFileList()
	html := Renderer(adminHtml, H{"title": CONFIG.Web.Title, "info": info, "backupsFileList": backupFileList, "fileName": param2})
	c.Writer.Write(html)
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
		admin.GET("/:param_1/", Getconf)
		admin.POST("/:param_1/", Postconf)
		admin.GET("/backups/:param_2/", GetBackups)
	}
	return r
}

func main() {
	if !isExist("conf/svnconf.yaml") {
		fmt.Println("未找到svnadmin配置文件 conf/svnconf.yaml")
	}
	r := NewWeb()
	r = Route(r)
	fmt.Println(CONFIG.Server.Listen)
	r.Run(CONFIG.Server.Listen)

}
