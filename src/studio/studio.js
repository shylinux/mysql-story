Volcanos(chat.ONIMPORT, {
	_init: function(can, msg) {
		can.ui = can.onappend.layout(can), can.onimport._project(can, msg)
		can.sup.onimport._field = function(sup, msg, cb) { msg.Table(function(value) {
			can.onimport._content(can, [msg.Option(aaa.SESS), msg.Option("database"), msg.Option("table"), value.index], value)
		}); return true }
	},
	_project: function(can, msg) { var _select
		msg.Table(function(value) {
			value.icon = icon.sess, value.nick = `${value.sess}(${value.host}:${value.port})`
			var target = can.onimport.item(can, value, function(event, value, show) {
				show == undefined && can.run(event, [value.sess], function(msg) {
					can.onimport._database(can, msg, value.sess, target)
				})
			}); _select = _select||target, value.sess == can.db.hash[0] && (_select = target)
		}), _select && _select.click()
	},
	_database: function(can, msg, sess, target) {
		can.onimport.itemlist(can, msg.Table(function(value, index) {
			value.icon = icon.database, value.nick = value.database
			value._select = sess == can.db.hash[0] && value.database == can.db.hash[1]
			return value
		}), function(event, value, show) { var target = event.currentTarget
			show == undefined && can.run(event, [sess, value.database], function(msg) {
				can.onimport._table(can, msg, sess, value.database, target)
			})
		}, function() {

		}, target)
	},
	_table: function(can, msg, sess, database, target) {
		can.onimport.itemlist(can, msg.Table(function(value, index) {
			value.icon = icon.table, value.nick = `${value.table}(${value.total})`
			value._select = sess == can.db.hash[0] && database == can.db.hash[1] && value.table == can.db.hash[2]
			return value
		}), function(event, value) { if (value._tabs) { return value._tabs.click() }
			value._tabs = can.onimport._content(can, [sess, database, value.table, "query"], {index: "web.code.mysql.query", args: [sess, database, value.table]}, event.currentTarget)
		}, function() {

		}, target)
	},
	_content: function(can, keys, meta, target) { var key = keys.join(".")
		var _icon = icon.table; can.base.endWith(meta.index, code.XTERM) && (_icon = icon.xterm)
		return can.onimport.tabs(can, [{icon: _icon, nick: keys.slice(1, 3).join("."), title: key}], function() { can.onexport.hash(can, keys)
			can.Option({sess: keys[0], database: keys[1], table: keys[2]})
			target && can.page.Select(can, can.ui.project, html.DIV_ITEM, function(target) { can.page.ClassList.del(can, target, html.SELECT) })
			for (var p = target; p; p = p.parentNode.previousElementSibling) { can.page.ClassList.add(can, p, html.SELECT) }
			if (can.onmotion.cache(can, function(save, load) { save({_content_plugin: can.ui._content_plugin})
				load(key, function(bak) { can.ui._content_plugin = bak._content_plugin }); return key
			}, can.ui.content)) { return }
			can.onappend.plugin(can, meta, function(sub) {
				sub.onexport.output = function(_sub, msg) {
					can.page.Select(sub, sub._option, "div.item.text.id", function(target) {
						can.onmotion.toggle(can, target, msg.append && msg.append.indexOf(mdb.ID) > -1)
					})
				}
				can.ui._content_plugin = sub, can.onimport.layout(can)
			}, can.ui.content)
		}, function() {})
	},
	layout: function(can) { can.ui.layout(can.ConfHeight(), can.ConfWidth(), 0, function(height, width) {
		can.ui._content_plugin && can.ui._content_plugin.onimport.size(can.ui._content_plugin, height, width, false)
	})},
})
Volcanos(chat.ONEXPORT, {
	link: function(can) {
		return can.misc.MergePodCmd(can, {pod: can.ConfSpace()||can.misc.Search(can, ice.POD), cmd: can.ConfIndex()}, true)
	},
})
