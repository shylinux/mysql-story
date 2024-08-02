package client

import (
	"path"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

type client struct {
	ice.Hash
	short  string `data:"sess"`
	action string `data:"xterm"`
	field  string `data:"time,sess,driver,database,host,port,username,password"`
	create string `name:"create sess*=biz driver*=mysql database*=mysql host*=localhost port*=10001 username*=root password*=root"`
	list   string `name:"list sess auto"`
}

func (s client) Xterm(m *ice.Message, arg ...string) {
	if kit.HasPrefixList(arg, ctx.RUN) {
		m.ProcessXterm("", "", arg...)
		return
	}
	msg := m.Cmd(s, m.Option(aaa.SESS))
	database := kit.Select(msg.Append(DATABASE), m.Option(DATABASE))
	m.ProcessXterm(kit.Format("%s(%s:%s)", m.Option(aaa.SESS), msg.Append(tcp.HOST), msg.Append(tcp.PORT)),
		kit.Format("%s -h %s -P %s -u %s -p%s %s", path.Join(ice.USR_LOCAL_DAEMON, msg.Append(tcp.PORT), "bin/mysql"),
			msg.Append(tcp.HOST), msg.Append(tcp.PORT), msg.Append(aaa.USERNAME), msg.Append(aaa.PASSWORD), database), arg...)

}
func (s client) List(m *ice.Message, arg ...string) {
	s.Hash.List(m, arg...)
	if len(arg) == 1 {
		m.EchoScript(kit.Format("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;", m.Append(DATABASE)))
	}
}

func init() { ice.CodeModCmd(client{}) }

func (s client) open(m *ice.Message, sess string, db string, cb func(*Driver)) {
	msg := m.Cmd(s, sess)
	Open(m, msg.Append(DRIVER), kit.Format("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4", msg.Append(aaa.USERNAME), msg.Append(aaa.PASSWORD), msg.Append(tcp.HOST), msg.Append(tcp.PORT), kit.Select(msg.Append(DATABASE), db)), cb)
}

type Client struct{ client }

func init() { ice.CodeModCmd(Client{}) }
