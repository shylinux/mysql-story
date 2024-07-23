package db

import (
	"path"

	"gorm.io/driver/sqlite"
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/web"
)

type Sqlite struct{ Driver }

func (s Sqlite) Init(m *ice.Message, arg ...string) {
	m.Cmd(ctx.CONFIG, web.COMPILE, "env.CGO_ENABLED", "1")
}
func (s Sqlite) BeforeMigrate(m *ice.Message, arg ...string) {
	p := path.Join("var/db/", m.PrefixKey()+".db")
	s.Driver.Init(m, func() Dialector { m.MkdirAll(path.Dir(p)); return sqlite.Open(p) })
}

func init() { ice.Cmd(prefixKey(), Sqlite{}) }
