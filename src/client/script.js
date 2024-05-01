Volcanos(chat.ONIMPORT, {
	_init: function(can, msg, cb) {
		can.require(["/plugin/local/code/vimer.js"], function(can) {
			can.onimport._last_init(can, msg, function() {
				can.onmotion.hidden(can, can.ui.project)
				can.onaction.list = ["save", "exec"], cb && cb(msg)
			})
		})
	},
})