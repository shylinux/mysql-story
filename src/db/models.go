package db

import (
	"reflect"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
)

type models struct {
	ice.Hash
	short string `data:"model"`
	field string `data:"time,model,index"`
	list  string `name:"list model auto"`
}

func (s models) Init(m *ice.Message, target ...ice.Any) {
	kit.For(target, func(target ice.Any) {
		t := reflect.TypeOf(target)
		m.Cmd(s, mdb.CREATE, "model", kit.LowerCapital(kit.Select("", kit.Split(t.String(), "."), -1)), ctx.INDEX, m.PrefixKey(), kit.Dict(mdb.TARGET, target))
	})
}
func (s models) Select(m *ice.Message, arg ...string) {
	m.Optionv(mdb.TARGET, mdb.HashSelectTarget(m.Message, kit.Hashs(arg[0]), nil))
}
func (s models) List(m *ice.Message, arg ...string) {
	s.Hash.List(m, arg...)
}

type Models struct{ models }

func init() { ice.Cmd(prefixKey(), models{}) }
func init() { ice.Cmd(prefixKey(), Models{}) }

func (s Models) Bind(m *ice.Message, model string) ice.Any {
	return m.Cmd(s, s.Select, model).Optionv(mdb.TARGET)
}
