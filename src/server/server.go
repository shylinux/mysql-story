package server

import (
	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/cli"
	"github.com/shylinux/icebergs/core/code"
	kit "github.com/shylinux/toolkits"

	"path"
	"runtime"
)

const (
	MYSQL  = "mysql"
	SERVER = "server"
)

var Index = &ice.Context{Name: MYSQL, Help: "mysql",
	Configs: map[string]*ice.Config{
		SERVER: {Name: SERVER, Help: "服务器", Value: kit.Data(
			"windows", "https://mirrors.tuna.tsinghua.edu.cn/mysql/downloads/MySQL-5.6/mysql-5.6.48.zip",
			"darwin", "https://mirrors.tuna.tsinghua.edu.cn/mysql/downloads/MySQL-5.6/mysql-5.6.48.tar.gz",
			"linux", "https://mirrors.tuna.tsinghua.edu.cn/mysql/downloads/MySQL-5.6/mysql-5.6.48.tar.gz",

			"build", []string{
				"-DCMAKE_INSTALL_PREFIX=./install",
				"-DDEFAULT_COLLATION=utf8_general_ci",
				"-DDEFAULT_CHARSET=utf8",
				"-DEXTRA_CHARSETS=all",
			},
			"start", []string{
				"--basedir=./", "--datadir=./data", "--plugin-dir=./lib/plugin",
				"--log-error=./mysqld.log", "--pid-file=./mysqld.pid",
				"--socket=./mysqld.socket",
			},
		)},
	},
	Commands: map[string]*ice.Command{
		ice.CTX_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) { m.Load() }},
		ice.CTX_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) { m.Save() }},

		SERVER: {Name: "server port=auto path=auto auto 启动:button 构建:button 下载:button", Help: "服务器", Action: map[string]*ice.Action{
			"download": {Name: "download", Help: "下载", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(code.INSTALL, "download", m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS)))
			}},
			"build": {Name: "build", Help: "构建", Hand: func(m *ice.Message, arg ...string) {
				m.Optionv("prepare", func(p string) {
					m.Option(cli.CMD_DIR, p)
					m.Cmdy(cli.SYSTEM, "cmake", "./", m.Confv(SERVER, "meta.build"))
				})
				m.Cmdy(code.INSTALL, "build", m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS)))
			}},
			"start": {Name: "start", Help: "启动", Hand: func(m *ice.Message, arg ...string) {
				m.Optionv("prepare", func(p string) []string {
					m.Option(cli.CMD_DIR, p)
					m.Cmd(cli.SYSTEM, "./scripts/mysql_install_db", "--datadir=./data")
					return []string{"--port", path.Base(p)}
				})
				m.Cmdy(code.INSTALL, "start", m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS)),
					"bin/mysqld", m.Confv(SERVER, "meta.start"))
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			m.Cmdy(code.INSTALL, path.Base(m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS))), arg)
		}},
	},
}

func init() { code.Index.Register(Index, nil) }
