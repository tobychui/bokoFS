package lsblk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// BlockDevice represents a block device and its attributes.
type BlockDevice struct {
	Name       string        `json:"name"`                 //e.g. sda (disk) or sda1 (partition)
	Size       int64         `json:"size"`                 // Size in bytes (manufacturer size)
	Type       string        `json:"type"`                 // Type of the block device (e.g., disk, part)
	MountPoint string        `json:"mountpoint,omitempty"` // Mount point of the device
	Children   []BlockDevice `json:"children,omitempty"`   // List of child devices (e.g., partitions)
}

// parseLSBLKJSONOutput parses the JSON output of the `lsblk` command into a slice of BlockDevice structs.
func parseLSBLKJSONOutput(output string) ([]BlockDevice, error) {
	var result struct {
		BlockDevices []BlockDevice `json:"blockdevices"`
	}

	err := json.Unmarshal([]byte(output), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse lsblk JSON output: %w", err)
	}

	return result.BlockDevices, nil
}

// GetLSBLKOutput runs the `lsblk` command with JSON output and returns its output as a slice of BlockDevice structs.
func GetLSBLKOutput() ([]BlockDevice, error) {
	cmd := exec.Command("lsblk", "-o", "NAME,SIZE,TYPE,MOUNTPOINT", "-b", "-J")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return parseLSBLKJSONOutput(out.String())
}

// GetBlockDeviceInfoFromDevicePath retrieves block device information for a given device path.
func GetBlockDeviceInfoFromDevicePath(devname string) (*BlockDevice, error) {
	devname = strings.TrimPrefix(devname, "/dev/")
	if strings.Contains(devname, "/") {
		return nil, fmt.Errorf("invalid device name: %s", devname)
	}

	// Get the block device info using lsblk
	// and filter for the specified device name.
	devices, err := GetLSBLKOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get block device info: %w", err)
	}

	for _, device := range devices {
		if device.Name == devname {
			return &device, nil
		} else if device.Children != nil {
			for _, child := range device.Children {
				if child.Name == devname {
					return &child, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("device %s not found", devname)
}
