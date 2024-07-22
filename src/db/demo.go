package db

import (
	"shylinux.com/x/ice"
)

type Demo struct {
	Model
	Name string
}
type demo struct {
	Table
	driver string `data:"sqlite"`
	create string `name:"create name*"`
	list   string `name:"list id auto"`
}

func (s demo) Init(m *ice.Message, arg ...string) {
	s.Table.Init(m, &Demo{})
}

func init() { ice.Cmd(prefixKey(), demo{}) }
