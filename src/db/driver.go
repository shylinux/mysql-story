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
	DSN    = "dsn"
	DB     = "db"
)

type driver struct {
	ice.Hash
	short string `data:"driver"`
	field string `data:"time,driver,index"`
	list  string `name:"list driver auto"`
}

type Dialector interface{ gorm.Dialector }

func (s driver) Init(m *ice.Message, cb func() Dialector, arg ...string) {
	var err error
	var db *gorm.DB
	m.Cmd(s, mdb.CREATE, DRIVER, kit.Select(m.CommandKey(), arg, 0), ctx.INDEX, m.PrefixKey(), kit.Dict(mdb.TARGET, func() *gorm.DB {
		defer m.Lock(m.PrefixKey())()
		kit.If(db == nil, func() {
			db, err = gorm.Open(cb(), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
			// db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4")以
			m.Warn(err)
		})
		return db
	}))
}
func (s driver) Select(m *ice.Message, arg ...string) {
	m.Optionv(mdb.TARGET, mdb.HashSelectTarget(m.Message, kit.Hashs(arg[0]), nil))
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

func (s Driver) open(m *ice.Message, d string) *gorm.DB {
	return m.Cmd(s, s.Select, d).Optionv(mdb.TARGET).(func() *gorm.DB)()
}
