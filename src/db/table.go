package db

import (
	"gorm.io/gorm"
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
	"sync"
)

const (
	MODEL = "model"
)

type Model struct {
	gorm.Model
	Creator  string
	Operator string
}

type Table struct {
	list string `name:"list id auto"`
}

func (s Table) Init(m *ice.Message, t ice.Any) {
	m.Cmd(Database{}, mdb.CREATE, ctx.INDEX, m.PrefixKey(), DRIVER, m.Config(DRIVER), kit.Dict(mdb.TARGET, t))
	m.TransInput("created_at", "创建时间", "deleted_at", "删除时间", "updated_at", "更新时间", "creator", "创建人", "operator", "操作人")
}

var once = &sync.Once{}

func (s Table) Open(m *ice.Message) *gorm.DB {
	once.Do(func() { m.Cmd(Database{}, "migrate") })
	// return m.Configv("db").(*gorm.DB)
	return m.Configv("db").(*gorm.DB).Model(m.Configv("model"))
}
func (s Table) Create(m *ice.Message, arg ...string) {
	res := s.Open(m).Create(kit.Dict("created_at", m.Time(), "creator", m.Option(ice.MSG_USERNAME), arg))
	m.Warn(res.Error)
}
func (s Table) Modify(m *ice.Message, arg ...string) {
	res := s.Open(m).Where("id = ?", m.Option(mdb.ID)).Updates(kit.Dict("updated_at", m.Time(), "operator", m.Option(ice.MSG_USERNAME), arg))
	m.Warn(res.Error)
}
func (s Table) Remove(m *ice.Message, arg ...string) {
	res := s.Open(m).Where("id = ?", m.Option(mdb.ID)).Updates(kit.Dict("deleted_at", m.Time(), arg))
	m.Warn(res.Error)
}
func (s Table) Select(m *ice.Message, stm string, arg ...ice.Any) {
	s.Show(m, s.Open(m).Where(stm, arg...))
}
func (s Table) List(m *ice.Message, arg ...string) {
	db := s.Open(m)
	if len(arg) > 0 {
		m.FieldsSetDetail()
		db = db.Where("id = ?", arg[0])
	} else {
		defer m.Action(s.Create)
	}
	s.Show(m, db)
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
			if head[i] == "deleted_at" {
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

func prefixKey() string {
	return kit.Keys("web.code", kit.PathName(-1), kit.FileName(-1))
}
