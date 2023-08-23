package client

import (
	_sql "database/sql"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
)

type Driver struct {
	*_sql.DB
}

func Open(m *ice.Message, driver, dsn string, cb func(*Driver)) *ice.Message {
	if db, e := _sql.Open(driver, dsn); m.Assert(e) {
		defer db.Close()
		cb(&Driver{db})
	}
	return m
}

func (s Driver) Exec(m *ice.Message, stm string, arg ...ice.Any) *ice.Message {
	// m.Logs(mdb.MODIFY, "dsn", dsn, "stm", stm, "arg", arg)
	m.Push(mdb.TIME, m.Time())
	if res, err := s.DB.Exec(stm, arg...); err != nil {
		m.Push("", kit.UnMarshal(kit.Format(err)))
	} else {
		if i, e := res.LastInsertId(); e == nil {
			m.Push("lastInsertId", i)
		}
		if i, e := res.RowsAffected(); e == nil {
			m.Push("rowsAffected", i)
		}
	}
	return m
}
func (s Driver) Query(m *ice.Message, stm string, arg ...ice.Any) *ice.Message {
	// m.Logs(mdb.SELECT, "dsn", dsn, "stm", stm, "arg", arg)
	if rows, err := s.DB.Query(stm, arg...); m.Assert(err) {
		head, err := rows.Columns()
		m.Assert(err)
		var data ice.List
		for _, _ = range head {
			var item ice.Any
			data = append(data, &item)
		}
		defer m.StatusTimeCount()
		for rows.Next() {
			rows.Scan(data...)
			for i, v := range data {
				v = *(v.(*ice.Any))
				switch v := v.(type) {
				case []byte:
					m.Push(head[i], string(v))
				default:
					m.Push(head[i], kit.Format("%v", v))
				}
			}
		}
	}
	return m
}
func (s Driver) Total(m *ice.Message, where string, arg ...string) string {
	if len(arg) > 2 {
		msg := s.Query(m.Spawn(), kit.Format("select count(*) as total from %s %s", kit.Keys(arg[1], arg[2]), where))
		return msg.Append(mdb.TOTAL)
	}
	return ""
}
