package es

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/tcp"
	"shylinux.com/x/icebergs/base/web"
	kit "shylinux.com/x/toolkits"
)

type client struct {
	ice.Hash
	ice.Rest
	short string `data:"session"`
	field string `data:"time,session,host,port"`

	create string `name:"create session=biz host=localhost port=10005" help:"创建"`
	list   string `name:"list session@key method=GET,PUT,POST,DELETE path run create text"`
}

func (s client) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case tcp.PORT:
		m.Cmdy(tcp.SERVER).Cut("port,status,time")
	default:
		s.Hash.Inputs(m, arg...)
	}
}
func (s client) List(m *ice.Message, arg ...string) {
	if s.Hash.List(m, arg...); len(arg) > 2 && arg[0] != "" && arg[2] != "" {
		url, args := kit.Format("http://%s:%s%s", m.Append(tcp.HOST), m.Append(tcp.PORT), arg[2]), []string{}
		switch arg[1] {
		case web.SPIDE_PUT, web.SPIDE_POST:
			args = []string{web.SPIDE_DATA, arg[3]}
		}
		m.SetAppend().Cmdy(arg[1], url, args)
	}
}

func init() { ice.CodeCtxCmd(client{}) }
