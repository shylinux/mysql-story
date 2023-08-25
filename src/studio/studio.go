package client

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/mysql-story/src/client"
)

type Studio struct {
	client.Client
	list string `name:"list refresh create" help:"数据库" icon:"mysql.png"`
}

func (s Studio) List(m *ice.Message, arg ...string) {
	m.Cmdy(s.Client, arg).Options(ice.MSG_ACTION, "").Display("")
}

func init() { ice.CodeModCmd(Studio{}) }
