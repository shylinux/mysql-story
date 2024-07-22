package db

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
)

type Database struct {
	ice.Hash
	Driver Driver
	driver string `data:"sqlite"`
	short  string `data:"index"`
	field  string `data:"time,index,driver"`
	list   string `name:"list index auto"`
}

func (s Database) Migrate(m *ice.Message, arg ...string) {
	driver := kit.Select(m.Config(DRIVER), arg, 0)
	mdb.HashSelectValue(m.Message, func(value ice.Map) {
		db := s.Driver.open(m, kit.Select(driver, value[DRIVER]))
		db.AutoMigrate(mdb.Confv(m.Message, value[ctx.INDEX], kit.Keym(MODEL), value[mdb.TARGET]))
		mdb.Confv(m.Message, value[ctx.INDEX], kit.Keym(DB), db)
	})
}
func (s Database) List(m *ice.Message, arg ...string) {
	s.Hash.List(m, arg...).Action(s.Migrate)
	m.StatusTimeCount(m.ConfigSimple(DRIVER))
}
func init() { ice.Cmd(prefixKey(), Database{}) }
