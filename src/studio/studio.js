(function() { const DATABASE = "database", TABLE = "table"
Volcanos(chat.ONIMPORT, {
	_init: function(can, msg) { can.onmotion.clear(can), can.ui = can.onappend.layout(can)
		can.onmotion.hidden(can, can.ui.profile), can.onmotion.hidden(can, can.ui.display), can.onmotion.hidden(can, can._status)
		can.db._hash = can.misc.SearchHash(can), can.onimport._session(can, msg, can.ui.project)
		can.sup.onimport._field = function(msg) { can.onimport._plugin(can, msg) }
	},
	_item: function(can, msg, target, cb, key, opts, stack, nick) {
		if (stack && stack.length > 0) { target = can.onimport.itemlist(can, [], function() {}, function() {}, target) }
		var _select; msg.Table(function(value) { value.nick = nick? nick(value): value[key], can.base.Copy(value, opts)
			var _target = can.onimport.item(can, value, function(event) { _target._list || cb(event, value, _target) }, null, target);
			(!_select || can.base.beginWith(can.db._hash, stack.concat(value[key]))) && (_select = _target)
		}), _select && _select.click()
	},
	_session: function(can, msg, target) {
		can.onimport._item(can, msg, target, function(event, value, target) {
			can.run(event, [value.sess], function(msg) { can.onimport._database(can, msg, target, value.sess) })
		}, aaa.SESS, {}, [], function(value) { return `${value.sess}(${value.host}:${value.port})` })
	},
	_database: function(can, msg, target, sess) {
		can.onimport._item(can, msg, target, function(event, value, target) {
			can.run(event, [sess, value.database], function(msg) { can.onimport._table(can, msg, target, sess, value.database) })
		}, DATABASE, {sess: sess}, [sess])
	},
	_table: function(can, msg, target, sess, database) {
		can.onimport._item(can, msg, target, function(event, value, target) { if (value._tabs) { return value._tabs.click() }
			value.title = [sess, database, value.table].join(nfs.PT), value._tabs = can.onimport.tabs(can, [value], function() {
				can.db.current = value, can.misc.SearchHash(can, sess, database, value.table)
				if (can.onmotion.cache(can, function() { return [sess, database, value.table].join(nfs.DF) }, can.ui.content)) { return can.onimport.layout(can) }
				can.onappend.plugin(can, {index: "web.code.mysql.query", args: [sess, database, value.table]}, function(sub) {
					value._sub = sub, sub.onexport.output = function() { can.onimport.layout(can) }
				}, can.ui.content)
			}, function() {})
		}, TABLE, {sess: sess, database: database}, [sess, database], function(value) { return `${value.table} ${value.total}`})
	},
	_plugin: function(can, msg) { msg.Table(function(value) {
		value.nick = value.index, value._tabs = can.onimport.tabs(can, [value], function() { can.db.current = value
			if (can.onmotion.cache(can, function() { return value.index }, can.ui.content)) { return can.onimport.layout(can) }
			can.onappend.plugin(can, value, function(sub) {
				value._sub = sub, sub.onexport.output = function() { can.onimport.layout(can) }
				sub.onexport.title = function(_, title) {
					can.page.Select(can, value._tabs, "span.name", function(target) {
						target.innerHTML = title
					})
				}
			}, can.ui.content)
		})
	}) },
	layout: function(can) { can.ui.layout(can.ConfHeight(), can.ConfWidth(), 0, function(height, width) {
		var sub = can.db.current && can.db.current._sub; sub && sub.onimport.size(sub, height-40, width-40, false)
	}) },
}, [""])
})()
