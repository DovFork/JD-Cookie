package main

import (
	"flag"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/scjtqs/jd_cookie/config"
	"github.com/scjtqs/jd_cookie/util"
	"github.com/scjtqs/jd_cookie/web"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"os"
	"os/signal"
	"path"
	"time"
	"go.uber.org/dig"
)

var (
	h       bool
	d       bool
	Version = "v2.0.0"
	Build   string
	configPath = "config.json"
)

func init() {
	var debug bool
	flag.BoolVar(&d, "d", false, "running as a daemon")
	flag.BoolVar(&debug, "D", false, "debug mode")
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&configPath,"c","config.json","config file path default is config.json")
	flag.Parse()
	logFormatter := &easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[%time%] [%lvl%]: %msg% \n",
	}
	w, err := rotatelogs.New(path.Join("logs", "%Y-%m-%d.log"), rotatelogs.WithRotationTime(time.Hour*24))
	if err != nil {
		log.Errorf("rotatelogs init err: %v", err)
		panic(err)
	}
	LogLevel := "info"
	if debug {
		log.SetReportCaller(true)
		LogLevel = "debug"
	}
	log.AddHook(util.NewLocalHook(w, logFormatter, util.GetLogLevel(LogLevel)...))
}

func main() {
	if h {
		help()
	}
	if d {
		web.Daemon()
	}
	conf:=config.GetConfigFronPath(configPath)
	container := dig.New()
	container.Provide(func() *config.Conf {
		return conf
	})
	conf.Save(configPath)
	log.Infof("欢迎使用jdcookie提取器 by scjtqs %s", Version)
	log.Info("当前开源版本：获取cookie成功后，不会自动提交到挂机服务器，需要自行修改了")
	web.HTTPServer.Run(":29099",container)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}

// help cli命令行-h的帮助提示
func help() {
	log.Infof(`jd_cookie service
version: %s

Usage:

server [OPTIONS]

Options:
`, Version, Build)
	flag.PrintDefaults()
	os.Exit(0)
}
