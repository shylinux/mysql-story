package db

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/web"
	"shylinux.com/x/icebergs/base/web/html"
	"shylinux.com/x/icebergs/misc/xterm"
	kit "shylinux.com/x/toolkits"
)

type models struct {
	ice.Hash
	short string `data:"name"`
	field string `data:"time,name,index"`
	cmds  string `data:"/opt/10001/bin/mysql -S /opt/10001/data/mysqld.socket -u root -proot"`
	path  string `data:"/opt/10001/"`
	list  string `name:"list name auto" help:"模型"`
}

func (s models) Exit(m *ice.Message, arg ...string) {
	m.Confv(m.PrefixKey(), mdb.HASH, "")
}
func (s models) Select(m *ice.Message, arg ...string) {
	m.Optionv(mdb.TARGET, s.Hash.Target(m, arg[0], nil))
}
func (s models) List(m *ice.Message, arg ...string) {
	s.Hash.List(m, arg...).Action(s.Xterm, s.AutoCreate)
}
func (s models) AutoCreate(m *ice.Message, arg ...string) {
	cmds := []string{}
	list := map[string]bool{}
	m.Cmd("").Table(func(value ice.Maps) {
		if db := kit.Split(value[mdb.NAME], ".")[0]; !list[db] {
			cmds = append(cmds, kit.Format("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;", db))
			list[db] = true
		}
	})
	cmds = append(cmds, "quit")
	args := kit.Split(m.Config(ctx.CMDS))
	m.Info("auto create %v", args)
	m.Info("auto create %v", cmds)
	if p, e := xterm.Command(m.Message, m.Config(nfs.PATH), args[0], args[1:]...); !m.Warn(e) {
		xterm.PushShell(m.Message, p, cmds, func(res string) {
			m.Option(ice.MSG_TITLE, m.ActionKey())
			web.PushNoticeGrow(m.Options(ctx.DISPLAY, html.PLUGIN_XTERM, ice.MSG_COUNT, "0", ice.MSG_DEBUG, ice.FALSE, ice.LOG_DISABLE, ice.TRUE).Message, res)
			kit.If(kit.HasPrefix(res, "quit"), func() { p.Close() })
		})
		m.ProcessHold()
	}
}
func (s models) Xterm(m *ice.Message, arg ...string) {
	m.ProcessXterm("AutoCreate", []string{mdb.TYPE, m.Config(ctx.CMDS), nfs.PATH, m.Config(nfs.PATH)}, arg...)
}

func init() { ice.Cmd(prefixKey(), models{}) }

type Models struct {
	models
	Database string
	Tables   []ice.Any
}

func (s Models) Init(m *ice.Message, arg ...string) {
	s.Hash.Init(m, arg...)
	if s.Database != "" {
		s.Register(m, s.Database, s.Tables...)
	}
}
func CmdModels(db string, tables ...ice.Any) {
	ice.Cmd(kit.Keys("web.code.db", kit.ModName(-1), kit.ModPath(-1), MODELS), Models{
		Database: kit.Select(kit.Split(kit.ModPath(-1), "./", "./")[0], db), Tables: tables,
	})
}
func (s Models) Register(m *ice.Message, domain string, target ...ice.Any) {
	kit.For(target, func(target ice.Any) {
		m.Cmd(s.models, s.Create, mdb.NAME, kit.Keys(domain, kit.TypeName(target)), ctx.INDEX, m.PrefixKey(), kit.Dict(mdb.TARGET, target))
	})
}
func (s Models) Target(m *ice.Message, name string) ice.Any {
	return m.Cmd(s.models, s.Select, name).Optionv(mdb.TARGET)
}
