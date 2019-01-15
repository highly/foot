package orm

import (
	"github.com/go-sql-driver/mysql"
	"github.com/highly/foot/config"
	"time"
)

const (
	DefaultCharset         = "utf8"
	DefaultCollation       = "utf8_general_ci"
	DefaultMaxIdleConns    = 2
	DefaultMaxOpenConns    = 0
	DefaultConnMaxLifetime = 0
)

type Options struct {
	User            string
	Passwd          string
	Addr            string
	DBName          string
	Charset         string
	Collation       string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

func (o *Options) DSN() string {
	c := &mysql.Config{
		User:                 o.User,
		Passwd:               o.Passwd,
		Net:                  "tcp",
		Addr:                 o.Addr,
		DBName:               o.DBName,
		ParseTime:            true,
		Loc:                  time.Local,
		Collation:            o.Collation,
		AllowNativePasswords: true,
		MaxAllowedPacket:     4 << 20,
		Params: map[string]string{
			"charset": o.Charset,
		},
	}
	return c.FormatDSN()
}

func OptionsFromConfig(scope string) *Options {
	return &Options{
		User:            config.String("mysql." + scope + ".user"),
		Passwd:          config.String("mysql." + scope + ".password"),
		Addr:            config.String("mysql." + scope + ".addr"),
		DBName:          config.String("mysql." + scope + ".dbname"),
		Charset:         config.DefaultString("mysql."+scope+".charset", DefaultCharset),
		Collation:       config.DefaultString("mysql."+scope+".collation", DefaultCollation),
		MaxIdleConns:    config.DefaultInt("mysql."+scope+".maxIdleConns", DefaultMaxIdleConns),
		MaxOpenConns:    config.DefaultInt("mysql."+scope+".maxOpenConns", DefaultMaxOpenConns),
		ConnMaxLifetime: config.DefaultDuration("mysql."+scope+".connMaxLifetime", DefaultConnMaxLifetime),
	}
}
