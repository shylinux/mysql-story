package server

import (
	"path"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

type server struct {
	ice.Code
	source string `data:"http://mirrors.tencent.com/ubuntu/pool/universe/m/mysql-5.6/mysql-5.6_5.6.33.orig.tar.gz"`
	start  string `name:"start port=10000 username=root password=root" help:"启动"`
	list   string `name:"list port path auto start build download" help:"数据库"`
}

func (s server) Build(m *ice.Message, arg ...string) {
	s.Code.Build(m, "", func(p string) {
		s.Code.System(m, p, "cmake", "./", "-DCMAKE_INSTALL_PREFIX=_install", "-DDEFAULT_COLLATION=utf8_general_ci", "-DDEFAULT_CHARSET=utf8", "-DEXTRA_CHARSETS=all")
	})
}
func (s server) Start(m *ice.Message, arg ...string) {
	s.Code.Start(m, "", "bin/mysqld", "--basedir=./", "--datadir=data", "--plugin-dir=lib/plugin",
		"--log-error=mysqld.log", "--pid-file=mysqld.pid", "--socket=mysqld.socket", func(p string) []string {
			s.Code.System(m.Spawn(), p, "scripts/mysql_install_db", "--datadir=data")
			return []string{"--port", path.Base(p)}
		})

	// 设置密码
	m.Sleep3s()
	s.Code.System(m, m.Option(cli.CMD_DIR), "bin/mysql", "-S", "data/mysqld.socket", "-u", m.Option(aaa.USERNAME),
		"-e", kit.Format("set password for %s@%s = password('%s')", m.Option(aaa.USERNAME), tcp.LOCALHOST, m.Option(aaa.PASSWORD)))
}
func (s server) List(m *ice.Message, arg ...string) {
	s.Code.List(m, "", arg...)
}

func init() { ice.CodeModCmd(server{}) }
