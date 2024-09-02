package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/tcp"
	"shylinux.com/x/mysql-story/src/client"
	"shylinux.com/x/mysql-story/src/db"
	kit "shylinux.com/x/toolkits"
)

type Mysql struct {
	db.Driver
	Client client.Client
}

func (s Mysql) BeforeMigrate(m *ice.Message, arg ...string) {
	m.Cmd(s.Client).Table(func(value ice.Maps) {
		s.Driver.Register(m, func(db string) db.Dialector {
			dsn := kit.Format("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", value[aaa.USERNAME], value[aaa.PASSWORD], value[tcp.HOST], value[tcp.PORT], db)
			m.Info("open db %s", dsn)
			return mysql.Open(dsn)
		}, value[aaa.SESS])
		s.AutoCreate(m.Options(value))
	})
}
func (s Mysql) AutoCreate(m *ice.Message, arg ...string) {
	dsn := kit.Format("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", m.Option(aaa.USERNAME), m.Option(aaa.PASSWORD), m.Option(tcp.HOST), m.Option(tcp.PORT), m.Option("database"))
	m.Info("open db %s", dsn)
	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}); !m.Warn(err) {
		list := map[string]bool{}
		m.Cmd("web.code.db.models").Table(func(value ice.Maps) {
			if name := kit.Split(value[mdb.NAME], ".")[0]; !list[name] {
				db.Exec(kit.Format("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;", name))
				list[name] = true
			}
		})
	}
}
func (s Mysql) List(m *ice.Message, arg ...string) {
	m.Cmdy(s.Client).PushAction(s.AutoCreate)
}

func init() { ice.Cmd("web.code.db.mysql", Mysql{}) }
