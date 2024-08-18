package db

import (
	"errors"
	"reflect"
	"strings"
	"time"
	"unicode"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	NAME       = "name"
	TYPE       = "type"
	ROLE       = "role"
	STATUS     = "status"
	AVATAR     = "avatar"
	ADDRESS    = "address"
)

type Time time.Time

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
	ice.Hash
	database database
	// export        string `data:"true"`
	beforeMigrate string `name:"beforeMigrate" event:"web.code.db.migrate.before"`
	afterMigrate  string `name:"afterMigrate" event:"web.code.db.migrate.after"`
	create        string `name:"create name*"`
	rename        string `name:"rename name*"`
	list          string `name:"list id auto"`
	find          string `name:"find uid*"`
}

func (s Table) BeforeMigrate(m *ice.Message, arg ...string) {
	s.database.Register(m)
}
func (s Table) AfterMigrate(m *ice.Message, arg ...string) {
}
func (s Table) Inputs(m *ice.Message, arg ...string) {
	if strings.HasSuffix(arg[0], "_uid") {
		s.Fields(m, UID, NAME)
		m.Cmdy(m.Prefix(strings.TrimSuffix(arg[0], "_uid"))).RenameAppend(UID, arg[0])
		m.DisplayInputKeyNameIconTitle()
	}
}

