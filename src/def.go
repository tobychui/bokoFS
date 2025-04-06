package main

import (
	"flag"

	"imuslab.com/bokofs/bokofsd/mod/netstat"
)

var (
	/* Start Flags */
	httpPort    = flag.Int("p", 9000, "Port to serve on (Plain HTTP)")
	devMode     = flag.Bool("dev", false, "Enable development mode")
	config      = flag.String("c", "./config", "Path to the config folder")
	serveSecure = flag.Bool("s", false, "Serve HTTPS. Default false")

	/* Runtime Variables */
	netstatBuffer *netstat.NetStatBuffers
)
