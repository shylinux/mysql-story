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
	MYSQL  = "mysql"
)

type Grant struct {
	client
	grants string `name:"grants sess* method*=all target*='*.*' host*=% username* password*"`
	revoke string `name:"revoke sess* method*=all target*='*.*' host username*"`
	list   string `name:"list sess auto"`
}

func (s Grant) Grants(m *ice.Message, arg ...string) {
	s.open(m, m.Option(aaa.SESS), MYSQL, func(db *Driver) {
		db.Exec(m.Spawn(), kit.Format("create user '%s'@'%s' identified by '%s'", m.Option(aaa.USERNAME), m.Option(tcp.HOST), m.Option(aaa.PASSWORD)))
		db.Exec(m.Spawn(), kit.Format("grant %s on %s to '%s'@'%s'", m.Option(METHOD), m.Option(TARGET), m.Option(aaa.USERNAME), m.Option(tcp.HOST)))
	})
}
func (s Grant) Revoke(m *ice.Message, arg ...string) {
	s.open(m, m.Option(aaa.SESS), MYSQL, func(db *Driver) {
		db.Exec(m.Spawn(), kit.Format("revoke %s on %s from '%s'@'%s'", m.Option(METHOD), m.Option(TARGET), m.Option(aaa.USERNAME), m.Option(tcp.HOST)))
	})
}
func (s Grant) Drop(m *ice.Message, arg ...string) {
	s.open(m, m.Option(aaa.SESS), MYSQL, func(db *Driver) {
		db.Exec(m.Spawn(), kit.Format("drop user '%s'@'%s'", m.Option(aaa.USERNAME), m.Option(tcp.HOST)))
	})
}
func (s Grant) List(m *ice.Message, arg ...string) {
	if len(arg) == 0 || arg[0] == "" {
		m.Cmdy(s.client)
		return
	}
	s.open(m, arg[0], MYSQL, func(db *Driver) {
		db.Query(m, kit.Format("select User,Host from user")).ToLowerAppend().RenameAppend(aaa.USER, aaa.USERNAME).Table(func(value ice.Maps) {
			msg := db.Query(m.Spawn(), kit.Format("show grants for '%s'@'%s'", value[aaa.USERNAME], value[tcp.HOST]))
			m.Push("stm", msg.Append(""))
		}).Sort("username,host").PushAction(s.Revoke, s.Drop).Action(s.Grants)
	})
}
func init() { ice.CodeModCmd(Grant{}) }
