package web

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"time"
)

var timeout = time.Second * 5

var ua = "jdapp;android;10.0.5;11;0393465333165363-5333430323261366;network/wifi;model/M2102K1C;osVer/30;appBuild/88681;partner/lc001;eufv/1;jdSupportDarkMode/0;Mozilla/5.0 (Linux; Android 11; M2102K1C Build/RKQ1.201112.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045534 Mobile Safari/537.36"

// 获取二维码
func (s *httpServer) getQrcode(c *gin.Context) {
	s.GetclientIP(c)
	log.Warn("start get qrcode")
	_, err := s.step1(c)
	if err != nil {
		c.JSON(200, MSG{
			"err": 1,
			"msg": err,
		})
		return
	}
	qrurl, err := s.setp2(c)
	if err != nil {
		c.JSON(200, MSG{
			"err": 1,
			"msg": err,
		})
		return
	}
	log.Warnf("get qrcode url = %s", qrurl)
	c.JSON(200, MSG{
		"err":    0,
		"qrcode": qrurl,
	})
}

func (s *httpServer) praseSetCookies(c *gin.Context, rsp string, cookie *cookiejar.Jar) {
	json := gjson.Parse(rsp)
	token := s.getToken(c)
	token.Stoken = json.Get("s_token").String()
	u, _ := url.Parse("https://plogin.m.jd.com")
	a := cookie.Cookies(u)
	for _, v := range a {
		if v.Name == "guid" {
			token.Guid = v.Value
		}
		if v.Name == "lsid" {
			token.Lsid = v.Value
		}
		if v.Name == "lstoken" {
			token.Lstoken = v.Value
		}
	}
	token.Cookies = "guid=" + token.Guid + "; lang=chs; lsid=" + token.Lsid + "; lstoken=" + token.Lstoken + "; "
	s.updateToken(c, token)
	//log.Warnf("cookies=%s", cookies)
}

