//
// File: app.go
//
//

package main


func app_begin(obj *objtype) {
	msg := "Begin"
	display_debug_message(obj, msg)
}

func app_end(obj *objtype) {
	msg := "End"
	display_debug_message(obj, msg)
}

func app_run(obj *objtype) {
	data := loadconfig(obj, *configfile)
	process(obj, *data)
}
