package smartgo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anatol/smart.go"
)

// SMART data structure for SATA disks
type SATAAttrData struct {
	Id      uint8
	Name    string
	Type    int
	RawVal  uint64
	Current uint8
	Worst   uint8
}

type SMARTData struct {
	ModelNumber  string
	SerialNumber string
	Size         uint64
	Temperature  int

	/* NVME specific fields */
	NameSpaceUtilizations []uint64
	PowerOnHours          uint64
	PowerCycles           uint64
	UnsafeShutdowns       uint64
	MediaErrors           uint64

	/* SATA specific fields */
	SATAAttrs []*SATAAttrData
}

// Get SMART data of a particular disk / device
func GetSMARTData(disk string) (*SMARTData, error) {
	if !strings.HasPrefix(disk, "/dev/") {
		disk = filepath.Join("/dev", disk)
	}

	//Check if the disk exists
	if _, err := os.Stat(disk); os.IsNotExist(err) {
		return nil, fmt.Errorf("disk %s does not exist", disk)
	}

	// Check if the disk is NVMe or SATA
	isNVMe := strings.HasPrefix(disk, "/dev/nvme")
	isSATA := strings.HasPrefix(disk, "/dev/sd")

	// If the disk is not NVMe or SATA, return an error
	if !isNVMe && !isSATA {
		return nil, fmt.Errorf("disk %s is not an NVMe or SATA disk", disk)
	}

	if isNVMe {
		return getNVMESMART(disk)
	} else if isSATA {
		return getSATASMART(disk)
	}

	return nil, errors.New("unsupported disk type")
}

// Check if the path is the disk path instead of partition path
func IsRootDisk(deviceFilePath string) bool {
	deviceFilename := filepath.Base(deviceFilePath)
	if !(strings.HasPrefix(deviceFilename, "sd") || strings.HasPrefix(deviceFilename, "nvme")) {
		return false
	}
	if strings.HasPrefix(deviceFilename, "sd") && len(deviceFilename) > 3 {
		return false
	}
	if strings.HasPrefix(deviceFilename, "nvme") && len(deviceFilename) > 5 {
		return false
	}
	return true
}

// Check if this SMART implementation supports the disk type
func IsDiskSupportedType(disk string) bool {
	return IsNVMeDevice(disk) || IsSATADevice(disk)
}

// Check if the disk is a SATA device
func IsSATADevice(disk string) bool {
	if !strings.HasPrefix(disk, "/dev/") {
		disk = filepath.Join("/dev", disk)
	}

	_, err := smart.OpenSata(disk)
	return err == nil
}

// Check if the disk is a NVMe device
func IsNVMeDevice(disk string) bool {
	if !strings.HasPrefix(disk, "/dev/") {
		disk = filepath.Join("/dev", disk)
	}

	_, err := smart.OpenNVMe(disk)
	return err == nil
}
