package client

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/core/code"
	kit "shylinux.com/x/toolkits"
)

const (
	SQL = "sql"
)

type sql struct {
	ice.Lang
	ice.Code
	Client
	regexp  string `data:"sql"`
	command string `data:"sql"`

	list string `name:"list sess@key database@key auto" help:"脚本"`
}

func (s sql) Init(m *ice.Message, arg ...string) {
	s.Lang.Init(m, code.SPLIT, kit.Dict(code.SPACE, "\t ", code.OPERATOR, ""),
		code.PREFIX, kit.Dict("<!-- ", code.COMMENT), code.PREPARE, kit.Dict(
			code.KEYWORD, kit.Simple(
				"create", "table", "if", "not", "exists",
				"show",
				"index", "on",
			),
			code.CONSTANT, kit.Simple(
				"InnoDB", "utf8mb4", "null",
			),
			code.DATATYPE, kit.Simple(
				"unsigned",
				"tinyint",
				"bigint",
				"varchar",
				"datetime",
			),
			code.FUNCTION, kit.Simple(
				"ENGINE", "CHARSET",
				"DEFAULT", "comment",
				"auto_increment", "primary",
			),
		))
}
func (s sql) Render(m *ice.Message, arg ...string) {}
func (s sql) Engine(m *ice.Message, arg ...string) {}
func (s sql) List(m *ice.Message, arg ...string) {
	if len(arg) < 2 || arg[0] == "" || arg[1] == "" {
		m.Cmdy(s.Client, arg)
		return
	}
	s.Code.ListScript(m)
}
func (s sql) CatScript(m *ice.Message) {
	m.Cmdy(nfs.CAT, m.Option(nfs.PATH))
}
func (s sql) RunScript(m *ice.Message) {
	s.open(m, m.Option("sess"), m.Option("database"), func(db *Driver) {
		db.Exec(m, "")
	})
}

func init() { ice.CodeModCmd(sql{}) }
