package client

import (
	"strings"

	"shylinux.com/x/ice"
	kit "shylinux.com/x/toolkits"
)

type cleanup struct {
	query
	list string `name:"list sess database auto" help:"清理"`
}

func (s cleanup) Cleanup(m *ice.Message, arg ...string) {
	m.Cmd(s.query, arg[0]).Table(func(value ice.Maps) {
		if kit.IsIn(value[DATABASE], "information_schema", "performance_schema", "mysql") {
			return
		}
		s.open(m, arg[0], value[DATABASE], func(db *Driver) {
			msg := db.Query(m.Spawn(), "show tables")
			msg.RenameAppend(kit.Select("", msg.Appendv(ice.MSG_APPEND), 0), TABLE).Table(func(val ice.Maps) {
				db.Exec(m, kit.Format("delete from %s where deleted_at IS NOT NULL", val[TABLE]))
			})
		})
	})
}
func (s cleanup) List(m *ice.Message, arg ...string) {
	if len(arg) == 0 {
		m.Cmdy(s.query, arg)
	} else if len(arg) == 1 {
		stats := map[string]int{}
		m.Cmd(s.query, arg[0]).Table(func(value ice.Maps) {
			if kit.IsIn(value[DATABASE], "information_schema", "performance_schema", "mysql") {
				return
			}
			s.open(m, arg[0], value[DATABASE], func(db *Driver) {
				msg := db.Query(m.Spawn(), "show tables")
				table, place, data, delete := msg.Length(), "", 0, 0
				msg.RenameAppend(kit.Select("", msg.Appendv(ice.MSG_APPEND), 0), TABLE).Table(func(val ice.Maps) {
					if strings.HasPrefix(val[TABLE], "user_") {
						place = strings.TrimPrefix(val[TABLE], "user_")
					}
					msg := db.Query(m.Spawn(), kit.Format("select count(*) from %s", val[TABLE]))
					data += kit.Int(msg.Append(msg.Appendv(ice.MSG_APPEND)[0]))
					msg = db.Query(m.Spawn(), kit.Format("select count(*) from %s where deleted_at IS NOT NULL", val[TABLE]))
					delete += kit.Int(msg.Append(msg.Appendv(ice.MSG_APPEND)[0]))
				})
				if place == "" {
					return
				}
				m.Push(DATABASE, value[DATABASE])
				m.Push(TABLE, table)
				m.Push("data", data)
				m.Push("deleted", delete)
				stats[DATABASE]++
				stats[TABLE] += table
				stats["data"] += data
				stats["deleted"] += delete
				msg = db.Query(m.Spawn(), kit.Format("select count(*) from %s", place))
				m.Push("place", kit.Int(msg.Append(msg.Appendv(ice.MSG_APPEND)[0])))
				stats["place"] += kit.Int(msg.Append(msg.Appendv(ice.MSG_APPEND)[0]))
				msg = db.Query(m.Spawn(), kit.Format("select count(distinct(user_uid)) from %s", "user_"+place))
				m.Push("user", kit.Int(msg.Append(msg.Appendv(ice.MSG_APPEND)[0])))
				stats["user"] += kit.Int(msg.Append(msg.Appendv(ice.MSG_APPEND)[0]))
			})
		})
		for k, v := range stats {
			m.Push(k, v)
		}
		m.Action(s.Cleanup)
		m.SortIntR("data")
	} else {
		s.open(m, arg[0], arg[1], func(db *Driver) {
			msg := db.Query(m.Spawn(), "show tables")
			msg.RenameAppend(kit.Select("", msg.Appendv(ice.MSG_APPEND), 0), TABLE).Table(func(val ice.Maps) {
				msg := db.Query(m.Spawn(), kit.Format("select count(*) from %s", val[TABLE]))
				data := kit.Int(msg.Append(msg.Appendv(ice.MSG_APPEND)[0]))
				msg = db.Query(m.Spawn(), kit.Format("select count(*) from %s where deleted_at IS NOT NULL", val[TABLE]))
				delete := kit.Int(msg.Append(msg.Appendv(ice.MSG_APPEND)[0]))
				m.Push(TABLE, val[TABLE])
				m.Push("data", data)
				m.Push("deleted", delete)
			})
		})
		m.SortIntR("data")
	}
}

func init() { ice.CodeModCmd(cleanup{}) }
