package client

import (
	"path"

	ice "shylinux.com/x/icebergs"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/core/code"
	"shylinux.com/x/mysql-story/src/server"
	kit "shylinux.com/x/toolkits"
)

const SQL = "sql"

func init() {
	server.Index.Merge(&ice.Context{
		Configs: map[string]*ice.Config{
			SQL: {Name: SQL, Help: "语句", Value: kit.Data(
				code.PLUG, kit.Dict(
					code.PREFIX, kit.Dict(
						"-- ", code.COMMENT,
					),
					"_keyword", kit.Dict(
						code.KEYWORD, kit.Simple(
							"CREATE", "DROP", "USE", "IF",
						),
						code.FUNCTION, kit.Simple(
							"DEFAULT", "COMMENT",
							"DATABASE", "TABLE",
						),
						code.DATATYPE, kit.Simple(
							"int", "varchar",
							"datetime",
						),
					),
					code.KEYWORD, kit.Dict(),
				),
			)},
		},
		Commands: map[string]*ice.Command{
			ice.CTX_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
				m.Cmd(mdb.PLUGIN, mdb.CREATE, SQL, m.Prefix(SQL))
				m.Cmd(mdb.RENDER, mdb.CREATE, SQL, m.Prefix(SQL))
				code.LoadPlug(m, SQL)
			}},
			SQL: {Name: "sql", Help: "语句", Action: map[string]*ice.Action{
				mdb.PLUGIN: {Hand: func(m *ice.Message, arg ...string) {
					m.Echo(m.Conf(SQL, kit.Keym(code.PLUG)))
				}},
				mdb.RENDER: {Hand: func(m *ice.Message, arg ...string) {
					m.Cmdy(nfs.CAT, path.Join(arg[2], arg[1]))
				}},
			}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			}},
		},
	})
}
