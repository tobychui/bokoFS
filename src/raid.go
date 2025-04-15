package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"imuslab.com/bokofs/bokofsd/mod/diskinfo/blkstat"
	"imuslab.com/bokofs/bokofsd/mod/utils"
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
		case "reassemble":
			// Reassemble all RAID devices
			raidManager.HandleForceAssembleReload(w, r)
			return

		case "test":
			devname, err := utils.GetPara(r, "dev")
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			bs, err := blkstat.GetInstalledBus(devname)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			js, _ := json.Marshal(bs)
			utils.SendJSONResponse(w, string(js))
			return
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
	}))
}
