package db

import (
	"gorm.io/driver/mysql"
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/tcp"
	"shylinux.com/x/mysql-story/src/client"
	kit "shylinux.com/x/toolkits"
)

type Mysql struct {
	Driver
	Client client.Client
}

func (s Mysql) BeforeMigrate(m *ice.Message, arg ...string) {
	m.Cmd(s.Client).Table(func(value ice.Maps) {
		dsn := kit.Format("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
			value[aaa.USERNAME], value[aaa.PASSWORD], value[tcp.HOST], value[tcp.PORT], value[client.DATABASE])
		s.Driver.Register(m, func() Dialector { return mysql.Open(dsn) }, value[aaa.SESS])
	})
}

func init() { ice.Cmd(prefixKey(), Mysql{}) }
