package web

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/vmihailenco/msgpack"
	"net/http/cookiejar"
)

// 将部分参数存于session中。拒绝全局变量

// 获取客户端的ip地址
func (s *httpServer) GetclientIP(c *gin.Context) string {
	session := sessions.Default(c)
	var ip string
	if session.Get("clientip") != nil {
		ip = session.Get("clientip").(string)
	}
	if ip == "" {
		ip = c.ClientIP()
		session.Set("clientip", ip)
		session.Save()
	}
	return ip
}

// 读取cookieJar
func (s *httpServer) getCookieJar(c *gin.Context) *cookiejar.Jar {
	jar, _ := cookiejar.New(nil)
	session := sessions.Default(c)
	if session.Get("cookieJar") != nil {
		msgpack.Unmarshal(session.Get("cookieJar").([]byte),jar)
	}
	return jar
}

// 写入cookieJar
func (s *httpServer) updateCookieJar(c *gin.Context, jar *cookiejar.Jar) error {
	session := sessions.Default(c)
	res,_:=msgpack.Marshal(jar)
	session.Set("cookieJar", res)
	return session.Save()
}

type Token struct {
	Stoken     string `json:"s_token"`
	Cookies    string `json:"cookies"`
	Guid       string `json:"guid"`
	Lsid       string `json:"lsid"`
	Lstoken    string `json:"lstoken"`
	Okl_token  string `json:"okl_token"`
	Token      string `json:"token"`
	UserCookie string `json:"user_cookie"`
}

func (s *httpServer) getToken(c *gin.Context) *Token {
	session := sessions.Default(c)
	token := &Token{}
	if session.Get("token") != nil {
		json.Unmarshal([]byte(session.Get("token").(string)),token)
	}
	return token
}

// 更新 token那一批变量
func (s *httpServer) updateToken(c *gin.Context, token *Token) (*Token, error) {
	session := sessions.Default(c)
	u := &Token{}
	if session.Get("token") != nil {
		json.Unmarshal([]byte(session.Get("token").(string)), u)
	}
	if token.Token != "" {
		u.Token = token.Token
	}
	if token.Cookies != "" {
		u.Cookies = token.Cookies
	}
	if token.Guid != "" {
		u.Guid = token.Guid
	}
	if token.Lsid != "" {
		u.Lsid = token.Lsid
	}
	if token.Lstoken != "" {
		u.Lstoken = token.Lstoken
	}
	if token.Okl_token != "" {
		u.Okl_token = token.Okl_token
	}
	if token.Stoken != "" {
		u.Stoken = token.Stoken
	}
	if token.UserCookie != "" {
		u.UserCookie = token.UserCookie
	}
	set, _ := json.Marshal(u)
	session.Set("token", string(set))
	err := session.Save()
	return u, err
}

func (s *httpServer) cleanSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
}
