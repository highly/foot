package orm

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
)

var mu = &sync.Mutex{}
var gormDBs = make(map[string]*gorm.DB)

func R(scope string) *gorm.DB {
	db, _ := New(scope)
	return db
}

func D() *gorm.DB {
	db, _ := New("default")
	return db
}

func New(scope string) (*gorm.DB, error) {
	mu.Lock()
	defer mu.Unlock()

	if _, exist := gormDBs[scope]; !exist {
		o := OptionsFromConfig(scope)
		db, err := gorm.Open("mysql", o.DSN())
		if err != nil {
			return nil, err
		}
		db.LogMode(true)
		db.SetLogger(logger{})
		db.DB().SetMaxIdleConns(o.MaxIdleConns)
		db.DB().SetMaxOpenConns(o.MaxOpenConns)
		db.DB().SetConnMaxLifetime(o.ConnMaxLifetime)
		addGormCallbacks(db)
		gormDBs[scope] = db
	}
	return gormDBs[scope], nil
}
