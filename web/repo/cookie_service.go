package repo

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/scjtqs/jd_cookie/config"
	log "github.com/sirupsen/logrus"
)

type rdbms struct {
	db *xorm.EngineGroup
}

var rdbmsInst *rdbms

func (r *rdbms) DB() *xorm.EngineGroup {
	return r.db
}

func (r *rdbms) Transaction(f func(*xorm.Session) (interface{}, error)) (interface{}, error) {
	return r.db.Transaction(f)
}

func InitRDBMS(ds config.DbConf) error {
	var dsn []string = make([]string, 0)
	var err error

	rdbmsInst = &rdbms{}

	if ds.DbEnable {
		d := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", ds.DbUser, ds.DbPass, ds.DbHost, ds.DbPort, ds.DbDatabase)
		log.Info(d)
		dsn = append(dsn, d)
	}
	driver := ds.DbType
	if driver != "mysql" {
		driver = "mysql"
	}
	rdbmsInst.db, err = xorm.NewEngineGroup(driver, dsn)
	rdbmsInst.db.ShowSQL(true)
	if err != nil {
		log.Errorf("init db error = %s", err.Error())
		return err
	}

	return nil
}

func getRDBMSInstance() (*rdbms, error) {
	if rdbmsInst == nil {
		log.Errorf("persist storage was not been initialized.")
		return nil, errors.New("persist storage was not been initialized.")
	}

	return rdbmsInst, nil
}
