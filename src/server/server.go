package server

import (
	"path"
	"runtime"

	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/aaa"
	"github.com/shylinux/icebergs/base/cli"
	"github.com/shylinux/icebergs/base/gdb"
	"github.com/shylinux/icebergs/base/tcp"
	"github.com/shylinux/icebergs/base/web"
	"github.com/shylinux/icebergs/core/code"
	kit "github.com/shylinux/toolkits"
)

const (
	MYSQL = "mysql"
)
const (
	MYSQL_SERVER_START = "mysql.server.start"
)
const SERVER = "server"

var Index = &ice.Context{Name: MYSQL, Help: "数据库",
	Configs: map[string]*ice.Config{
		SERVER: {Name: SERVER, Help: "服务器", Value: kit.Data(
			cli.WINDOWS, "https://mirrors.tuna.tsinghua.edu.cn/mysql/downloads/MySQL-5.6/mysql-5.6.51.zip",
			cli.DARWIN, "https://mirrors.tuna.tsinghua.edu.cn/mysql/downloads/MySQL-5.6/mysql-5.6.51.tar.gz",
			cli.LINUX, "https://mirrors.tuna.tsinghua.edu.cn/mysql/downloads/MySQL-5.6/mysql-5.6.51.tar.gz",

			gdb.BUILD, []string{
				"-DCMAKE_INSTALL_PREFIX=./_install",
				"-DDEFAULT_COLLATION=utf8_general_ci",
				"-DDEFAULT_CHARSET=utf8",
				"-DEXTRA_CHARSETS=all",
			},
			gdb.START, []string{
				"--basedir=./", "--datadir=./data", "--plugin-dir=./lib/plugin",
				"--log-error=./mysqld.log", "--pid-file=./mysqld.pid",
				"--socket=./mysqld.socket",
			},
			aaa.USERNAME, "root", aaa.PASSWORD, "root",
		)},
	},
	Commands: map[string]*ice.Command{
		ice.CTX_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			m.Load()
		}},
		ice.CTX_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			m.Save()
		}},

		SERVER: {Name: "server port path auto start build download", Help: "服务器", Action: map[string]*ice.Action{
			web.DOWNLOAD: {Name: "download", Help: "下载", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(code.INSTALL, web.DOWNLOAD, m.Conf(SERVER, kit.Keym(runtime.GOOS)))
			}},
			gdb.BUILD: {Name: "build", Help: "构建", Hand: func(m *ice.Message, arg ...string) {
				m.Optionv(code.PREPARE, func(p string) {
					m.Option(cli.CMD_DIR, p)
					m.Cmdy(cli.SYSTEM, "cmake", "./", m.Confv(SERVER, kit.Keym(gdb.BUILD)))
				})
				m.Cmdy(code.INSTALL, gdb.BUILD, m.Conf(SERVER, kit.Keym(runtime.GOOS)))
			}},
			gdb.START: {Name: "start", Help: "启动", Hand: func(m *ice.Message, arg ...string) {
				m.Optionv(code.PREPARE, func(p string) []string {
					m.Option(cli.CMD_DIR, p)
					m.Cmd(cli.SYSTEM, "./scripts/mysql_install_db", "--datadir=./data")
					return []string{"--port", path.Base(p)}
				})
				m.Echo(m.Cmdx(code.INSTALL, gdb.START, m.Conf(SERVER, kit.Keym(runtime.GOOS)),
					"bin/mysqld", m.Confv(SERVER, kit.Keym(gdb.START))))

				// 设置密码
				m.Sleep("1s")
				username := m.Conf(SERVER, kit.Keym(aaa.USERNAME))
				password := m.Conf(SERVER, kit.Keym(aaa.PASSWORD))
				m.Cmd(cli.SYSTEM, "bin/mysql", "-S", "data/mysqld.socket", "-u", username,
					"-e", kit.Format("set password for %s@%s = password('%s')", username, tcp.LOCALHOST, password))

				// 触发事件
				m.Event(MYSQL_SERVER_START, aaa.USERNAME, username, aaa.PASSWORD, password,
					tcp.HOST, tcp.LOCALHOST, tcp.PORT, path.Base(m.Option(cli.CMD_DIR)))
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			m.Cmdy(code.INSTALL, m.Conf(SERVER, kit.Keym(runtime.GOOS)), arg)
		}},
	},
}

func init() { code.Index.Register(Index, nil) }
