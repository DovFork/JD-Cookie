package web

import (
	"crypto/tls"
	"embed"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"github.com/scjtqs/jd_cookie/config"
	"github.com/scjtqs/jd_cookie/web/repo"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"go.uber.org/dig"
	"html/template"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"time"
)

type httpServer struct {
	engine      *gin.Engine
	HTTP        *http.Server
	ct          *dig.Container
	Conf        *config.Conf
	cookiesRepo repo.CookiesRepository
}

type jumpLogin struct {
	Tk        *Token
	CookieJar *cookiejar.Jar
}

var HTTPServer = &httpServer{}

var ckChan = make(chan jumpLogin, 10)

func (s *httpServer) Run(addr string, ct *dig.Container) {
	s.ct = ct
	ct.Invoke(func(conf *config.Conf) {
		s.Conf = conf
	})
	var f embed.FS
	ct.Invoke(func(file embed.FS) {
		f = file
	})
	gin.SetMode(gin.ReleaseMode)
	s.engine = gin.New()
	// 创建基于 内存 的存储引擎，secret 参数是用于加密的密钥
	store := memstore.NewStore([]byte("scjtqsnb"))
	// 设置session中间件，参数mysession，指的是session的名字，也是cookie的名字
	// store是前面创建的存储引擎，我们可以替换成其他存储引擎
	s.engine.Use(sessions.Sessions("mysession", store))

	s.engine.Use(func(c *gin.Context) {
		if c.Request.Method != "GET" && c.Request.Method != "POST" {
			log.Warnf("已拒绝客户端 %v 的请求: 方法错误", c.Request.RemoteAddr)
			c.Status(404)
			return
		}
		c.Next()
	})
	// 自动加载模板
	t := template.New("tmp")
	//func 函数映射 全局模板可用
	t.Funcs(template.FuncMap{
		"getYear":        GetYear,
		"formatAsDate":   FormatAsDate,
		"getDate":        GetDate,
		"getavator":      Getavator,
		"getServerInfo":  GetServerInfo,
		"formatFileSize": FormatFileSize,
	})

	//从二进制中加载模板（后缀必须.html)
	templ := template.Must(template.New("").ParseFS(f, "template/html/*.html"))
	s.engine.SetHTMLTemplate(templ)
	//静态资源
	//s.engine.Static("/assets", "./template/assets")
	//s.engine.StaticFS("/public", http.FS(f))
	s.engine.GET("/", func(c *gin.Context) {
		s.GetclientIP(c)
		var v string
		ct.Invoke(func(version string) {
			v = version
		})
		c.HTML(http.StatusOK, "upcookie.html", gin.H{
			"version": v,
		})
	})
	// 静态文件处理
	s.engine.GET("assets/*action", func(c *gin.Context) {
		c.FileFromFS("template/assets/"+c.Param("action"), http.FS(f))
	})

	// 路由
	// 获取二维码
	s.engine.GET("/qrcode", s.getQrcode)
	s.engine.GET("/qrcode_jumplogin", s.getQrcode_jumplogin)
	// 获取返回的cookie信息
	s.engine.GET("/cookie", s.getCookie)
	// 获取各种配置文件api
	s.engine.GET("/api/config/:key")
	// 保存配置
	s.engine.POST("/api/upsave", s.upsave)
	s.engine.POST("/api/save")
	s.engine.GET("/home")
	s.engine.POST("/auth")
	//s.engine.GET("/test",s.test)

	// 初始化db
	s.initdb()

	go func() {
		log.Infof("jdcookie提取 服务器已启动: %v", addr)
		log.Info("请用浏览器打开url: http://公网ip或者域名%s", addr)
		log.Warn("请务必使用公网访问，否则读取到的客户端Ip会是内网Ip，不是公网Ip.")
		log.Warnf("v3.x 版本 是服务端部署版本。客户端需要使用浏览器打开，让浏览器和手机在同一个网络下（或者直接用手机打开浏览器）")
		s.HTTP = &http.Server{
			Addr:    addr,
			Handler: s.engine,
		}
		if err := s.HTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(err)
			log.Infof("HTTP 服务启动失败, 请检查端口是否被占用.")
			log.Warnf("将在五秒后退出.")
			time.Sleep(time.Second * 5)
			os.Exit(1)
		}
	}()
	go s.backgroundRun()
}

func (s *httpServer) initdb() {
	if s.Conf.DbConf.DbEnable {
		var err error
		err = repo.InitRDBMS(s.Conf.DbConf)
		if err != nil {
			log.Fatalf("faild to init db error= %s", err.Error())
		}
		s.cookiesRepo, err = repo.NewCookieRepo()
		if err != nil {
			log.Fatalf("faild to get initd db error= %s", err.Error())
		}
		s.cookiesRepo.InitTables()
	}
}

