package repo

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type Cookies struct {
	PtPin      string    `xorm:"unique varchar(100) notnull 'pt_pin'"`
	UserCookie string    `xorm:"text notnull 'user_cookie'"`
	PtKey      string    `xorm:"varchar(255) notnull 'pt_key'"`
	CreateTime time.Time `xorm:"datetime 'createtime'"`
}

// TableName 数据库名称
func (s *Cookies) TableName() string {
	return "cookies"
}

type CookiesRepo struct {
	db *rdbms
}

// NewCookieRepo 创建cookies仓库
func NewCookieRepo() (*CookiesRepo, error) {
	var cks CookiesRepo
	var err error
	cks.db, err = getRDBMSInstance()
	if err != nil {
		return nil, err
	}

	return &cks, nil
}

// 初始化db
var sql = "CREATE TABLE IF NOT EXISTS `cookies` (`pt_pin` varchar(100) NOT NULL,`user_cookie` text NOT NULL COMMENT '用户cookie', `pt_key` varchar(255) NOT NULL,`createtime` datetime DEFAULT CURRENT_TIMESTAMP,PRIMARY KEY (`pt_pin`))ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4"

func (cksp *CookiesRepo) InitTables() {
	ext, err := cksp.db.DB().IsTableExist("cookies")
	if err != nil {
		panic("db error: " + err.Error())
	}
	if !ext {
		res, err := cksp.db.DB().Exec(sql)
		if err != nil {
			log.Errorf("faild to init db err=%v", err)
		}
		log.Infof("init table success res= %v", res)
	}
}

func (cksp *CookiesRepo) UpdateCookie(pt_pin, pt_key, usercookie string) (*Cookies, error) {
	var cks Cookies
	createTime := time.Unix(time.Now().Unix(), 0)
	_, err := cksp.db.DB().Exec("REPLACE INTO cookies (`pt_pin`,`user_cookie`,`pt_key`,`createtime`) VALUES(?,?,?,?)", pt_pin, usercookie, pt_key, createTime)
	if err != nil {
		return nil, err
	}
	cks.PtPin = pt_pin
	cks.PtKey = pt_key
	cks.UserCookie = usercookie
	cks.CreateTime = createTime
	return &cks, nil
}

func (cksp *CookiesRepo) GetCookieByPtPin(pt_pin string) (*Cookies, error) {
	var cks Cookies
	has, err := cksp.db.DB().Where("pt_pin = ?", pt_pin).Get(&cks)
	if err != nil {
		return nil, err
	}
	if !has {
		log.Warnf("%s record not exits", pt_pin)
	}
	return &cks, nil
}

func (cksp *CookiesRepo) DeleteCookieByPtPin(pt_pin string) (int64, error) {
	var cks Cookies
	return cksp.db.DB().Where("pt_pin = ?", pt_pin).Delete(&cks)
}
