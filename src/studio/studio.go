package client

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/mysql-story/src/client"
)

type studio struct {
	client client.Client
	query  client.Query
	create string `name:"create sess*=biz username*=root password*=root host*=127.0.0.1 port*=10001 database*=mysql driver*=mysql"`
	list   string `name:"list refresh" icon:"mysql.png"`
}

func (s studio) Create(m *ice.Message, arg ...string) {
	m.Cmdy(s.client, m.ActionKey(), arg)
}
func (s studio) List(m *ice.Message, arg ...string) {
	m.Cmdy(s.query, arg).Display("")
}

func init() { ice.CodeModCmd(studio{}) }
