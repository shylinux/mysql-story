package db

import (
	"os"
	"path"

	"gorm.io/driver/sqlite"
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/web"
)

type sqlite3 struct{ Driver }

func (s sqlite3) Init(m *ice.Message, arg ...string) {
	m.Cmd(ctx.CONFIG, web.COMPILE, "env.CGO_ENABLED", "1")
	p := "var/db/" + m.PrefixKey() + ".db"
	s.Driver.Init(m, func() Dialector { os.MkdirAll(path.Dir(p), 0755); return sqlite.Open(p) })
}

func init() { ice.Cmd(prefixKey(), sqlite3{}) }
