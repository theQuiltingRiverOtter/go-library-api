package main

import (
	"flag"

	"library-api/router"

	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	debug := flag.Bool("debug", true, "sets log level to debug")

	flag.Parse()

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	r := router.SetupRouter("debug")

	r.Run("localhost:8000")

}
