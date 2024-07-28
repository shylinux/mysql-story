package db

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
)

type models struct {
	ice.Hash
	short string `data:"name"`
	field string `data:"time,name,index"`
	list  string `name:"list name auto" help:"模型"`
}

func (s models) Select(m *ice.Message, arg ...string) {
	m.Optionv(mdb.TARGET, s.Hash.Target(m, arg[0], nil))
}

func init() { ice.Cmd(prefixKey(), models{}) }

type Models struct{ models }

func init() { ice.Cmd(prefixKey(), Models{}) }

func (s Models) Register(m *ice.Message, domain string, target ...ice.Any) {
	kit.For(target, func(target ice.Any) {
		m.Cmd(s, s.Create, mdb.NAME, kit.Keys(domain, kit.TypeName(target)), ctx.INDEX, m.PrefixKey(), kit.Dict(mdb.TARGET, target))
	})
}
func (s Models) Target(m *ice.Message, name string) ice.Any {
	return m.Cmd(s, s.Select, name).Optionv(mdb.TARGET)
}
