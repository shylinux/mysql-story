package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
)

const (
	DRIVER = "driver"
	DB     = "db"
)

type driver struct {
	ice.Hash
	short string `data:"name"`
	field string `data:"time,name,index"`
	list  string `name:"list name auto" help:"驱动"`
}

func (s driver) Select(m *ice.Message, arg ...string) {
	m.Optionv(mdb.TARGET, s.Hash.Target(m, arg[0], nil))
}

func init() { ice.Cmd(prefixKey(), driver{}) }

type Driver struct {
	driver
	beforeMigrate string `name:"beforeMigrate" event:"web.code.db.migrate.before"`
	afterMigrate  string `name:"afterMigrate" event:"web.code.db.migrate.after"`
}

func (s Driver) BeforeMigrate(m *ice.Message, arg ...string) {}
func (s Driver) AfterMigrate(m *ice.Message, arg ...string)  {}

func init() { ice.Cmd(prefixKey(), Driver{}) }

type Dialector interface{ gorm.Dialector }

func (s Driver) Register(m *ice.Message, cb func() Dialector, arg ...string) {
	var err error
	var db *gorm.DB
	m.Cmd(s, s.Create, mdb.NAME, kit.Select(m.CommandKey(), arg, 0), ctx.INDEX, m.PrefixKey(), kit.Dict(mdb.TARGET, func() *gorm.DB {
		defer m.Lock()()
		kit.If(db == nil, func() {
			db, err = gorm.Open(cb(), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
			// db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4")
			m.Warn(err)
		})
		return db
	}))
}
func (s Driver) Target(m *ice.Message, d string) *gorm.DB {
	return m.Cmd(s, s.Select, d).Optionv(mdb.TARGET).(func() *gorm.DB)()
}
