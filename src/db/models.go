package db

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
)

type models struct {
	ice.Hash
	short string `data:"name"`
	field string `data:"time,name,index"`
	list  string `name:"list name auto" help:"模型"`
}

func (s models) Exit(m *ice.Message, arg ...string) {
	m.Confv(m.PrefixKey(), mdb.HASH, "")
}
func (s models) Select(m *ice.Message, arg ...string) {
	m.Optionv(mdb.TARGET, s.Hash.Target(m, arg[0], nil))
}
func (s models) List(m *ice.Message, arg ...string) {
	s.Hash.List(m, arg...)
	list := map[string]bool{}
	m.Table(func(value ice.Maps) {
		if db := kit.Split(value[mdb.NAME], ".")[0]; !list[db] {
			m.EchoScript(kit.Format("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;", db))
			list[db] = true
		}
	})
	m.EchoScript("./bin/mysql -S ./data/mysqld.socket -u root -p")
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
	ice.Cmd(kit.Keys("web.code.db", kit.ModName(-1), kit.ModPath(-1), "models"), Models{
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
