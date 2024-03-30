package sqlite

import "shylinux.com/x/ice"

const (
	SQLITE3 = "sqlite3"
)

type server struct {
	ice.Code
	source string `data:"https://sqlite.org/2023/sqlite-autoconf-3420000.tar.gz"`
	list   string `name:"list path auto xterm build download" help:"数据库"`
}

func (s server) Build(m *ice.Message, arg ...string) { s.Code.Build(m); s.Code.Order(m) }
func (s server) Xterm(m *ice.Message, arg ...string) { s.Code.Xterm(m, "", SQLITE3, arg...) }
func (s server) List(m *ice.Message, arg ...string)  { s.Code.Source(m, "", arg...) }

func init() { ice.CodeCtxCmd(server{}) }
