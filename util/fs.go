package util

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

// WriteAllText 将给定text写入给定path
func WriteAllText(path, text string) error {
	return ioutil.WriteFile(path, []byte(text), 0o644)
}

// GetCurrentPath 获取当前文件的路径，直接返回string
func GetCurrentPath() string {
	cwd, e := os.Getwd()
	if e != nil {
		panic(e)
	}
	return cwd
}

// PathExists 判断给定path是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// ReadAllText 读取给定path对应文件，无法读取时返回空值
func ReadAllText(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error(err)
		return ""
	}
	return string(b)
}