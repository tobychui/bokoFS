package main

import (
	"net/http"
	"strings"
)

/*
	raid.go

	This file handles the RAID management and monitoring API routing
*/

func HandleRAIDCalls() http.Handler {
	return http.StripPrefix("/raid/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(r.URL.Path, "/")

		switch pathParts[0] {
		case "list":
			// List all RAID devices
			raidManager.HandleListRaidDevices(w, r)
			return
		case "info":
			// Handle loading the detail of a given RAID array, require "dev=md0" as a query parameter
			raidManager.HandleLoadArrayDetail(w, r)
			return
		case "overview":
			// Render the RAID overview page
			raidManager.HandleRenderOverview(w, r)
			return
		case "sync":
			// Get the RAID sync state, require "dev=md0" as a query parameter
			raidManager.HandleGetRAIDSyncState(w, r)
			return
		case "start-resync":
			// Activate a RAID device, require "dev=md0" as a query parameter
			raidManager.HandleSyncPendingToReadWrite(w, r)
			return
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
	}))
}