// 获取二维码第一步
func (s *httpServer) step1(c *gin.Context) (*cookiejar.Jar, error) {
	ip := s.GetclientIP(c)
	jar := s.getCookieJar(c)
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	getUrl := "https://plogin.m.jd.com/cgi-bin/mm/new_login_entrance?lang=chs&appid=300&returnurl=https://wq.jd.com/passport/LoginRedirect?state=" + timeStamp + "&returnurl=https://home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&/myJd/home.action&source=wq_passport"
	client := &http.Client{
		// example, custom cookie jar implements
		Jar: jar,
		// example, ignore self-signed certificate
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	method := "GET"
	req, err := http.NewRequest(method, getUrl, nil)

	if err != nil {
		log.Errorf("get qrcode step1 faild err=%s", err.Error())
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Language", "zh-cn")
	req.Header.Add("Referer", "https://plogin.m.jd.com/cgi-bin/mm/new_login_entrance?lang=chs&appid=300&returnurl=https://wq.jd.com/passport/LoginRedirect?state="+timeStamp+"&returnurl=https://home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&/myJd/home.action&source=wq_passport")
	//req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36")
	req.Header.Add("User-Agent", ua)
	req.Header.Add("Host", "plogin.m.jd.com")
	req.Header.Set("X-Forwarded-For", ip)
	req.Header.Set("Proxy-Client-IP", ip)
	req.Header.Set("WL-Proxy-Client-IP", ip)
	req.Header.Set("CLIENT-IP", ip)
	res, err := client.Do(req)
	if err != nil {
		log.Errorf("get qrcode step1 faild err=%s", err.Error())
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("get qrcode step1 faild err=%s", err.Error())
		return nil, err
	}
	s.praseSetCookies(c, string(body), jar)
	//log.Warnf("url=%s,cookiejar=%+v,res=%s", getUrl, jar, string(body))
	s.updateCookieJar(c, jar)
	return jar, nil
}

// 获取 二维码第二步
func (s *httpServer) setp2(c *gin.Context) (string, error) {
	ip := s.GetclientIP(c)
	token := s.getToken(c)
	jar := s.getCookieJar(c)
	if token.Cookies == "" {
		return "", errors.New("empty cookies")
	}
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	getUrl := "https://plogin.m.jd.com/cgi-bin/m/tmauthreflogurl?s_token=" + token.Stoken + "&v=" + timeStamp + "&remember=true"
	client := &http.Client{
		Jar:       jar,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	var res string
	err := gout.New(client).
		POST(getUrl).
		//Debug(true).
		SetJSON(
			gout.H{
				"lang":      "chs",
				"appid":     300,
				"returnurl": "https://wqlogin2.jd.com/passport/LoginRedirect?state=" + timeStamp + "&returnurl=//home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&/myJd/home.action",
				"source":    "wq_passport",
			},
		).
		BindBody(&res).
		SetHeader(gout.H{
			"Connection":   "Keep-Alive",
			"Content-Type": "application/x-www-form-urlencoded; Charset=UTF-8",
			"Accept":       "application/json, text/plain, */*",
			"Cookie":       token.Cookies,
			"Referer":      "https://plogin.m.jd.com/login/login?appid=300&returnurl=https://wqlogin2.jd.com/passport/LoginRedirect?state=" + timeStamp + "&returnurl=//home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&/myJd/home.action&source=wq_passport",
			//"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36",
			"User-Agent":         ua,
			"X-Forwarded-For":    ip,
			"Proxy-Client-IP":    ip,
			"WL-Proxy-Client-IP": ip,
			"CLIENT-IP":          ip,
		}).
		SetTimeout(timeout).
		F().Retry().Attempt(5).
		WaitTime(time.Millisecond * 500).MaxWaitTime(time.Second * 5).
		Do()
	s.updateCookieJar(c, jar)
	if err != nil {
		return "", err
	}
	resjson := gjson.Parse(res)
	token.Token = resjson.Get("token").String()
	u, _ := url.Parse("https://plogin.m.jd.com")
	for _, v := range jar.Cookies(u) {
		if v.Name == "okl_token" {
			token.Okl_token = v.Value
		}
	}
	qrUrl := "https://plogin.m.jd.com/cgi-bin/m/tmauth?appid=300&client_type=m&token=" + token.Token
	s.updateToken(c, token)
	//log.Warnf("url=%s,cookiejar=%+v,res=%s", getUrl, jar, res)
	return qrUrl, nil
}

// 获取返回的cookie信息
func (s *httpServer) getCookie(c *gin.Context) {
	//session := sessions.Default(c)
	//cookies=session.Get("cookies").(string)
	jar := s.getCookieJar(c)
	check, err := s.checkLogin(c, jar)
	if err != nil {
		c.JSON(200, MSG{
			"err": 1,
			"msg": err,
		})
		return
	}
	checkJson := gjson.Parse(check)
	if checkJson.Get("errcode").Int() == 0 {
		//获取cookie
		ucookie := s.getJdCookie(check, jar, c)
		c.JSON(200, MSG{
			"err":    0,
			"cookie": ucookie,
		})
		return
	} else {
		c.JSON(200, MSG{
			"err": checkJson.Get("errcode").Int(),
			"msg": checkJson.Get("message").String(),
		})
	}

}

// 校验登录状态
func (s *httpServer) checkLogin(c *gin.Context, jar *cookiejar.Jar) (string, error) {
	token := s.getToken(c)
	ip := s.GetclientIP(c)
	if token.Cookies == "" {
		return "", errors.New("empty cookies")
	}
	//jar := s.getCookieJar(c)
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	getUrl := "https://plogin.m.jd.com/cgi-bin/m/tmauthchecktoken?&token=" + token.Token + "&ou_state=0&okl_token=" + token.Okl_token
	client := &http.Client{
		Jar:       jar,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	var res string
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
			"Cookie":       token.Cookies,
			"Connection":   "Keep-Alive",
			"Content-Type": "application/x-www-form-urlencoded; Charset=UTF-8",
			"Accept":       "application/json, text/plain, */*",
			//"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36",
			"User-Agent":         ua,
			"X-Forwarded-For":    ip,
			"Proxy-Client-IP":    ip,
			"WL-Proxy-Client-IP": ip,
			"CLIENT-IP":          ip,
		}).
		SetTimeout(timeout).
		F().Retry().Attempt(5).
		WaitTime(time.Millisecond * 500).MaxWaitTime(time.Second * 5).
		Do()
	s.updateCookieJar(c, jar)
	if err != nil {
		return "", err
	}
	//log.Warnf("checkLogin res=%s",res)
	return res, nil
}

// 解析用户的cookie
func (s *httpServer) getJdCookie(resp string, cookie *cookiejar.Jar, c *gin.Context) string {
	u, _ := url.Parse("https://plogin.m.jd.com")
	token := s.getToken(c)
	var TrackerID, pt_key, pt_pin, pt_token, pwdt_id, s_key, s_pin = "", "", "", "", "", "", ""
	for _, v := range cookie.Cookies(u) {
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
	s.updateToken(c, token)
	log.Info("############  登录成功，获取到 Cookie  #############")
	log.Infof("Cookie1=%s", token.UserCookie)
	log.Info("####################################################")
	return token.UserCookie
}

func (s *httpServer) upsave(c *gin.Context) {
	//log.Warnf("更新到挂机服务器 res=%v", res)
	// 清空缓存参数
	token := s.getToken(c)
	s.cleanSession(c)
	////发送数据给 挂机服务器
	postUrl := os.Getenv("UPSAVE")
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
			c.JSON(200, MSG{
				"err":   1,
				"title": "更新到挂机服务器失败",
				"msg":   err.Error(),
			})
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
			c.JSON(200, MSG{
				"err":   errcode,
				"title": title,
				"msg":   fmt.Sprintf("%s, cookie= %s", msg, token.UserCookie),
			})
		}
		return
	}

	c.JSON(200, MSG{
		"err":   0,
		"title": "提取cookie成功",
		"msg":   "cookie= " + token.UserCookie,
	})
}
