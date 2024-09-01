package db

import (
	"sync"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
)

type database struct {
	ice.Hash
	models models
	Models Models
	Driver Driver
	driver string `data:"mysql"`
	short  string `data:"index"`
	field  string `data:"time,index,models,domain,driver"`
	list   string `name:"list index auto" help:"存储"`
}

func (s database) Exit(m *ice.Message, arg ...string) {
	m.Confv(m.PrefixKey(), mdb.HASH, "")
}
func (s database) Migrate(m *ice.Message, arg ...string) {
	driver := kit.Select(m.Config(DRIVER), arg, 0)
	mdb.HashSelectValue(m.Message, func(value ice.Map) {
		db := s.Driver.Target(m, kit.Select(driver, value[DRIVER]), kit.Format(value[DOMAIN]))
		m.Warn(db.AutoMigrate(mdb.Confv(m.Message, value[ctx.INDEX], kit.Keym(MODEL), value[mdb.TARGET])), value[ctx.INDEX])
		mdb.Confv(m.Message, value[ctx.INDEX], kit.Keym(DB), db)
	})
}
func (s database) List(m *ice.Message, arg ...string) {
	s.Hash.List(m, arg...).Action(s.Migrate)
	m.StatusTimeCount(m.ConfigSimple(DRIVER))
}
func init() { ice.Cmd(prefixKey(), database{}) }

func (s database) Register(m *ice.Message) {
	models := kit.Select(m.CommandKey(), m.Config(MODELS))
	domain := kit.Select(kit.Select("", kit.Split(m.PrefixKey(), "."), 2), m.Config(DOMAIN))
	target := s.Models.Target(m, kit.Keys(domain, models))
	m.WarnNotFound(target == nil, kit.Keys(domain, models))
	m.Cmd(s, s.Create, ctx.INDEX, m.PrefixKey(), MODELS, kit.Keys(domain, models), DOMAIN, domain, DRIVER, m.Config(DRIVER), kit.Dict(mdb.TARGET, target))
}
func (s database) OnceMigrate(m *ice.Message) {
	once.Do(func() {
		m.Spawn(ice.Maps{ice.MSG_DAEMON: ""}).Cmd(s.models, s.models.AutoCreate)
		m.Event("web.code.db.migrate.before")
		defer m.GoSleep("30ms", func() { m.Event("web.code.db.migrate.after") })
		m.Cmd(s, s.Migrate)
	})
}

var once = &sync.Once{}
