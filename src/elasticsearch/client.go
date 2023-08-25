package es

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/tcp"
	"shylinux.com/x/icebergs/base/web"
	kit "shylinux.com/x/toolkits"
)

type client struct {
	ice.Hash
	short  string `data:"sess"`
	field  string `data:"time,sess,host,port"`
	create string `name:"create sess*=biz host=localhost port*=10004"`
	list   string `name:"list sess@key index mapping id auto" help:"搜索引擎" icon:"elastic.png"`
}

func (s client) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case tcp.PORT:
		m.Cmdy(tcp.PORT, mdb.INPUTS, arg)
		m.Push(arg[0], "9200")
	default:
		s.Hash.Inputs(m, arg...)
	}
}
func (s client) List(m *ice.Message, arg ...string) {
	if len(arg) == 0 || arg[0] == "" {
		s.Hash.List(m, arg...).Action(mdb.CREATE)
	} else if len(arg) == 1 || arg[1] == "" {
		data := s.Get(m, arg[0], "/_cat/indices", "format", "json")
		kit.For(data, func(index int, value ice.Map) {
			m.Push("", value, kit.Split("index,uuid,health,status,docs.count,store.size"))
		})
	} else if len(arg) == 2 || arg[2] == "" {
		data := s.Get(m, arg[0], kit.Format("/%s/_mapping", arg[1]), "format", "json")
		kit.For(kit.Value(data, kit.Keys(arg[1], "mappings")), func(key string, value ice.Map) {
			m.Push("mapping", key).Push("properties", kit.Join(kit.SortedKey(kit.Value(value, "properties"))))
		})
	} else {
		data := s.Get(m, arg[0], kit.Format("/%s/%s/_search", arg[1], arg[2]), "q", "name:*")
		kit.For(kit.Value(data, "hits.hits"), func(index int, value ice.Map) { m.Push("", value["_source"]) })
		m.StatusTimeCountTotal(kit.Value(data, "hits.total.value"))
	}
}

func init() { ice.CodeCtxCmd(client{}) }

func (s client) Get(m *ice.Message, sess, url string, arg ...ice.Any) ice.Any {
	msg := s.Hash.List(m.Spawn(), kit.Select(m.Option(aaa.SESS), sess))
	return web.SpideGet(m.Message, web.HostPort(m.Message, msg.Append(tcp.HOST), msg.Append(tcp.PORT))+url, arg)
}
