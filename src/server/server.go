package server

import (
	"path"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

type server struct {
	ice.Code
	action  string `data:"xterm"`
	linux   string `data:"https://cdn.mysql.com/archives/mysql-5.6/mysql-5.6.33-linux-glibc2.5-x86_64.tar.gz"`
	darwin  string `data:"https://cdn.mysql.com/archives/mysql-5.6/mysql-5.6.33-osx10.11-x86_64.tar.gz"`
	windows string `data:"https://cdn.mysql.com/archives/mysql-5.6/mysql-5.6.33-winx64.zip"`
	source  string `data:"https://cdn.mysql.com/archives/mysql-5.6/mysql-5.6.33.tar.gz"`
	start   string `name:"start port*=10001 username*=root password*=root" help:"启动"`
	list    string `name:"list port path auto start install build download" help:"数据库"`
}

func (s server) Build(m *ice.Message, arg ...string) {
	s.Code.Build(m, "", func(p string) {
		s.Code.System(m, p, "cmake", "./", "-DCMAKE_INSTALL_PREFIX=_install", "-DDEFAULT_COLLATION=utf8_general_ci", "-DDEFAULT_CHARSET=utf8", "-DEXTRA_CHARSETS=all")
	})
	s.Code.Order(m)
}
func (s server) Start(m *ice.Message, arg ...string) {
	s.Code.Start(m, "", "bin/mysqld", "--basedir=./", "--datadir=./data", "--plugin-dir=./lib/plugin", "--socket=mysqld.socket", "--log-error=mysqld.log", "--pid-file=mysqld.pid", func(p string) []string {
		s.Code.System(m.Spawn(), p, "scripts/mysql_install_db", "--datadir=data")
		return []string{"--port", path.Base(p)}
	})
	m.Sleep3s()
	s.Code.System(m, m.Option(cli.CMD_DIR), "bin/mysql", "-S", "./data/mysqld.socket", "-u", m.Option(aaa.USERNAME),
		"-e", kit.Format("set password for %s@%s = password('%s')", m.Option(aaa.USERNAME), tcp.LOCALHOST, m.Option(aaa.PASSWORD)))

}
func (s server) Xterm(m *ice.Message, arg ...string) {
	s.Code.Xterm(m, "", []string{mdb.TYPE, kit.Format("bin/mysql -h 127.0.0.1 -P %s", m.Option(tcp.PORT)), nfs.PATH, kit.Path(m.Option(nfs.DIR)) + nfs.PS}, arg...)
}
func (s server) List(m *ice.Message, arg ...string) {
	s.Code.List(m, "", arg...)
}

func init() { ice.CodeModCmd(server{}) }
