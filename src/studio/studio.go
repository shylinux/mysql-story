package client

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/mysql-story/src/client"
)

type studio struct {
	query  client.Query
	client client.Client
	script client.Script
	tools  string `data:"web.code.mysql.server,web.code.mysql.client"`
	create string `name:"create sess*=biz driver*=mysql database*=mysql host*=127.0.0.1 port*=10001 username*=root password*=root"`
	list   string `name:"list refresh" icon:"studio.png"`
}

func (s studio) Inputs(m *ice.Message, arg ...string) { m.Cmdy(s.client, m.ActionKey(), arg) }
func (s studio) Create(m *ice.Message, arg ...string) { m.Cmdy(s.client, m.ActionKey(), arg) }
func (s studio) Remove(m *ice.Message, arg ...string) { m.Cmdy(s.client, m.ActionKey(), arg) }
func (s studio) Xterm(m *ice.Message, arg ...string)  { m.Cmdy(s.client, m.ActionKey(), arg) }
func (s studio) Script(m *ice.Message, arg ...string) { m.Cmdy(s.script, arg) }
func (s studio) List(m *ice.Message, arg ...string) {
	if m.Cmdy(s.query, arg); m.Length() == 0 {
		m.EchoInfoButton("please create sess", s.Create)
		return
	} else if len(arg) == 0 {
		m.PushAction(s.client.Remove)
	} else if len(arg) == 1 {
		m.PushAction(s.client.Xterm)
	}
	m.Display("").DisplayCSS("")
}

func init() { ice.CodeModCmd(studio{}) }
