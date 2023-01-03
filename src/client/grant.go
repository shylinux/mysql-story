package client

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

const (
	METHOD = "method"
	TARGET = "target"
)

type Grant struct {
	Client
	grants string `name:"grants sess method=all target=*.* host=% username password" help:"授权"`
	revoke string `name:"revoke sess method=all target=*.* host username" help:"撤销"`
	drop   string `name:"drop" help:"删除"`
	list   string `name:"list sess auto" help:"权限"`
}

func (s Grant) Grants(m *ice.Message, arg ...string) {
	_sql_exec(m.Spawn(), s.meta(m, m.Option(aaa.SESS), MYSQL), kit.Format("grant %s on %s to '%s'@'%s' identified by '%s'",
		m.Option(METHOD), m.Option(TARGET), m.Option(aaa.USERNAME), m.Option(tcp.HOST), m.Option(aaa.PASSWORD)))
}
func (s Grant) Revoke(m *ice.Message, arg ...string) {
	_sql_exec(m.Spawn(), s.meta(m, m.Option(aaa.SESS), MYSQL), kit.Format("revoke %s on %s from '%s'@'%s'",
		m.Option(METHOD), m.Option(TARGET), m.Option(aaa.USERNAME), m.Option(tcp.HOST)))
}
func (s Grant) Drop(m *ice.Message, arg ...string) {
	_sql_exec(m.Spawn(), s.meta(m, m.Option(aaa.SESS), MYSQL), kit.Format("drop user '%s'@'%s'",
		m.Option(aaa.USERNAME), m.Option(tcp.HOST)))
}
func (s Grant) List(m *ice.Message, arg ...string) {
	if len(arg) < 1 || arg[0] == "" {
		m.Cmdy(s.Client)
		return
	}
	_sql_query(m, s.meta(m, arg[0], MYSQL), kit.Format("select User,Host from user")).ToLowerAppend().RenameAppend(aaa.USER, aaa.USERNAME).Tables(func(value ice.Maps) {
		msg := _sql_query(m.Spawn(), s.meta(m, arg[0], MYSQL), kit.Format("show grants for '%s'@'%s'", value[aaa.USERNAME], value[tcp.HOST]))
		m.Push("stm", msg.Append(""))
	}).Sort("username,host").PushAction(s.Revoke, s.Drop).Action(s.Grants)
}
func init() { ice.CodeModCmd(Grant{}) }
