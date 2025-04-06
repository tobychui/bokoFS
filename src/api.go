package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"imuslab.com/bokofs/bokofsd/mod/diskinfo"
	"imuslab.com/bokofs/bokofsd/mod/diskinfo/lsblk"
	"imuslab.com/bokofs/bokofsd/mod/diskinfo/smart"
	"imuslab.com/bokofs/bokofsd/mod/netstat"
)

/*
	API Router

	This module handle routing of the API calls
*/

// Primary handler for the API router
func HandlerAPIcalls() http.Handler {
	return http.StripPrefix("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the disk ID from the URL path
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 2 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		diskID := pathParts[1]
		if diskID == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		switch diskID {
		case "info":
			// Request to /api/info/*
			HandleInfoAPIcalls().ServeHTTP(w, r)
			return
		case "smart":
			// Request to /api/smart/*
			HandleSMARTCalls().ServeHTTP(w, r)
			return
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
	}))
}

// Handler for SMART API calls
func HandleSMARTCalls() http.Handler {
	return http.StripPrefix("/smart/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 2 {
			http.Error(w, "Bad Request - Missing disk name", http.StatusBadRequest)
			return
		}
		subPath := pathParts[0]
		diskName := pathParts[1]
		if diskName == "" {
			http.Error(w, "Bad Request - Missing disk name", http.StatusBadRequest)
			return
		}
		switch subPath {
		case "health":
			if diskName == "all" {
				// Get the SMART information for all disks
				allDisks, err := diskinfo.GetAllDisks()
				if err != nil {
					log.Println("Error getting all disks:", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				// Create a map to hold the SMART information for each disk
				diskInfoMap := []*smart.DriveHealthInfo{}
				for _, disk := range allDisks {
					diskName := disk.Name
					health, err := smart.GetDiskSMARTHealthSummary(diskName)
					if err != nil {
						log.Println("Error getting disk health:", err)
						continue
					}

					diskInfoMap = append(diskInfoMap, health)
				}
				// Convert the disk information to JSON and write it to the response
				js, _ := json.Marshal(diskInfoMap)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(js)
				return
			}

			// Get the health status of the disk
			health, err := smart.GetDiskSMARTHealthSummary(diskName)
			if err != nil {
				log.Println("Error getting disk health:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			// Convert the health status to JSON and write it to the response
			js, _ := json.Marshal(health)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(js)
			return
		case "info":
			// Handle SMART API calls
			dt, err := smart.GetDiskType(diskName)
			if err != nil {
				log.Println("Error getting disk type:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if dt == smart.DiskType_SATA {
				// Get SATA disk information
				sataInfo, err := smart.GetSATAInfo(diskName)
				if err != nil {
					log.Println("Error getting SATA disk info:", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				// Convert the SATA info to JSON and write it to the response
				js, _ := json.Marshal(sataInfo)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(js)
			} else if dt == smart.DiskType_NVMe {
				// Get NVMe disk information
				nvmeInfo, err := smart.GetNVMEInfo(diskName)
				if err != nil {
					log.Println("Error getting NVMe disk info:", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				// Convert the NVMe info to JSON and write it to the response
				js, _ := json.Marshal(nvmeInfo)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(js)
			} else {
				log.Println("Unknown disk type:", dt)
				http.Error(w, "Bad Request - Unknown disk type", http.StatusBadRequest)
				return
			}
			return
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
	}))
}

// Handler for info API calls
func HandleInfoAPIcalls() http.Handler {
	return http.StripPrefix("/info/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Check the next part of the URL
		pathParts := strings.Split(r.URL.Path, "/")
		subPath := pathParts[0]
		switch subPath {
		case "netstat":
			// Get the current network statistics
			netstatBuffer.HandleGetBufferedNetworkInterfaceStats(w, r)
			return
		case "iface":
			// Get the list of network interfaces
			netstat.HandleListNetworkInterfaces(w, r)
			return
		case "list":
			// List all block devices and their partitions
			blockDevices, err := lsblk.GetLSBLKOutput()
			if err != nil {
				log.Println("Error getting block devices:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			disks := make([]*diskinfo.Disk, 0)
			for _, device := range blockDevices {
				if device.Type == "disk" {
					disk, err := diskinfo.GetDiskInfo(device.Name)
					if err != nil {
						log.Println("Error getting disk info:", err)
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
						return
					}
					disks = append(disks, disk)
				}
			}
			// Convert the block devices to JSON and write it to the response
			js, _ := json.Marshal(disks)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(js)
		case "disk":
			// Get the disk info for a particular disk, e.g. sda
			if len(pathParts) < 2 {
				http.Error(w, "Bad Request - Invalid disk name", http.StatusBadRequest)
				return
			}
			diskID := pathParts[1]
			if diskID == "" {
				http.Error(w, "Bad Request - Invalid disk name", http.StatusBadRequest)
				return
			}

			if !diskinfo.DevicePathIsValidDisk(diskID) {
				log.Println("Invalid disk ID:", diskID)
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			// Handle diskinfo API calls
			targetDiskInfo, err := diskinfo.GetDiskInfo(diskID)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Convert the disk info to JSON and write it to the response
			js, _ := json.Marshal(targetDiskInfo)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(js)
			return
		case "part":
			// Get the partition info for a particular partition, e.g. sda1
			if len(pathParts) < 2 {
				http.Error(w, "Bad Request - Missing parition name", http.StatusBadRequest)
				return
			}
			partID := pathParts[1]
			if partID == "" {
				http.Error(w, "Bad Request - Missing parition name", http.StatusBadRequest)
				return
			}

			if !diskinfo.DevicePathIsValidPartition(partID) {
				log.Println("Invalid partition name:", partID)
				http.Error(w, "Bad Request - Invalid parition name", http.StatusBadRequest)
				return
			}

			// Handle partinfo API calls
			targetPartInfo, err := diskinfo.GetPartitionInfo(partID)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Convert the partition info to JSON and write it to the response
			js, _ := json.Marshal(targetPartInfo)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(js)

			return
		default:
			fmt.Println("Unknown API call:", subPath)
			http.Error(w, "Not Found", http.StatusNotFound)
		}

	}))
}
