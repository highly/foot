package orm

import (
	"encoding/json"
	"fmt"
	"github.com/highly/foot/config"
	"github.com/highly/foot/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
	"strings"
	"sync"
)

var mu = &sync.Mutex{}
var gormDBs = make(map[string]*gorm.DB)

func NewGorm() (*gorm.DB, error) {
	return NewGormWithName("default")
}

func NewGormWithName(scope string) (*gorm.DB, error) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := gormDBs[scope]; !ok {
		dsn := dsn(scope)
		log.Debugf("connect dsn: %s", dsn)
		db, err := gorm.Open("mysql", dsn)
		if err != nil {
			return nil, err
		}
		db.LogMode(true)
		db.SetLogger(gormLogger{})
		if maxIdleConns := config.GetInt(fmt.Sprintf("Mysql.%s.MaxIdleConns", scope)); maxIdleConns > 0 {
			log.Debugf("set msyql MaxIdleConns: %d", maxIdleConns)
			db.DB().SetMaxIdleConns(maxIdleConns)
		}
		if maxOpenConns := config.GetInt(fmt.Sprintf("Mysql.%s.MaxOpenConns", scope)); maxOpenConns > 0 {
			log.Debugf("set mysql MaxOpenConns: %d", maxOpenConns)
			db.DB().SetMaxOpenConns(maxOpenConns)
		}
		if connMaxLifetime := config.GetDuration(fmt.Sprintf("Mysql.%s.ConnMaxLifetime", scope)); connMaxLifetime > 0 {
			log.Debugf("set mysql ConnMaxLifetime: %d", connMaxLifetime)
			db.DB().SetConnMaxLifetime(connMaxLifetime)
		}
		gormDBs[scope] = db
		addGormCallbacks(db)
	}
	return gormDBs[scope], nil
}

// AddGormCallbacks adds callbacks for tracing, you should call Wrap to make them work
func addGormCallbacks(db *gorm.DB) {
	callbacks := newCallbacks()
	registerCallbacks(db, "create", callbacks)
	registerCallbacks(db, "query", callbacks)
	registerCallbacks(db, "update", callbacks)
	registerCallbacks(db, "delete", callbacks)
	registerCallbacks(db, "row_query", callbacks)
}

type callbacks struct{}

func newCallbacks() *callbacks {
	return &callbacks{}
}

func (c *callbacks) beforeCreate(scope *gorm.Scope)   { c.before(scope) }
func (c *callbacks) afterCreate(scope *gorm.Scope)    { c.after(scope, "INSERT") }
func (c *callbacks) beforeQuery(scope *gorm.Scope)    { c.before(scope) }
func (c *callbacks) afterQuery(scope *gorm.Scope)     { c.after(scope, "SELECT") }
func (c *callbacks) beforeUpdate(scope *gorm.Scope)   { c.before(scope) }
func (c *callbacks) afterUpdate(scope *gorm.Scope)    { c.after(scope, "UPDATE") }
func (c *callbacks) beforeDelete(scope *gorm.Scope)   { c.before(scope) }
func (c *callbacks) afterDelete(scope *gorm.Scope)    { c.after(scope, "DELETE") }
func (c *callbacks) beforeRowQuery(scope *gorm.Scope) { c.before(scope) }
func (c *callbacks) afterRowQuery(scope *gorm.Scope)  { c.after(scope, "") }

func (c *callbacks) before(scope *gorm.Scope) {
}

func (c *callbacks) after(scope *gorm.Scope, operation string) {
	var vars, table, method, errs, count, result zap.Field
	if operation == "" {
		operation = strings.ToUpper(strings.Split(scope.SQL, " ")[0])
	}
	if buf, err := json.Marshal(scope.SQLVars); err == nil {
		vars = zap.String("db.vars", string(buf))
	} else {
		vars = zap.String("db.vars", "")
	}
	table = zap.String("db.table", scope.TableName())
	method = zap.String("db.method", operation)
	errs = zap.Bool("db.err", scope.HasError())
	count = zap.Int64("db.count", scope.DB().RowsAffected)
	if buf, err := json.Marshal(scope.DB().Value); err == nil {
		result = zap.String("db.result", string(buf))
	} else {
		result = zap.String("db.result", "")
	}
	log.Info("sql record", vars, table, method, errs, count, result)
	if scope.HasError() && !gorm.IsRecordNotFoundError(scope.DB().Error) {
		log.Error("sql error",
			zap.String("sql", scope.SQL),
			zap.Any("vars", scope.SQLVars),
			zap.String("error", scope.DB().Error.Error()),
			zap.String("traceid", ""),
			zap.String("spanid", ""),
		)
	}
}

func registerCallbacks(db *gorm.DB, name string, c *callbacks) {
	beforeName := fmt.Sprintf("tracing:%v_before", name)
	afterName := fmt.Sprintf("tracing:%v_after", name)
	gormCallbackName := fmt.Sprintf("gorm:%v", name)
	// gorm does some magic, if you pass CallbackProcessor here - nothing works
	switch name {
	case "create":
		db.Callback().Create().Before(gormCallbackName).Register(beforeName, c.beforeCreate)
		db.Callback().Create().After(gormCallbackName).Register(afterName, c.afterCreate)
	case "query":
		db.Callback().Query().Before(gormCallbackName).Register(beforeName, c.beforeQuery)
		db.Callback().Query().After(gormCallbackName).Register(afterName, c.afterQuery)
	case "update":
		db.Callback().Update().Before(gormCallbackName).Register(beforeName, c.beforeUpdate)
		db.Callback().Update().After(gormCallbackName).Register(afterName, c.afterUpdate)
	case "delete":
		db.Callback().Delete().Before(gormCallbackName).Register(beforeName, c.beforeDelete)
		db.Callback().Delete().After(gormCallbackName).Register(afterName, c.afterDelete)
	case "row_query":
		db.Callback().RowQuery().Before(gormCallbackName).Register(beforeName, c.beforeRowQuery)
		db.Callback().RowQuery().After(gormCallbackName).Register(afterName, c.afterRowQuery)
	}
}

type gormLogger struct{}

func (gormLogger) Print(v ...interface{}) {
	log.Debug(fmt.Sprint(v...))
}

func dsn(scope string) string {
	username := config.GetString(fmt.Sprintf("Mysql.%s.Username", scope))
	password := config.GetString(fmt.Sprintf("Mysql.%s.Password", scope))
	host := config.GetString(fmt.Sprintf("Mysql.%s.Host", scope))
	port := config.GetInt(fmt.Sprintf("Mysql.%s.Port", scope))
	database := config.GetString(fmt.Sprintf("Mysql.%s.Database", scope))
	charset := config.GetString(fmt.Sprintf("Mysql.%s.Charset", scope))
	format := "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local"
	return fmt.Sprintf(format, username, password, host, port, database, charset)
}
