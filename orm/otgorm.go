package orm

import (
	"encoding/json"
	"fmt"
	"github.com/highly/foot/log"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"strings"
)

func addGormCallbacks(db *gorm.DB) {
	callbacks := &callbacks{}
	registerCallbacks(db, "create", callbacks)
	registerCallbacks(db, "query", callbacks)
	registerCallbacks(db, "update", callbacks)
	registerCallbacks(db, "delete", callbacks)
	registerCallbacks(db, "row_query", callbacks)
}

type callbacks struct{}

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

func (c *callbacks) before(scope *gorm.Scope) {}

func (c *callbacks) after(scope *gorm.Scope, operation string) {
	var vars, table, method, errs, count, result zap.Field
	if operation == "" {
		operation = strings.ToUpper(strings.Split(scope.SQL, " ")[0])
	}
	vars = zap.String("db.vars", "")
	if len(scope.SQLVars) > 0 {
		if buf, err := json.Marshal(scope.SQLVars); err == nil {
			vars = zap.String("db.vars", string(buf))
		}
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
		)
	}
}

func registerCallbacks(db *gorm.DB, name string, c *callbacks) {
	beforeName := fmt.Sprintf("gorm:%v_before", name)
	afterName := fmt.Sprintf("gorm:%v_after", name)
	gormCallbackName := fmt.Sprintf("gorm:%v", name)

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
