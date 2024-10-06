package client

import (
	_sql "database/sql"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
	// _ "shylinux.com/x/go-sql-mysql"
)

type Driver struct {
	dsn string
	*_sql.DB
}

func Open(m *ice.Message, driver, dsn string, cb func(*Driver)) {
	if db, e := _sql.Open(driver, dsn); !m.Warn(e) {
		defer db.Close()
		cb(&Driver{DB: db})
	}
}
func (s Driver) Exec(m *ice.Message, stm string, arg ...ice.Any) {
	m.Logs(mdb.MODIFY, "dsn", s.dsn, "stm", stm, "arg", arg)
	m.Push(mdb.TIME, m.Time()).Push("stm", stm)
	if res, err := s.DB.Exec(stm, arg...); m.Warn(err) {
		m.Push("err", kit.Format(err))
		m.Push("lastInsertId", "")
		m.Push("rowsAffected", "")
	} else {
		m.Push("err", "")
		m.Push("lastInsertId", kit.Ignore(res.LastInsertId()))
		m.Push("rowsAffected", kit.Ignore(res.RowsAffected()))
	}
}
func (s Driver) Query(m *ice.Message, stm string, arg ...ice.Any) *ice.Message {
	m.Logs(mdb.SELECT, "dsn", s.dsn, "stm", stm, "arg", arg)
	if rows, err := s.DB.Query(stm, arg...); !m.Warn(err) {
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
				switch v = *(v.(*ice.Any)); v := v.(type) {
				case nil:
					m.Push(head[i], "")
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
	return s.Query(m.Spawn(), kit.Format("select count(*) as total from %s %s", kit.Keys(arg[1], arg[2]), where)).Append(mdb.TOTAL)
}