func (s Table) Open(m *ice.Message) *gorm.DB {
	s.database.OnceMigrate(m)
	db, ok := m.Optionv(DB).(*gorm.DB)
	kit.If(!ok, func() { db, ok = m.Configv(DB).(*gorm.DB) })
	return db.Model(m.Configv(MODEL)).WithContext(m)
}
func (s Table) OpenID(m *ice.Message, id string) *gorm.DB {
	return s.Open(m).Where("id = ?", id)
}
func (s Table) Create(m *ice.Message, arg ...string) {
	for i := 0; i < len(arg); i += 2 {
		if strings.HasSuffix(arg[i], "_time") {
			t, _ := time.ParseInLocation("2006-01-02 15:04:05", arg[i+1], time.Local)
			arg[i+1] = t.UTC().Format("2006-01-02 15:04:05")
		}
	}
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
	m.ProcessRefresh()
}
func (s Table) Modify(m *ice.Message, arg ...string) {
	m.Warn(s.OpenID(m, m.Option(mdb.ID)).Updates(ice.Map{UPDATED_AT: s.now(m), OPERATOR: m.Option(ice.MSG_USERNAME), arg[0]: arg[1]}).Error)
}
func (s Table) Remove(m *ice.Message, arg ...string) {
	m.Warn(s.OpenID(m, m.Option(mdb.ID)).Updates(kit.Dict(DELETED_AT, s.now(m), arg)).Error)
}
func (s Table) List(m *ice.Message, arg ...string) *ice.Message {
	if len(arg) == 0 {
		m.OptionDefault(mdb.ORDER, "id desc")
		s.Show(m, s.Open(m)).PushAction(s.Remove).Action(s.Create)
	} else {
		m.FieldsSetDetail()
		s.Show(m, s.OpenID(m, arg[0])).PushAction(s.Remove)
	}
	return m
}
func (s Table) Find(m *ice.Message, arg ...string) {
	s.Select(m, arg...)
}
func (s Table) SelectJoin(m *ice.Message, target ice.Any, arg ...string) *ice.Message {
	kit.If(len(arg) == 0, func() { arg = append(arg, NAME) })
	model := s.ToLower(kit.TypeName(target))
	s.ClearOption(m)
	list := []string{}
	m.Table(func(value ice.Maps) { list = kit.AddUniq(list, value[model+"_uid"]) })
	users := m.CmdMap(target, s.SelectList, UID, list, UID)
	m.Table(func(value ice.Maps) {
		user := users[value[model+"_uid"]]
		kit.For(arg, func(k string) {
			if kit.HasSuffix(k, "_uid") {
				m.Push(k, user[k])
			} else {
				m.Push(model+"_"+k, user[k])
			}
		})
	})
	return m
}
func (s Table) SelectList(m *ice.Message, arg ...string) {
	s.Select(m, kit.Format(`%s in ("%v")`, arg[0], kit.Join(arg[1:], `","`)))
}
func (s Table) SelectForUpdate(m *ice.Message, arg ...string) *ice.Message {
	return s.Select(m.Options("query_option", "FOR UPDATE"), arg...)
}
func (s Table) SelectDetail(m *ice.Message, arg ...string) *ice.Message {
	return s.Select(m.FieldsSetDetail(), arg...).Action()
}
func (s Table) Select(m *ice.Message, arg ...string) *ice.Message {
	db := s.Open(m)
	switch table := m.Optionv(mdb.TABLE).(type) {
	case []ice.Any:
		kit.For(table, func(table ice.Any) { db = db.Joins(s.LeftJoin(table)) })
	case []string:
		kit.For(table, func(table string) { db = db.Joins(s.LeftJoin(table)) })
	case string:
		db = db.Joins(s.LeftJoin(table))
	case ice.Any:
		db = db.Joins(s.LeftJoin(table))
	case nil:
		m.OptionDefault(mdb.ORDER, "id desc")
	default:
		m.ErrorNotImplement(table)
	}
	if m.Option("query_option") == "FOR UPDATE" {
		db = db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	s.Show(m, s.Where(m, db, arg...))
	if !m.FieldsIsDetail() {
		m.Action(s.Create)
	}
	s.ClearOption(m)
	m.PushAction()
	return m
}

func (s Table) Update(m *ice.Message, data ice.Map, arg ...string) {
	m.Warn(s.Where(m, s.Open(m), arg...).Updates(data).Error)
	m.ProcessRefresh()
}
func (s Table) Rename(m *ice.Message, arg ...string) *ice.Message {
	s.Update(m, kit.Dict(NAME, m.Option(NAME)), arg...)
	return m
}
func (s Table) Delete(m *ice.Message, arg ...string) *ice.Message {
	res := s.Where(m, s.Open(m), arg...).Delete(m.Configv(MODEL))
	m.Push(mdb.COUNT, res.RowsAffected).Warn(res.Error)
	m.ProcessRefresh()
	return m
}
func (s Table) Transaction(m *ice.Message, cb func()) {
	s.Open(m).Transaction(func(tx *gorm.DB) error {
		m.Optionv(DB, tx)
		if cb(); m.IsErr() {
			return errors.New(m.Result())
		}
		return nil
	})
	m.Option(DB, "")
	s.ClearOption(m)
}
func (s Table) AddCount(m *ice.Message, arg ...string) {
	s.Exec(m, kit.Format("Update %s set %s = %s + %s where uid = '%s'",
		s.TableName(kit.TypeName(m.Configv(MODEL))), arg[0], arg[0], arg[1], arg[2]))
}
func (s Table) Exec(m *ice.Message, arg ...string) {
	db := s.Open(m)
	db.Exec(arg[0])
}
func (s Table) Show(m *ice.Message, db *gorm.DB) *ice.Message {
	fields := kit.Simple(m.Optionv(mdb.SELECT))
	kit.If(len(fields) > 0, func() { db = db.Select(fields) })
	kit.If(m.Option(mdb.ORDER), func() {
		db = db.Order(kit.Join(kit.Simple(m.Optionv(mdb.ORDER)), ","))
	})
	rows, err := db.Offset(kit.Int(m.OptionDefault(mdb.OFFSET, "0"))).Limit(kit.Int(m.OptionDefault(mdb.LIMIT, "30"))).Rows()
	if m.Warn(err) {
		return m
	}
	defer rows.Close()
	head, err := rows.Columns()
	if m.Warn(err) {
		return m
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
			case nil:
				m.Push(head[i], "")
			case []byte:
				m.Push(head[i], string(v))
			default:
				if v != nil && (kit.IsIn(head[i], CREATED_AT, UPDATED_AT) || kit.HasSuffix(head[i], "_time")) {
					if t, e := time.Parse("2006-01-02 15:04:05 -0700 UTC", kit.Format("%v", v)); e == nil {
						v = t.Local().Format("2006-01-02 15:04:05")
					}
				}
				m.Push(head[i], kit.Format("%v", v))
			}
		}
	}
	return m
}
func (s Table) ClearOption(m *ice.Message, arg ...string) *ice.Message {
	m.Optionv(mdb.TABLE, []string{})
	m.Optionv(mdb.SELECT, []string{})
	m.Optionv(mdb.ORDER, []string{})
	m.Optionv("query_option", []string{})
	return m
}
func (s Table) Tables(m *ice.Message, target ...ice.Any) Table {
	m.Optionv(mdb.TABLE, target...)
	return s
}
func (s Table) FieldsWithCreatedAT(m *ice.Message, target ice.Any, arg ...ice.Any) Table {
	if target == nil || target == "" {
		target = m.Configv(MODEL)
	}
	s.Fields(m, append([]ice.Any{
		s.AS(s.Key(target, CREATED_AT), CREATED_AT),
		s.AS(s.Key(target, UID), UID),
	}, arg...)...).Orders(m, s.Desc(CREATED_AT))
	return s
}
func (s Table) Fields(m *ice.Message, arg ...ice.Any) Table {
	for i, v := range arg {
		switch v := v.(type) {
		case string:
			kit.For([]string{NAME, TYPE, ROLE, AVATAR, STATUS, ADDRESS}, func(suffix string) {
				if !strings.Contains(v, " ") && strings.HasSuffix(v, "_"+suffix) {
					arg[i] = s.TableName(strings.TrimSuffix(v, "_"+suffix)) + "." + suffix + " AS " + v
				}
			})
		}
	}
	m.Optionv(mdb.SELECT, arg...)
	return s
}
func (s Table) Orders(m *ice.Message, arg ...ice.Any) Table {
	m.Optionv(mdb.ORDER, arg...)
	return s
}
func (s Table) Limit(m *ice.Message, limit int) Table {
	m.Option(mdb.LIMIT, limit)
	return s
}
func (s Table) Where(m *ice.Message, db *gorm.DB, arg ...string) *gorm.DB {
	if len(arg) == 0 {
		return db
	}
	if len(arg) == 1 || strings.Contains(arg[0], "?") {
		db = db.Where(arg[0], kit.TransArgs(arg[1:])...)
	} else {
		params := kit.Dict()
		kit.For(arg, func(k, v string) { params[k] = v })
		db = db.Where(params)
	}
	return db
}
func (s Table) LeftJoin(target ice.Any) string {
	model, models := "", ""
	switch target := target.(type) {
	case string:
		if strings.Contains(target, " ") {
			return target
		}
		model = target
	default:
		model = kit.TypeName(target)
		model = s.ToLower(model)
	}
	models = s.TableName(model)
	return kit.Format("left join %s on %s_uid = %s.uid", models, model, models)
}
func (s Table) ToLower(model string) string {
	list, begin, last := []string{}, 0, false
	for i, v := range model {
		if i == len(model)-1 {
			list = append(list, strings.ToLower(model[begin:]))
		} else if unicode.IsUpper(v) && last {
			list, begin = append(list, strings.ToLower(model[begin:i])), i
		}
		last = unicode.IsLower(v)
	}
	model = kit.Join(list, "_")
	return model
}
func (s Table) TableName(model string) string {
	model = s.ToLower(model)
	if kit.HasSuffix(model, "y") {
		model = model[:len(model)-1] + "ies"
	} else if kit.HasSuffix(model, "s") {
		model = model + "es"
	} else {
		model = model + "s"
	}
	return model
}
func (s Table) now(m *ice.Message) string {
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}
func (s Table) Keys(target ice.Any, k string) string {
	return s.ToLower(kit.TypeName(target)) + "_" + k
}
func (s Table) Key(target ice.Any, k string) string {
	return kit.Keys(s.TableName(kit.TypeName(target)), k)
}
func (s Table) AS(from, to string) string {
	return from + " AS " + to
}
func (s Table) Desc(k string) string {
	return k + " DESC"
}

func prefixKey() string { return kit.Keys("web.code", kit.PathName(-1), kit.FileName(-1)) }
