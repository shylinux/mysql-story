package sqlite

import (
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/mysql-story/src/client"
	kit "shylinux.com/x/toolkits"

	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	list string `name:"list path tbl_name id auto" icon:"sqlite.png" help:"存储"`
}

func (s Client) Drop(m *ice.Message, arg ...string) {
	client.Open(m, SQLITE3, m.Option(nfs.PATH), func(db *client.Driver) {
		db.Exec(m, kit.Format("DROP TABLE %s", m.Option("tbl_name")))
		kit.If(m.Append(ice.ERR) == "", func() { m.ProcessRefresh() })
	})
}
func (s Client) List(m *ice.Message, arg ...string) {
	if len(arg) == 0 || strings.HasSuffix(arg[0], nfs.PS) {
		m.Cmdy(nfs.DIR, arg)
	} else {
		client.Open(m, SQLITE3, arg[0], func(db *client.Driver) {
			if len(arg) == 1 {
				db.Query(m, kit.Format("SELECT * FROM sqlite_master WHERE type = 'table'"))
				m.PushAction(s.Drop)
			} else if len(arg) == 2 {
				db.Query(m, kit.Format("SELECT * from %s LIMIT 100", arg[1]))
			} else {
				m.FieldsSetDetail()
				db.Query(m, kit.Format("SELECT * from %s WHERE id = %s", arg[1], arg[2]))
			}
		})
	}
}

func init() { ice.CodeCtxCmd(Client{}) }
