package whgoxy

import "log"

var Opt *Options

func New(options *Options) {
	Opt = options

	if options.MongoUse {
		InitMongoDatabase(options)
	} else {
		// No database selected
		log.Fatalln("Error: No database selected.")
		return
	}
}
