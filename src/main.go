package main

import (
	"embed"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"imuslab.com/bokofs/bokofsd/mod/bokofs"
	"imuslab.com/bokofs/bokofsd/mod/bokofs/bokoworker"
)

//go:embed web/*
var embeddedFiles embed.FS

func main() {
	flag.Parse()

	// Start the application
	err := initialization()
	if err != nil {
		panic(err)
	}

	// Capture termination signals and call cleanup
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		fmt.Println("Received termination signal, cleaning up...")
		cleanup()
		os.Exit(0)
	}()

	//DEBUG
	wds, err := bokofs.NewWebdavInterfaceServer("/disk/", "/thumb/")
	if err != nil {
		panic(err)
	}

	test, err := bokoworker.NewFSWorker(&bokoworker.Options{
		NodeName:       "test",
		ServePath:      "./test",
		ThumbnailStore: "./tmp/test/",
	})
	if err != nil {
		panic(err)
	}
	wds.AddWorker(test)

	test2, err := bokoworker.NewFSWorker(&bokoworker.Options{
		NodeName:       "test2",
		ServePath:      "./mod",
		ThumbnailStore: "./tmp/mod/",
	})
	if err != nil {
		panic(err)
	}
	wds.AddWorker(test2)

	//END DEBUG

	/* Static Web Server */
	http.Handle("/", csrfMiddleware(tmplMiddleware(http.FileServer(webfs))))

	/* WebDAV Handlers */
	http.Handle("/disk/", wds.FsHandler())     //Note the trailing slash
	http.Handle("/thumb/", wds.ThumbHandler()) //Note the trailing slash

	/* REST API Handlers */
	http.Handle("/meta", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement handler logic for /meta
		fmt.Fprintln(w, "Meta handler not implemented yet")
	}))

	http.Handle("/api/", csrfMiddleware(HandlerAPIcalls()))

	addr := fmt.Sprintf(":%d", *httpPort)
	fmt.Printf("Starting static web server on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
		os.Exit(1)
	}
}
