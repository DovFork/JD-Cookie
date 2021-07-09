package repo

import (
	"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
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
	if !ds.DbEnable {
		return errors.New("db not enable")
	}
	rdbmsInst = &rdbms{}
	switch ds.DbType {
	case "mysql":
		d := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", ds.DbUser, ds.DbPass, ds.DbHost, ds.DbPort, ds.DbDatabase)
		log.Info(d)
		dsn = append(dsn, d)
		break
	case "postgres":
		d := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", ds.DbUser, ds.DbPass, ds.DbHost, ds.DbPort, ds.DbDatabase)
		log.Info(d)
		dsn = append(dsn, d)
		break
	case "sqlite3":
		d := fmt.Sprintf("%s", ds.DbHost)
		log.Info(d)
		dsn = append(dsn, d)
		break
	case "mssql":
		d := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&connection+timeout=30", ds.DbUser, ds.DbPass, ds.DbHost, ds.DbPort, ds.DbDatabase)
		log.Info(d)
		dsn = append(dsn, d)
	default:
		log.Fatalf("not supported db type ! only for mysql、postgres、sqlite3、mssql")
	}

	rdbmsInst.db, err = xorm.NewEngineGroup(ds.DbType, dsn)
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
