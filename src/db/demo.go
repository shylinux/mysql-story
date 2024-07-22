package db

import (
	"shylinux.com/x/ice"
)

type demoTable struct {
	Model
	Name string
}
type demo struct {
	Table
	driver string `data:"sqlite"`
	create string `name:"create name*"`
	list   string `name:"list id auto" help:"示例"`
}

func (s demo) Init(m *ice.Message, arg ...string) {
	s.Table.Init(m, &demoTable{})
}

func init() { ice.Cmd(prefixKey(), demo{}) }
