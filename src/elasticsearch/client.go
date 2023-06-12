package es

import (
	"net/http"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/tcp"
	"shylinux.com/x/icebergs/base/web"
	kit "shylinux.com/x/toolkits"
)

type client struct {
	ice.Hash
	short string `data:"sess"`
	field string `data:"time,sess,host,port"`

	insert     string `name:"insert id* data*"`
	addIndex   string `name:"addIndex index*=demo"`
	addMapping string `name:"addMapping mapping*=user properties*:textarea"`
	create     string `name:"create sess*=biz host=localhost port*=10004" help:"创建"`
	list       string `name:"list sess@key index mapping query auto" help:"搜索引擎"`
}

func (s client) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case tcp.PORT:
		m.Cmdy(tcp.SERVER).Cut("port,status,time")
		m.Push(arg[0], "9200")
	case "properties":
		m.Push(arg[0], kit.Format(kit.Dict("name", kit.Dict("type", "string"))))
	case "data":
		m.Push(arg[0], kit.Format(kit.Dict("name", "hi")))
	default:
		s.Hash.Inputs(m, arg...)
	}
}
func (s client) request(m *ice.Message, sess, method, url, data string) ice.Any {
	msg := s.Hash.List(m.Spawn(), kit.Select(m.Option("sess"), sess))
	return kit.UnMarshal(m.Cmdx(method, kit.Format("http://%s:%s/", msg.Append(tcp.HOST), msg.Append(tcp.PORT))+url, web.SPIDE_DATA, data))
}
func (s client) Insert(m *ice.Message, arg ...string) {
	s.request(m, "", http.MethodPost, kit.Format("%s/%s/%s", m.Option("index"), m.Option("mapping"), m.Option("id")), m.Option("data"))
}
func (s client) Delete(m *ice.Message, arg ...string) {
	s.request(m, "", http.MethodDelete, kit.Format("%s/%s/%s", m.Option("index"), m.Option("mapping"), m.Option("id")), "")
}
func (s client) AddMapping(m *ice.Message, arg ...string) {
	data := kit.Format(kit.Dict(m.Option("mapping"), kit.Dict("properties", kit.UnMarshal(m.Option("properties")))))
	s.request(m, "", http.MethodPut, kit.Format("%s/_mapping/%s", m.Option("index"), m.Option("mapping")), data)
}
func (s client) DelMapping(m *ice.Message, arg ...string) {
	s.request(m, "", http.MethodDelete, kit.Format("%s/_mapping/%s", m.Option("index"), m.Option("mapping")), "")
}
func (s client) AddIndex(m *ice.Message, arg ...string) {
	s.request(m, "", http.MethodPut, m.Option("index"), "")
}
func (s client) DelIndex(m *ice.Message, arg ...string) {
	s.request(m, "", http.MethodDelete, m.Option("index"), "")
}
func (s client) List(m *ice.Message, arg ...string) {
	if len(arg) == 0 || arg[0] == "" {
		s.Hash.List(m, arg...).Action(mdb.CREATE)
	} else if len(arg) == 1 || arg[1] == "" {
		kit.For(s.request(m, arg[0], http.MethodGet, "_cat/indices?format=json", ""), func(index int, value ice.Map) {
			m.Push("index", value["index"])
			m.Push("uuid", value["uuid"])
			m.Push("docs.count", value["docs.count"])
			m.Push("store.size", value["store.size"])
		})
		m.PushAction(s.DelIndex).Action(s.AddIndex).StatusTimeCount()
	} else if len(arg) == 2 || arg[2] == "" {
		data := s.request(m, arg[0], http.MethodGet, kit.Format("%s/_mapping?format=json", arg[1]), "")
		m.Debug("what %v", data)
		kit.For(kit.Value(data, kit.Keys(arg[1], "mappings", "properties")), func(key string, value ice.Map) {
			m.Push("mapping", key).Push("properties", kit.Format(value))
		})
		m.PushAction(s.DelMapping).Action(s.AddMapping).StatusTimeCount()
	} else if len(arg) == 3 || arg[3] == "" {
		data := s.request(m, arg[0], http.MethodGet, kit.Format("%s/%s/_search?q=name:*", arg[1], arg[2]), "")
		kit.For(kit.Value(data, "hits.hits"), func(index int, value ice.Map) { m.Push("", value["_source"]) })
		m.PushAction(s.Delete).Action(s.Insert).StatusTimeCount(mdb.TOTAL, kit.Value(data, "hits.total.value"))
	} else {
		data := s.request(m, arg[0], http.MethodGet, kit.Format("%s/%s/_search?q=name:%s", arg[1], arg[2], arg[3]), "")
		kit.For(kit.Value(data, "hits.hits"), func(index int, value ice.Map) { m.Push("", value["_source"]) })
		m.PushAction(s.Delete).Action(s.Insert).StatusTimeCount(mdb.TOTAL, kit.Value(data, "hits.total.value"))
	}
}

func init() { ice.CodeCtxCmd(client{}) }
