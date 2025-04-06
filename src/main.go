package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"imuslab.com/bokofs/bokofsd/mod/bokofs"
	"imuslab.com/bokofs/bokofsd/mod/bokofs/bokoworker"
	"imuslab.com/bokofs/bokofsd/mod/netstat"
)

//go:embed web/*
var embeddedFiles embed.FS

func main() {
	flag.Parse()

	/* File system handler */
	var fileSystem http.FileSystem
	if *devMode {
		fmt.Println("Development mode enabled. Serving files from ./web directory.")
		fileSystem = http.Dir("./web")
	} else {
		fmt.Println("Production mode enabled. Serving files from embedded filesystem.")
		subFS, err := fs.Sub(embeddedFiles, "web")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing embedded subdirectory: %v\n", err)
			os.Exit(1)
		}
		fileSystem = http.FS(subFS)
	}

	configFolderPath := "./config"
	if *config != "" {
		configFolderPath = *config
	}
	if _, err := os.Stat(configFolderPath); os.IsNotExist(err) {
		fmt.Printf("Config folder does not exist. Creating folder at %s\n", configFolderPath)
		if err := os.Mkdir(configFolderPath, os.ModePerm); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config folder: %v\n", err)
			os.Exit(1)
		}
	}

	/* Network statistics */
	nsb, err := netstat.NewNetStatBuffer(300)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating network statistics buffer: %v\n", err)
		os.Exit(1)
	}
	defer netstatBuffer.Close()
	netstatBuffer = nsb

	/* Package Check */
	if !checkRuntimeEnvironment() {
		fmt.Println("Runtime environment check failed. Please install the missing packages.")
		os.Exit(1)
	}

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

	http.Handle("/", http.FileServer(fileSystem))

	/* WebDAV Handlers */
	http.Handle("/disk/", wds.FsHandler())     //Note the trailing slash
	http.Handle("/thumb/", wds.ThumbHandler()) //Note the trailing slash

	/* REST API Handlers */
	http.Handle("/meta", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement handler logic for /meta
		fmt.Fprintln(w, "Meta handler not implemented yet")
	}))

	http.Handle("/api/", HandlerAPIcalls())

	addr := fmt.Sprintf(":%d", *httpPort)
	fmt.Printf("Starting static web server on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
		os.Exit(1)
	}
}
