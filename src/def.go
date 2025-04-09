package main

import (
	"flag"
	"net/http"

	"imuslab.com/bokofs/bokofsd/mod/disktool/raid"
	"imuslab.com/bokofs/bokofsd/mod/netstat"
)

const (
	CSRF_COOKIENAME = "bokofs-csrf"
)

var (
	/* Start Flags */
	httpPort = flag.Int("p", 9000, "Port to serve on (Plain HTTP)")
	devMode  = flag.Bool("dev", false, "Enable development mode")
	config   = flag.String("c", "./config", "Path to the config folder")

	//serveSecure = flag.Bool("s", false, "Serve HTTPS. Default false")

	/* Runtime Variables */
	sysuuid        string                          //System UUID (UUIDv4)
	webfs          http.FileSystem                 //The web filesystem for static files
	csrfMiddleware func(http.Handler) http.Handler //CSRF protection middleware

	/* Modules */
	netstatBuffer *netstat.NetStatBuffers
	raidManager   *raid.Manager
)
