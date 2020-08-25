package server

import (
	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/cli"
	"github.com/shylinux/icebergs/base/tcp"
	"github.com/shylinux/icebergs/core/code"
	kit "github.com/shylinux/toolkits"

	"path"
	"runtime"
	"strings"
)

const MYSQL = "mysql"
const _my_cnf = `
[mysqld]
basedir = ./
datadir = ./data
port = %s
socket = ./mysqld.socket

sql_mode=NO_ENGINE_SUBSTITUTION,STRICT_TRANS_TABLES 
[mysqld_safe]
log-error = ./mysqld.log
pid-file = ./mysqld.pid
`

var Index = &ice.Context{Name: MYSQL, Help: "mysql",
	Configs: map[string]*ice.Config{
		MYSQL: {Name: MYSQL, Help: "mysql", Value: kit.Data(
			"windows", "https://mirrors.tuna.tsinghua.edu.cn/mysql/downloads/MySQL-5.6/mysql-5.6.48.zip",
			"darwin", "https://mirrors.tuna.tsinghua.edu.cn/mysql/downloads/MySQL-5.6/mysql-5.6.48.tar.gz",
			"linux", "https://mirrors.tuna.tsinghua.edu.cn/mysql/downloads/MySQL-5.6/mysql-5.6.48.tar.gz",
		)},
	},
	Commands: map[string]*ice.Command{
		ice.CTX_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) { m.Load() }},
		ice.CTX_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) { m.Save() }},

		MYSQL: {Name: "mysql port=auto auto 启动:button 编译:button 下载:button cmd:textarea", Help: "mysql", Action: map[string]*ice.Action{
			"download": {Name: "download", Help: "下载", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(code.INSTALL, "download", m.Conf(MYSQL, kit.Keys(kit.MDB_META, runtime.GOOS)))
			}},
			"compile": {Name: "compile", Help: "编译", Hand: func(m *ice.Message, arg ...string) {
				name := path.Base(strings.TrimSuffix(strings.TrimSuffix(m.Conf(MYSQL, kit.Keys(kit.MDB_META, runtime.GOOS)), ".tar.gz"), "zip"))
				m.Option(cli.CMD_DIR, path.Join(m.Conf(code.INSTALL, kit.META_PATH), name))
				m.Cmdy(cli.SYSTEM, "cmake", "./", "-DCMAKE_INSTALL_PREFIX=./install",
					"-DDEFAULT_COLLATION=utf8_general_ci", "-DDEFAULT_CHARSET=utf8",
					"-DEXTRA_CHARSETS=all",
				)

				m.Cmdy(cli.SYSTEM, "make", "-j8")
				m.Cmdy(cli.SYSTEM, "make", "install")
			}},
			"start": {Name: "start", Help: "启动", Hand: func(m *ice.Message, arg ...string) {
				// 分配
				port := m.Cmdx(tcp.PORT, "select")
				p := path.Join(m.Conf(cli.DAEMON, kit.META_PATH), port)

				// 复制
				name := path.Base(strings.TrimSuffix(strings.TrimSuffix(m.Conf(MYSQL, kit.Keys(kit.MDB_META, runtime.GOOS)), ".tar.gz"), "zip"))
				from := kit.Path(path.Join(m.Conf(code.INSTALL, kit.META_PATH), name, "install"))
				m.Cmdy(cli.SYSTEM, "cp", "-r", from, p)

				// 生成
				m.Option(cli.CMD_DIR, p)
				m.Cmd(cli.SYSTEM, "./scripts/mysql_install_db", "--datadir=./data")
				if f, _, e := kit.Create(path.Join(p, "my.cnf")); m.Assert(e) {
					f.WriteString(kit.Format(_my_cnf, port))
				}

				// 启动
				m.Option(cli.CMD_STDOUT, path.Join(p, "data/mysqld.log"))
				m.Option(cli.CMD_STDERR, path.Join(p, "data/mysqld.log"))
				m.Cmdy(cli.DAEMON, "bin/mysqld",
					"--basedir=./", "--datadir=./data", "--plugin-dir=./lib/plugin",
					"--log-error=./mysqld.log", "--pid-file=./mysqld.pid",
					"--socket=./mysqld.socket", "--port="+port)
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			if len(arg) > 0 && arg[0] != "" {
				m.Option(cli.CMD_DIR, path.Join("var/daemon", arg[0]))
				m.Cmdy(cli.SYSTEM, "bin/mysql", "-S", "data/mysqld.socket", "-u", "root", "-e", kit.Select("show databases", arg, 1))
				return
			}

			m.Cmd(cli.DAEMON).Table(func(index int, value map[string]string, head []string) {
				if strings.HasPrefix(value[kit.MDB_NAME], "bin/mysql") {
					m.Push(kit.MDB_TIME, value[kit.MDB_TIME])
					m.Push(kit.MDB_PID, value[kit.MDB_PID])
					m.Push(kit.MDB_DIR, value[kit.MDB_DIR])
					m.Push(kit.MDB_PORT, path.Base(value[kit.MDB_DIR]))
					m.Push(kit.MDB_STATUS, value[kit.MDB_STATUS])
					m.Push(kit.MDB_NAME, value[kit.MDB_NAME])
				}
			})
			m.Sort("time", "time_r")
		}},
	},
}

func init() { code.Index.Register(Index, nil) }
