//go:build linux
// +build linux

package blkstat

/*
	blkstat.go

	This file extract the realtime disk I/O statistics from the Linux kernel
	by reading the /sys/block/<block_name>/stat file.

	Mostly you will only need to use ReadIOs, ReadSectors, WriteIOs and WriteSectors
	to get the I/O statistics. Note that the values are accumulated since the
	system booted, so you will need to calculate the difference between two
	consecutive calls to get the I/O rate.
*/
import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type BlockStat struct {
	ReadIOs      uint64
	ReadMerges   uint64
	ReadSectors  uint64
	ReadTicks    uint64
	WriteIOs     uint64
	WriteMerges  uint64
	WriteSectors uint64
	WriteTicks   uint64
	InFlight     uint64
	IoTicks      uint64
	TimeInQueue  uint64
}

type InstallPosition struct {
	PCIEBusAddress string // PCIe bus address of the device
	SATAPort       string // SATA port location of the device
	USBPort        string // USB port in {hub_id}-{port_num} format
	NVMESlot       string // NVMe slot information
}

// GetInstalledBus retrieves the PCIe bus address, SATA port location, USB port, and NVMe slot for a given block device.
func GetInstalledBus(blockName string) (*InstallPosition, error) {
	linkPath := fmt.Sprintf("/sys/block/%s", blockName)
	realPath, err := os.Readlink(linkPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read symlink for block device: %w", err)
	}

	// Extract PCIe bus address, SATA port location, USB port, and NVMe slot from the resolved path
	parts := strings.Split(realPath, "/")
	var pcieBusAddress, sataPort, usbPort, nvmeSlot string
	for i, part := range parts {
		if strings.HasPrefix(part, "pci") {
			pcieBusAddress = part
		} else if strings.HasPrefix(part, "ata") {
			sataPort = part
		} else if strings.HasPrefix(part, "usb") {
			if i+1 < len(parts) && strings.Contains(parts[i+1], ":") {
				usbPort = parts[i] // USB port in {hub_id}-{port_num} format
			}
		} else if strings.HasPrefix(part, "nvme") {
			if i+1 < len(parts) && strings.HasPrefix(parts[i+1], "nvme") {
				nvmeSlot = parts[i+1] // NVMe slot information
			}
		}
	}

	if pcieBusAddress == "" && sataPort == "" && usbPort == "" && nvmeSlot == "" {
		return nil, fmt.Errorf("failed to extract PCIe bus address, SATA port, USB port, or NVMe slot")
	}

	return &InstallPosition{
		PCIEBusAddress: pcieBusAddress,
		SATAPort:       sataPort,
		USBPort:        usbPort,
		NVMESlot:       nvmeSlot,
	}, nil
}

// GetBlockStat retrieves the block statistics for a given block device.
func GetBlockStat(blockName string) (*BlockStat, error) {
	statPath := fmt.Sprintf("/sys/block/%s/stat", blockName)
	data, err := os.ReadFile(statPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read stat file: %w", err)
	}

	fields := strings.Fields(string(data))
	if len(fields) < 11 {
		return nil, fmt.Errorf("unexpected stat file format")
	}

	values := make([]uint64, 11)
	for i := 0; i < 11; i++ {
		values[i], err = strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse stat value: %w", err)
		}
	}

	return &BlockStat{
		ReadIOs:      values[0],
		ReadMerges:   values[1],
		ReadSectors:  values[2],
		ReadTicks:    values[3],
		WriteIOs:     values[4],
		WriteMerges:  values[5],
		WriteSectors: values[6],
		WriteTicks:   values[7],
		InFlight:     values[8],
		IoTicks:      values[9],
		TimeInQueue:  values[10],
	}, nil
}
