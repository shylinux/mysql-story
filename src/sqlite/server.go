package sqlite

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/web/html"
)

const (
	SQLITE3 = "sqlite3"
)

type server struct {
	ice.Code
	source string `data:"https://sqlite.org/2023/sqlite-autoconf-3420000.tar.gz"`
	list   string `name:"list"`
}

func (s server) Build(m *ice.Message, arg ...string) {
	s.Code.Build(m)
	s.Code.Order(m)
}
func (s server) List(m *ice.Message, arg ...string) {
	if m.Exists(s.Path(m, "", "_install")) {
		m.ProcessXterm(SQLITE3, SQLITE3, arg...).Push(ctx.STYLE, html.OUTPUT)
	} else if m.Exists(s.Path(m, "")) {
		m.EchoInfoButton("please build sqlite3", s.Build)
	} else {
		m.EchoInfoButton("please download sqlite3", s.Download)
	}
}

func init() { ice.CodeCtxCmd(server{}) }
