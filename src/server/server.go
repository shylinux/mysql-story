package server

import (
	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/core/code"
	kit "github.com/shylinux/toolkits"

	"runtime"
)

const MYSQL = "mysql"

var Index = &ice.Context{Name: MYSQL, Help: "mysql",
	Configs: map[string]*ice.Config{
		MYSQL: {Name: MYSQL, Help: "mysql", Value: kit.Data(
			"windows", "https://mirrors.tuna.tsinghua.edu.cn/mysql/downloads/MySQL-5.6/mysql-5.6.48.zip",
			"darwin", "https://mirrors.tuna.tsinghua.edu.cn/mysql/downloads/MySQL-5.6/mysql-5.6.48.tar.gz",
			"linux", "https://mirrors.tuna.tsinghua.edu.cn/mysql/downloads/MySQL-5.6/mysql-5.6.48.tar.gz",
		)},
	},
	Commands: map[string]*ice.Command{
		ice.CTX_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},
		ice.CTX_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},

		MYSQL: {Name: "mysql port=auto auto 启动:button 编译:button 下载:button", Help: "mysql", Action: map[string]*ice.Action{
			"download": {Name: "download", Help: "下载", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(code.INSTALL, "download", m.Conf(MYSQL, kit.Keys(kit.MDB_META, runtime.GOOS)))
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
		}},
	},
}

func init() { code.Index.Register(Index, nil) }
