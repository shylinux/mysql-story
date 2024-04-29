(function() {
const DATABASE = "database", TABLE = "table"
Volcanos(chat.ONIMPORT, {
	_init: function(can, msg) {
		can.ui = can.onappend.layout(can), can.onimport._project(can, msg)
	},
	_project: function(can, msg) { var _select
		msg.Table(function(value) {
			var target = can.onimport.item(can, value, function(event, value, show) {
				show == undefined && can.run(event, [value.sess], function(msg) {
					can.onimport._database(can, msg, value.sess, target)
				})
			}); _select = _select||target, value.sess == can.db.hash[0] && (_select = target)
		}), _select.click()
	},
	_database: function(can, msg, sess, target) {
		can.onimport.itemlist(can, msg.Table(function(value, index) { value.nick = value.database
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
		can.onimport.itemlist(can, msg.Table(function(value, index) { value.nick = `${value.table}(${value.total})`
			value._select = sess == can.db.hash[0] && database == can.db.hash[1] && value.table == can.db.hash[2]
			return value
		}), function(event, value) { if (value._tabs) { return value._tabs.click() }
			value._tabs = can.onimport._content(can, [sess, database, value.table, "query"], {index: "web.code.mysql.query", args: [sess, database, value.table]}, event.currentTarget)
		}, function() {

		}, target)
	},
	_content: function(can, keys, meta, target) {
		return can.onimport.tabs(can, [{nick: keys.join(".")}], function() { can.onexport.hash(can, keys)
			can.page.Select(can, can.ui.project, html.DIV_ITEM, function(target) {
				can.page.ClassList.del(can, target, html.SELECT)
			})
			for (var p = target; p; p = p.parentNode.previousElementSibling) {
				can.page.ClassList.add(can, p, html.SELECT)
			}
			if (can.onmotion.cache(can, function(save, load) {
				save({
					_content_plugin: can.ui._content_plugin,
				}), load(keys.join("."), function(bak) {
					can.ui._content_plugin = bak._content_plugin
				}); return keys.join(".")
			}, can.ui.content)) { return }
			can.onappend.plugin(can, meta, function(sub) { can.ui._content_plugin = sub
				can.onimport.layout(can)
			}, can.ui.content)
		}, function() {})
	},
	layout: function(can) { can.ui.layout(can.ConfHeight(), can.ConfWidth(), 0, function(height, width) {
		can.ui._content_plugin && can.ui._content_plugin.onimport.size(can.ui._content_plugin, height-40, width-40, false)
	})},
}, [""])
})()
