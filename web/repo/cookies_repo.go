package repo

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type Cookies struct {
	Id         int64     `xorm:"bigint(20) pk notnull autoincr 'id'"`
	PtPin      string    `xorm:"notnull unique 'pt_pin'"`
	UserCookie string    `xorm:"text notnull 'user_cookie'"`
	PtKey      string    `xorm:"notnull 'pt_key'"`
	CreateTime time.Time `xorm:"datetime updated 'createtime'"`
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
//var sql = "CREATE TABLE IF NOT EXISTS `cookies` (`id` bigint(20) NOT NULL AUTO_INCREMENT,`pt_pin` varchar(100) NOT NULL,`user_cookie` text NOT NULL COMMENT '用户cookie', `pt_key` varchar(255) NOT NULL,`createtime` datetime DEFAULT CURRENT_TIMESTAMP,PRIMARY KEY (`id`),UNIQUE KEY (`pt_pin`))ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4"

func (cksp *CookiesRepo) InitTables() {
	ext, err := cksp.db.DB().IsTableExist("cookies")
	if err != nil {
		panic("db error: " + err.Error())
	}
	if !ext {
		err := cksp.db.DB().CreateTables(Cookies{})
		if err != nil {
			log.Errorf("faild to init db err=%v", err)
		}
	}
}

func (cksp *CookiesRepo) UpdateCookie(pt_pin, pt_key, usercookie string) (*Cookies, error) {
	var cks Cookies
	has, err := cksp.db.DB().Where("pt_pin = ?", pt_pin).Get(&cks)
	cksNew := &Cookies{
		PtPin:      pt_pin,
		PtKey:      pt_key,
		UserCookie: usercookie,
	}
	if err != nil {
		return nil, err
	}
	if !has {
		log.Warnf("%s record not exits", pt_pin)
		_, err = cksp.db.DB().Insert(cksNew)
		return cksNew, err
	}
	cks.UserCookie=usercookie
	cks.PtKey=pt_key
	cks.PtPin=pt_pin
	_, err = cksp.db.DB().Id(cks.Id).Update(cksNew)
	return &cks, err
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
