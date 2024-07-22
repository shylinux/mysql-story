package db

import (
	"gorm.io/gorm"
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
)

const (
	DRIVER = "driver"
	DSN    = "dsn"
)

type driver struct {
	ice.Hash
	short string `data:"driver"`
	field string `data:"time,driver,index"`
	list  string `name:"list driver auto"`
}

type Dialector interface{ gorm.Dialector }

func (s driver) Init(m *ice.Message, cb func() Dialector) {
	var err error
	var db *gorm.DB
	m.Cmd(s, mdb.CREATE, DRIVER, m.CommandKey(), ctx.INDEX, m.PrefixKey(), kit.Dict(mdb.TARGET, func() *gorm.DB {
		defer m.Lock(m.PrefixKey())()
		if db != nil {
			return db
		}
		db, err = gorm.Open(cb())
		m.Warn(err)
		return db
	}))
}
func (s driver) Select(m *ice.Message, arg ...string) {
	m.Optionv(mdb.TARGET, mdb.HashSelectTarget(m.Message, kit.Hashs(arg[0]), nil))
}

func init() { ice.Cmd(prefixKey(), driver{}) }

type Driver struct{ driver }

func openDB(m *ice.Message, d string) *gorm.DB {
	msg := m.Cmd(driver{}, mdb.SELECT, d)
	return msg.Optionv(mdb.TARGET).(func() *gorm.DB)()
}
