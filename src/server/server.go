package server

import (
	"os"
	"path"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

type Server struct {
	ice.Code

	source string `data:"http://mirrors.tencent.com/ubuntu/pool/universe/m/mysql-5.6/mysql-5.6_5.6.33.orig.tar.gz"`
	start  string `name:"start port" help:"启动"`

	username string `data:"root"`
	password string `data:"root"`
}

func (s Server) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case tcp.PORT:
		if s.List(m); m.Length() > 0 {
			m.Cut("port,status,time")
		} else {
			m.Cmdy(tcp.PORT)
		}
	}
}
func (s Server) Download(m *ice.Message, arg ...string) {
	s.Code.Download(m, m.Config(nfs.SOURCE), arg...)
}
func (s Server) Build(m *ice.Message, arg ...string) {
	s.Code.Prepare(m, func(p string) {
		s.Code.System(m, p, "cmake", "./",
			"-DCMAKE_INSTALL_PREFIX=./_install",
			"-DDEFAULT_COLLATION=utf8_general_ci",
			"-DDEFAULT_CHARSET=utf8",
			"-DEXTRA_CHARSETS=all")
	})
	s.Code.Build(m, s.Code.PathOther(m, m.Config(nfs.SOURCE)), arg...)
}
func (s Server) Start(m *ice.Message, arg ...string) {
	args := []string{"bin/mysqld",
		"--basedir=./", "--datadir=./data", "--plugin-dir=./lib/plugin",
		"--log-error=./mysqld.log", "--pid-file=./mysqld.pid",
		"--socket=./mysqld.socket",
	}

	if kit.Int(m.Option(tcp.PORT)) >= 10000 {
		p := kit.Path(m.Conf(cli.DAEMON, kit.Keym(nfs.PATH)), m.Option(tcp.PORT))
		if _, e := os.Stat(p); e == nil {
			s.Code.Daemon(m, p, append(args, "--port", m.Option(tcp.PORT))...)
		}
		return // 重启服务
	}

	// 启动服务
	s.Code.Prepare(m, func(p string) []string {
		s.Code.System(m, p, "./scripts/mysql_install_db", "--datadir=./data")
		return []string{"--port", path.Base(p)}
	})
	s.Code.Start(m, s.Code.PathOther(m, m.Config(nfs.SOURCE)), args...)

	// 设置密码
	m.Sleep("3s")
	username, password := m.Config(aaa.USERNAME), m.Config(aaa.PASSWORD)
	s.Code.System(m, m.Option(cli.CMD_DIR), "bin/mysql", "-S", "data/mysqld.socket", "-u", username,
		"-e", kit.Format("set password for %s@%s = password('%s')", username, tcp.LOCALHOST, password))
}
func (s Server) List(m *ice.Message, arg ...string) {
	s.Code.List(m, s.Code.PathOther(m, m.Config(nfs.SOURCE)), arg...)
}

func init() { ice.Cmd("web.code.mysql.server", Server{}) }
