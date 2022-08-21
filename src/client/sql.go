package client

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/core/code"
	kit "shylinux.com/x/toolkits"
)

type sql struct{ ice.Lang }

func (h sql) Init(m *ice.Message, arg ...string) {
	h.Lang.Init(m, code.SPLIT, kit.Dict(code.SPACE, "\t ", code.OPERATE, ""),
		code.PREFIX, kit.Dict("<!-- ", code.COMMENT), code.PREPARE, kit.Dict(
			code.KEYWORD, kit.Simple(
				"create", "table", "if", "not", "exists",
				"index",
				"on",
			),
			code.CONSTANT, kit.Simple(
				"InnoDB", "utf8mb4", "null",
			),
			code.DATATYPE, kit.Simple(
				"unsigned",
				"bigint",
				"tinyint",
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
func (h sql) Render(m *ice.Message, arg ...string) {}
func (h sql) Engine(m *ice.Message, arg ...string) {}

func init() { ice.CodeModCmd(sql{}) }