//直接唤起 京东 客户端 后台获取cookie
func (s *httpServer) backgroundRun() {
	tk := <-ckChan
	//jar := s.getCookieJar(c)
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	getUrl := "https://plogin.m.jd.com/cgi-bin/m/tmauthchecktoken?&token=" + tk.Tk.Token + "&ou_state=0&okl_token=" + tk.Tk.Okl_token
	client := &http.Client{
		Jar:       tk.CookieJar,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	var res string
	tm := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	ua := fmt.Sprintf("Mozilla/5.0 (iPhone; CPU iPhone OS 13_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 SP-engine/2.14.0 main/1.0 baiduboxapp/11.18.0.16 (Baidu; P2 13.3.1) NABar/0.0 TM/%s", tm)
	for i := 1; i <= 100; i++ {
		err := gout.New(client).
			POST(getUrl).
			//Debug(true).
			SetWWWForm(
				gout.H{
					"lang":      "chs",
					"appid":     300,
					"returnurl": "https://wqlogin2.jd.com/passport/LoginRedirect?state=1100399130787&returnurl=//home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&/myJd/home.action",
					"source":    "wq_passport",
				},
			).
			BindBody(&res).
			SetHeader(gout.H{
				"Referer":      "https://plogin.m.jd.com/login/login?appid=300&returnurl=https://wqlogin2.jd.com/passport/LoginRedirect?state=" + timeStamp + "&returnurl=//home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&/myJd/home.action&source=wq_passport",
				"Cookie":       tk.Tk.Cookies,
				"Connection":   "Keep-Alive",
				"Content-Type": "application/x-www-form-urlencoded; Charset=UTF-8",
				"Accept":       "application/json, text/plain, */*",
				//"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36",
				"User-Agent": ua,
			}).
			SetTimeout(timeout).
			F().Retry().Attempt(5).
			WaitTime(time.Millisecond * 500).MaxWaitTime(time.Second * 5).
			Do()
		if err != nil {
			return
		}
		checkJson := gjson.Parse(res)
		if checkJson.Get("errcode").Int() == 0 {
			//获取cookie
			u, _ := url.Parse("https://plogin.m.jd.com")
			token := tk.Tk
			var TrackerID, pt_key, pt_pin, pt_token, pwdt_id, s_key, s_pin = "", "", "", "", "", "", ""
			for _, v := range tk.CookieJar.Cookies(u) {
				if v.Name == "TrackerID" {
					TrackerID = v.Value
				}
				if v.Name == "pt_key" {
					pt_key = v.Value
				}
				if v.Name == "pt_pin" {
					pt_pin = v.Value
				}
				if v.Name == "pt_token" {
					pt_token = v.Value
				}
				if v.Name == "pwdt_id" {
					pwdt_id = v.Value
				}
				if v.Name == "s_key" {
					s_key = v.Value
				}
				if v.Name == "s_pin" {
					s_pin = v.Value
				}
			}
			token.Cookies = "TrackerID=" + TrackerID + "; pt_key=" + pt_key + "; pt_pin=" + pt_pin + "; pt_token=" + pt_token + "; pwdt_id=" + pwdt_id + "; s_key=" + s_key + "; s_pin=" + s_pin + "; wq_skey="
			token.UserCookie = "pt_key=" + pt_key + ";pt_pin=" + pt_pin + ";"
			token.PtPin = pt_pin
			token.PtKey = pt_key
			log.Info("############  登录成功，获取到 Cookie  #############")
			log.Infof("Cookie1=%s", token.UserCookie)
			log.Info("####################################################")
			//写db
			if s.Conf.DbConf.DbEnable {
				_, err := s.cookiesRepo.UpdateCookie(token.PtPin, token.PtKey, token.UserCookie)
				if err != nil {
					log.Errorf("save cookie to db faild %s", err.Error())
				}
			}
			////发送数据给 挂机服务器
			postUrl := s.Conf.UpSave
			if postUrl != "" {
				var res MSG
				code := 0
				err := gout.POST(postUrl).
					//Debug(true).
					SetWWWForm(
						gout.H{
							"userCookie": token.UserCookie,
						},
					).
					BindJSON(&res).
					SetHeader(gout.H{
						"Connection":   "Keep-Alive",
						"Content-Type": "application/x-www-form-urlencoded; Charset=UTF-8",
						"Accept":       "application/json, text/plain, */*",
						"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36",
					}).
					Code(&code).
					SetTimeout(timeout).
					F().Retry().Attempt(5).
					WaitTime(time.Millisecond * 500).MaxWaitTime(time.Second * 5).
					Do()
				if err != nil || code != 200 {
					log.Errorf("upsave notify post  usercookie to %s faild", postUrl)
				} else {
					errcode := res["err"]
					if errcode == nil {
						errcode = 0
					}
					title := res["title"]
					if title == nil {
						title = "更新到挂机服务成功"
					}
					msg := res["msg"]
					if msg == nil {
						msg = "cookie= " + token.UserCookie
					}
					log.Infof("errcode=%d,title=%s,msg=%s", errcode, title, msg)
				}
			}
			break
		} else {
			log.Errorf("获取cookie失败")
		}
		time.Sleep(time.Second * 1)
	}
}
