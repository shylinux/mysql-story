package client

import (
	"path"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/core/code"
	kit "shylinux.com/x/toolkits"
)

type sql struct{}

func (s sql) Init(m *ice.Message, arg ...string) {
	m.Conf("sql", kit.MDB_META, kit.Dict(
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
	))
	m.Cmd(mdb.PLUGIN, mdb.CREATE, "sql", m.Prefix("sql"))
	m.Cmd(mdb.RENDER, mdb.CREATE, "sql", m.Prefix("sql"))
	code.LoadPlug(m.Message, "sql")
	m.Load()
}
func (s sql) Plugin(m *ice.Message, arg ...string) {
	m.Echo(m.Conf("sql", kit.Keym(code.PLUG)))
}
func (s sql) Render(m *ice.Message, arg ...string) {
	m.Cmdy(nfs.CAT, path.Join(arg[2], arg[1]))
}

func init() { ice.Cmd("web.code.mysql.sql", sql{}) }
