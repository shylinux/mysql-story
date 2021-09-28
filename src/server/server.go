package server

import (
	"os"
	"path"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/tcp"
	"shylinux.com/x/icebergs/base/web"
	"shylinux.com/x/icebergs/core/code"
	kit "shylinux.com/x/toolkits"
)

type server struct {
	source   string `data:"http://mirrors.tencent.com/slackware/slackware64-14.0/patches/source/mysql/mysql-5.5.52.tar.xz"`
	username string `data:"root"`
	password string `data:"root"`

	inputs   string `name:"inputs" help:"补全"`
	download string `name:"download" help:"下载"`
	build    string `name:"build" help:"构建"`
	start    string `name:"start" help:"启动"`
	list     string `name:"list port path auto start build download" help:"服务器"`
}

func (s server) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case tcp.PORT:
		m.Cmdy(tcp.SERVER)
	}
}
func (s server) Download(m *ice.Message, arg ...string) {
	m.Cmdy(code.INSTALL, web.DOWNLOAD, m.Conf(tcp.SERVER, kit.META_SOURCE))
}
func (s server) Build(m *ice.Message, arg ...string) {
	m.Optionv(code.PREPARE, func(p string) {
		m.Option(cli.CMD_DIR, p)
		m.Cmdy(cli.SYSTEM, "cmake", "./",
			"-DCMAKE_INSTALL_PREFIX=./_install",
			"-DDEFAULT_COLLATION=utf8_general_ci",
			"-DDEFAULT_CHARSET=utf8",
			"-DEXTRA_CHARSETS=all")
	})
	m.Cmdy(code.INSTALL, cli.BUILD, m.Conf(tcp.SERVER, kit.META_SOURCE))
}
func (s server) Start(m *ice.Message, arg ...string) {
	args := []string{"bin/mysqld",
		"--basedir=./", "--datadir=./data", "--plugin-dir=./lib/plugin",
		"--log-error=./mysqld.log", "--pid-file=./mysqld.pid",
		"--socket=./mysqld.socket",
	}

	if kit.Int(m.Option(tcp.PORT)) >= 10000 {
		p := kit.Path(m.Conf(cli.DAEMON, kit.META_PATH), m.Option(tcp.PORT))
		if _, e := os.Stat(p); e == nil {
			m.Option(cli.CMD_DIR, p)
			m.Cmdy(cli.DAEMON, args, "--port", m.Option(tcp.PORT))
		}
		return
	}

	m.Optionv(code.PREPARE, func(p string) []string {
		m.Option(cli.CMD_DIR, p)
		m.Cmd(cli.SYSTEM, "./scripts/mysql_install_db", "--datadir=./data")
		return []string{"--port", path.Base(p)}
	})
	m.Echo(m.Cmdx(code.INSTALL, cli.START, m.Conf(tcp.SERVER, kit.META_SOURCE), args))

	// 设置密码
	m.Sleep("1s")
	username := m.Conf(tcp.SERVER, kit.Keym(aaa.USERNAME))
	password := m.Conf(tcp.SERVER, kit.Keym(aaa.PASSWORD))
	m.Cmd(cli.SYSTEM, "bin/mysql", "-S", "data/mysqld.socket", "-u", username,
		"-e", kit.Format("set password for %s@%s = password('%s')", username, tcp.LOCALHOST, password))

	// // 触发事件
	// m.Event(MYSQL_SERVER_START, aaa.USERNAME, username, aaa.PASSWORD, password,
	// 	tcp.HOST, tcp.LOCALHOST, tcp.PORT, path.Base(m.Option(cli.CMD_DIR)))
}
func (s server) List(m *ice.Message, arg ...string) {
	m.Cmdy(code.INSTALL, m.Conf(tcp.SERVER, kit.META_SOURCE), arg)
}

func init() { ice.Cmd("web.code.mysql.server", server{}) }
