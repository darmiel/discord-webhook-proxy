package whgoxy

import "log"

var opt *Options

func New(options *Options) {
	opt = options

	// Init database
	if err := InitDatabase(opt.DatabaseFile); err != nil {
		log.Fatalln("Error received:", err.Error())
		return
	}
}
