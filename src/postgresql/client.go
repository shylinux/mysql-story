package postgre

import "shylinux.com/x/ice"

type client struct {
	ice.Code

	list string `name:"list port path auto start order build download" help:"示例"`
}

func (s client) List(m *ice.Message, arg ...string) {
	s.Code.List(m, "", arg...)
}

func init() { ice.Cmd("web.code.postgre.client", client{}) }
