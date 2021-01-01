package whgoxy

import "flag"

type Options struct {
	DatabaseFile string `json:"database_file"`
}

func Parse() (opt *Options, err error) {
	opt = &Options{}

	// flags here \/
	flag.StringVar(&opt.DatabaseFile, "f", "./data.sqlite3", "File for database")
	// flags here /\

	flag.Parse()
	return opt, nil
}
