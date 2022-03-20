package client

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

const (
	TARGET = "target"
	METHOD = "method"
)

type Grant struct {
	Client

	grants string `name:"grants session method=all target=*.* host=% username password" help:"授权"`
	revoke string `name:"revoke session method=all target=*.* host username" help:"撤销"`
	drop   string `name:"drop" help:"删除"`
	list   string `name:"list session auto" help:"权限"`
}

func (g Grant) Grants(m *ice.Message, arg ...string) {
	_sql_exec(m, g.sql_meta(m, m.Option(SESSION), MYSQL), kit.Format("grant %s on %s to '%s'@'%s' identified by '%s'",
		m.Option(METHOD), m.Option(TARGET), m.Option(aaa.USERNAME), m.Option(tcp.HOST), m.Option(aaa.PASSWORD)))
	m.SetAppend()
}
func (g Grant) Revoke(m *ice.Message, arg ...string) {
	_sql_exec(m, g.sql_meta(m, m.Option(SESSION), MYSQL), kit.Format("revoke %s on %s from '%s'@'%s'",
		m.Option(METHOD), m.Option(TARGET), m.Option(aaa.USERNAME), m.Option(tcp.HOST)))
	m.SetAppend()
}
func (g Grant) Drop(m *ice.Message, arg ...string) {
	_sql_exec(m, g.sql_meta(m, m.Option(SESSION), MYSQL), kit.Format("drop user '%s'@'%s'", m.Option(aaa.USERNAME), m.Option(tcp.HOST)))
	m.SetAppend()
}
func (g Grant) List(m *ice.Message, arg ...string) {
	if len(arg) < 1 || arg[0] == "" { // 连接列表
		m.Action(g.Create)
		m.Cmdy(g.Client)
		return
	}

	m.Action(g.Grants)
	_sql_query(m, g.sql_meta(m, arg[0], MYSQL), kit.Format("select User,Host from user")).ToLowerAppend().RenameAppend("user", aaa.USERNAME).Table(func(index int, value map[string]string, head []string) {
		msg := _sql_query(m.Spawn(), g.sql_meta(m, arg[0], MYSQL), kit.Format("show grants for '%s'@'%s'", value[aaa.USERNAME], value[tcp.HOST]))
		m.Push("stm", msg.Append(""))
	})
	m.Sort("username,host")
	m.PushAction(g.Revoke, g.Drop)
	m.StatusTimeCount()
}
func init() { ice.CodeModCmd(Grant{}) }
