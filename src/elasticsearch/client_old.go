package es

import (
	"net/http"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/tcp"
	"shylinux.com/x/icebergs/base/web"
	kit "shylinux.com/x/toolkits"
)

type oldclient struct {
	ice.Hash
	short string `data:"sess"`
	field string `data:"time,sess,host,port"`

	create string `name:"create sess=biz host=localhost port=10004" help:"创建"`
	list   string `name:"list sess@key method=GET,PUT,POST,DELETE path run create text" help:"搜索引擎"`
}

func (s oldclient) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case tcp.PORT:
		m.Cmdy(tcp.SERVER).Cut("port,status,time")
	default:
		s.Hash.Inputs(m, arg...)
	}
}
func (s oldclient) List(m *ice.Message, arg ...string) {
	if s.Hash.List(m, arg...); len(arg) > 2 && arg[0] != "" && arg[2] != "" {
		url, args := kit.Format("http://%s:%s%s", m.Append(tcp.HOST), m.Append(tcp.PORT), arg[2]), []string{}
		switch arg[1] {
		case http.MethodPut, http.MethodPost:
			args = []string{web.SPIDE_DATA, arg[3]}
		}
		m.SetAppend().Cmdy(arg[1], url, args)
	}
}

func init() { ice.Cmd("web.code.elasticsearch.old.client", oldclient{}) }
