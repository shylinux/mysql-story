package client

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
)

type shell struct {
	client client
	list   string `name:"list sess auto" help:"终端"`
}

func (s shell) List(m *ice.Message, arg ...string) {
	if len(arg) == 0 {
		m.Cmdy(s.client)
		return
	}
	m.Option(aaa.SESS, arg[0])
	s.client.Xterm(m)
}

func init() { ice.CodeModCmd(shell{}) }
