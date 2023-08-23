package sqlite

import (
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/mysql-story/src/client"
	kit "shylinux.com/x/toolkits"
)

type Client struct {
	list string `name:"list path tbl_name id auto" help:"数据库" icon:"sqlite.png"`
}

func (s Client) List(m *ice.Message, arg ...string) {
	if len(arg) == 0 || strings.HasSuffix(arg[0], nfs.PS) {
		m.Cmdy(nfs.DIR, arg)
	} else {
		client.Open(m, SQLITE3, arg[0], func(db *client.Driver) {
			if len(arg) == 1 {
				db.Query(m, kit.Format("SELECT * FROM sqlite_master WHERE type = 'table'"))
			} else if len(arg) == 2 {
				db.Query(m, kit.Format("SELECT * from %s LIMIT 100", arg[1]))
			} else {
				m.OptionFields(mdb.DETAIL)
				db.Query(m, kit.Format("SELECT * from %s WHERE id = %s", arg[1], arg[2]))
			}
		})
	}
}

func init() { ice.CodeCtxCmd(Client{}) }
