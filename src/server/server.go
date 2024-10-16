package server

import (
	"path"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"

	"shylinux.com/x/mysql-story/src/client"
)

const (
	MYSQL      = "mysql"
	BIN_MYSQL  = "bin/mysql"
	BIN_MYSQLD = "bin/mysqld"
)

type server struct {
	ice.Code
	client  client.Client
	action  string `data:"xterm"`
	linux   string `data:"https://cdn.mysql.com/archives/mysql-5.6/mysql-5.6.33-linux-glibc2.5-x86_64.tar.gz"`
	darwin  string `data:"https://cdn.mysql.com/archives/mysql-5.6/mysql-5.6.33-osx10.11-x86_64.tar.gz"`
	windows string `data:"https://cdn.mysql.com/archives/mysql-5.6/mysql-5.6.33-winx64.zip"`
	source  string `data:"https://cdn.mysql.com/archives/mysql-5.6/mysql-5.6.33.tar.gz"`
	start   string `name:"start port*=10001 username=root password=root"`
}

func (s server) Init(m *ice.Message, arg ...string) {
	m.PackageCreateBinary(MYSQL)
}
func (s server) Build(m *ice.Message, arg ...string) {
	s.Code.Build(m, "", func(p string) {
		s.Code.System(m, p, "cmake", "./", "-DCMAKE_INSTALL_PREFIX=_install", "-DDEFAULT_COLLATION=utf8_general_ci", "-DDEFAULT_CHARSET=utf8", "-DEXTRA_CHARSETS=all")
	})
}
func (s server) Start(m *ice.Message, arg ...string) {
	s.Code.Start(m, "", BIN_MYSQLD, "--basedir=./", "--datadir=./data", "--plugin-dir=./lib/plugin", "--socket=mysqld.socket", "--log-error=mysqld.log", "--pid-file=mysqld.pid", func(p string) []string {
		s.Code.System(m.Spawn(), p, "scripts/mysql_install_db", "--datadir=data")
		return []string{"--port", path.Base(p)}
	})
	m.Sleep3s().OptionDefault(aaa.USERNAME, aaa.ROOT, aaa.PASSWORD, kit.HashsUniq())
	s.Code.System(m, m.Option(cli.CMD_DIR), BIN_MYSQL, "-S", "./data/mysqld.socket", "-u", m.Option(aaa.USERNAME), "-e", kit.Format("set password for %s@%s = password('%s')", m.Option(aaa.USERNAME), tcp.LOCALHOST, m.Option(aaa.PASSWORD)))
	m.Cmd(s.client, s.client.Create, aaa.SESS, m.Option(tcp.PORT), aaa.USERNAME, m.Option(aaa.USERNAME), aaa.PASSWORD, m.Option(aaa.PASSWORD), tcp.HOST, "127.0.0.1", tcp.PORT, m.Option(tcp.PORT), client.DATABASE, MYSQL, client.DRIVER, MYSQL)
}
func (s server) Xterm(m *ice.Message, arg ...string) {
	m.ProcessXterm(kit.Keys(MYSQL, m.Option(tcp.PORT)), kit.Format("%s/%s -h %s -P %s", m.Option(nfs.DIR), BIN_MYSQL, "127.0.0.1", m.Option(tcp.PORT)), arg...)
}
func (s server) List(m *ice.Message, arg ...string) {
	s.Code.List(m, "", arg...)
}

func init() { ice.CodeModCmd(server{}) }

type Server struct{ server }

func init() { ice.CodeModCmd(Server{}) }
