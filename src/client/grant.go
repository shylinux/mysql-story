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

func (s Grant) Grants(m *ice.Message, arg ...string) {
	_sql_exec(m, s.meta(m, m.Option(SESSION), MYSQL), kit.Format("grant %s on %s to '%s'@'%s' identified by '%s'",
		m.Option(METHOD), m.Option(TARGET), m.Option(aaa.USERNAME), m.Option(tcp.HOST), m.Option(aaa.PASSWORD))).SetAppend()
}
func (s Grant) Revoke(m *ice.Message, arg ...string) {
	_sql_exec(m, s.meta(m, m.Option(SESSION), MYSQL), kit.Format("revoke %s on %s from '%s'@'%s'",
		m.Option(METHOD), m.Option(TARGET), m.Option(aaa.USERNAME), m.Option(tcp.HOST))).SetAppend()
}
func (s Grant) Drop(m *ice.Message, arg ...string) {
	_sql_exec(m, s.meta(m, m.Option(SESSION), MYSQL), kit.Format("drop user '%s'@'%s'",
		m.Option(aaa.USERNAME), m.Option(tcp.HOST))).SetAppend()
}
func (s Grant) Remove(m *ice.Message, arg ...string) {
	m.Cmdy(s.Client, s.Remove, arg)
}
func (s Grant) List(m *ice.Message, arg ...string) {
	if len(arg) < 1 || arg[0] == "" { // 连接列表
		m.Cmdy(s.Client)
		return
	}

	_sql_query(m, s.meta(m, arg[0], MYSQL), kit.Format("select User,Host from user")).ToLowerAppend().RenameAppend("user", aaa.USERNAME).Tables(func(value ice.Maps) {
		msg := _sql_query(m.Spawn(), s.meta(m, arg[0], MYSQL), kit.Format("show grants for '%s'@'%s'", value[aaa.USERNAME], value[tcp.HOST]))
		m.Push("stm", msg.Append(""))
	}).Sort("username,host").PushAction(s.Revoke, s.Drop).Action(s.Grants)
}
func init() { ice.CodeModCmd(Grant{}) }
