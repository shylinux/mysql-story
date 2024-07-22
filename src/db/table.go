package db

import (
	"sync"

	"gorm.io/gorm"
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
)

const (
	MODEL      = "model"
	CREATED_AT = "created_at"
	UPDATED_AT = "updated_at"
	DELETED_AT = "deleted_at"
	CREATOR    = "creator"
	OPERATOR   = "operator"
)

type Model struct {
	gorm.Model
	Creator  string
	Operator string
}

type Table struct {
	models Models
	db     Database
	list   string `name:"list id auto"`

	beforeMigrate string `name:"beforeMigrate" event:"web.code.db.migrate.before"`
	afterMigrate  string `name:"afterMigrate" event:"web.code.db.migrate.after"`
}

func (s Table) BeforeMigrate(m *ice.Message, arg ...string) {
	s.Init(m, s.models.Bind(m, kit.Select(m.CommandKey(), m.Config("models"))))
}
func (s Table) AfterMigrate(m *ice.Message, arg ...string) {

}
func (s Table) Init(m *ice.Message, t ice.Any) *ice.Message {
	m.Cmd(s.db, s.db.Create, ctx.INDEX, m.PrefixKey(), DRIVER, m.Config(DRIVER), kit.Dict(mdb.TARGET, t))
	m.TransInput(CREATED_AT, "创建时间", UPDATED_AT, "更新时间", DELETED_AT, "删除时间", CREATOR, "创建人", OPERATOR, "操作人")
	return m
}

var once = &sync.Once{}

func (s Table) Open(m *ice.Message) *gorm.DB {
	once.Do(func() {
		defer m.Event("web.code.db.migrate.before")("web.code.db.migrate.after")
		m.Cmd(s.db, s.db.Migrate)
	})
	return m.Configv(DB).(*gorm.DB).Model(m.Configv(MODEL))
}
func (s Table) OpenID(m *ice.Message, id string) *gorm.DB {
	return s.Open(m).Where("id = ?", id)
}
func (s Table) Create(m *ice.Message, arg ...string) {
	res := s.Open(m).Create(kit.Dict(CREATED_AT, m.Time(), CREATOR, m.Option(ice.MSG_USERNAME), arg))
	m.Warn(res.Error)
}
func (s Table) Modify(m *ice.Message, arg ...string) {
	res := s.OpenID(m, m.Option(mdb.ID)).Updates(kit.Dict(UPDATED_AT, m.Time(), OPERATOR, m.Option(ice.MSG_USERNAME), arg))
	m.Warn(res.Error)
}
func (s Table) Remove(m *ice.Message, arg ...string) {
	res := s.OpenID(m, m.Option(mdb.ID)).Updates(kit.Dict(DELETED_AT, m.Time(), arg))
	m.Warn(res.Error)
}
func (s Table) Select(m *ice.Message, stm string, arg ...ice.Any) {
	s.Show(m, s.Open(m).Where(stm, arg...))
}
func (s Table) List(m *ice.Message, arg ...string) {
	if len(arg) > 0 {
		m.FieldsSetDetail()
		s.Show(m, s.OpenID(m, arg[0]))
	} else {
		m.Action(s.Create)
		s.Show(m, s.Open(m))
	}
}
func (s Table) Show(m *ice.Message, db *gorm.DB) {
	rows, err := db.Offset(kit.Int(m.OptionDefault(mdb.OFFSET, "0"))).Limit(kit.Int(m.OptionDefault(mdb.LIMIT, "30"))).Rows()
	if m.Warn(err) {
		return
	}
	defer rows.Close()
	head, err := rows.Columns()
	if m.Warn(err) {
		return
	}
	var data ice.List
	for _, _ = range head {
		var item ice.Any
		data = append(data, &item)
	}
	for rows.Next() {
		rows.Scan(data...)
		for i, v := range data {
			if head[i] == DELETED_AT {
				continue
			}
			switch v = *(v.(*ice.Any)); v := v.(type) {
			case []byte:
				m.Push(head[i], string(v))
			default:
				m.Push(head[i], kit.Format("%v", v))
			}
		}
	}
	m.PushAction(s.Remove)
}

func prefixKey() string { return kit.Keys("web.code", kit.PathName(-1), kit.FileName(-1)) }
