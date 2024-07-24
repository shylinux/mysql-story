package db

import (
	"reflect"
	"time"

	"gorm.io/gorm"
	"shylinux.com/x/ice"
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
	UID        = "uid"
)

type Model struct {
	gorm.Model
	Creator  string `gorm:"type:char(32)"`
	Operator string `gorm:"type:char(32)"`
}
type ModelWithUID struct {
	Model
	Uid string `gorm:"type:char(32);uniqueIndex"`
}

type Table struct {
	database      database
	beforeMigrate string `name:"beforeMigrate" event:"web.code.db.migrate.before"`
	afterMigrate  string `name:"afterMigrate" event:"web.code.db.migrate.after"`
	list          string `name:"list id auto"`
}

func (s Table) BeforeMigrate(m *ice.Message, arg ...string) {
	s.database.Register(m)
}
func (s Table) AfterMigrate(m *ice.Message, arg ...string) {
}
func (s Table) Open(m *ice.Message) *gorm.DB {
	s.database.OnceMigrate(m)
	db, ok := m.Optionv(DB).(*gorm.DB)
	kit.If(!ok, func() { db, ok = m.Configv(DB).(*gorm.DB) })
	return db.Model(m.Configv(MODEL))
}
func (s Table) OpenID(m *ice.Message, id string) *gorm.DB {
	return s.Open(m).Where("id = ?", id)
}
func (s Table) now(m *ice.Message) string {
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}
func (s Table) Create(m *ice.Message, arg ...string) {
	data := kit.Dict(CREATED_AT, s.now(m), CREATOR, m.Option(ice.MSG_USERNAME), arg)
	switch model := m.Configv(MODEL).(type) {
	case interface{ OnCreate(ice.Map) }:
		model.OnCreate(data)
	default:
		if data[UID] == nil {
			t := reflect.TypeOf(model)
			kit.If(t.Kind() == reflect.Ptr, func() { t = t.Elem() })
			if _, ok := t.FieldByName("Uid"); ok {
				data[UID] = kit.HashsUniq()
			}
		}
	}
	if !m.Warn(s.Open(m).Create(data).Error) {
		m.Echo(kit.Select(kit.Format(data[mdb.ID]), data[UID]))
	}
}
func (s Table) Modify(m *ice.Message, arg ...string) {
	m.Warn(s.OpenID(m, m.Option(mdb.ID)).Updates(kit.Dict(UPDATED_AT, s.now(m), OPERATOR, m.Option(ice.MSG_USERNAME), arg)).Error)
}
func (s Table) Remove(m *ice.Message, arg ...string) {
	m.Warn(s.OpenID(m, m.Option(mdb.ID)).Updates(kit.Dict(DELETED_AT, s.now(m), arg)).Error)
}
func (s Table) Select(m *ice.Message, arg ...string) {
	args := kit.List()
	kit.For(arg[1:], func(v string) { args = append(args, v) })
	s.Show(m, s.Open(m).Where(arg[0], args...))
}
func (s Table) List(m *ice.Message, arg ...string) {
	if len(arg) == 0 {
		s.Show(m, s.Open(m))
		m.Action(s.Create)
	} else {
		m.FieldsSetDetail()
		s.Show(m, s.OpenID(m, arg[0]))
	}
}
func (s Table) Show(m *ice.Message, db *gorm.DB) {
	fields := kit.Simple(m.Optionv(mdb.SELECT))
	kit.If(len(fields) > 0, func() { db = db.Select(fields) })
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
				if v != nil && kit.IsIn(head[i], CREATED_AT, UPDATED_AT) {
					if t, e := time.Parse("2006-01-02 15:04:05 -0700 UTC", kit.Format("%v", v)); e == nil {
						v = t.Local().Format("2006-01-02 15:04:05")
					}
				}
				m.Push(head[i], kit.Format("%v", v))
			}
		}
	}
	m.PushAction(s.Remove)
}

func prefixKey() string { return kit.Keys("web.code", kit.PathName(-1), kit.FileName(-1)) }
