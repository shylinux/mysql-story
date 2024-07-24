package sqlite

import (
	"path"

	"gorm.io/driver/sqlite"
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/web"
	"shylinux.com/x/mysql-story/src/db"
)

type Sqlite struct{ db.Driver }

func (s Sqlite) Init(m *ice.Message, arg ...string) {
	m.Cmd(ctx.CONFIG, web.COMPILE, "env.CGO_ENABLED", "1")
}
func (s Sqlite) BeforeMigrate(m *ice.Message, arg ...string) {
	p := path.Join("var/db/", m.PrefixKey()+".db")
	s.Driver.Register(m, func() db.Dialector { m.MkdirAll(path.Dir(p)); return sqlite.Open(p) })
}

func init() { ice.Cmd("web.code.db.sqlite", Sqlite{}) }
