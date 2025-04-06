package df

import (
	"bytes"
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

type DiskInfo struct {
	DevicePath string
	Blocks     int64
	Used       int64
	Available  int64
	UsePercent int
	MountedOn  string
}

// GetDiskUsageByPath retrieves disk usage information for a specific path.
// e.g. "/dev/sda1" or "sda1" will return the disk usage for the partition mounted on "/dev/sda1".
func GetDiskUsageByPath(path string) (*DiskInfo, error) {
	//Make sure the path has a prefix and a trailing slash
	if !strings.HasPrefix(path, "/dev/") {
		path = "/dev/" + path
	}

	path = strings.TrimSuffix(path, "/")

	diskUsages, err := GetDiskUsage()
	if err != nil {
		return nil, err
	}

	for _, diskInfo := range diskUsages {
		if strings.HasPrefix(diskInfo.DevicePath, path) {
			return &diskInfo, nil
		}
	}

	return nil, errors.New("disk usage not found for path: " + path)
}

// GetDiskUsage retrieves disk usage information for all mounted filesystems.
func GetDiskUsage() ([]DiskInfo, error) {
	cmd := exec.Command("df", "-k")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(out.String(), "\n")
	if len(lines) < 2 {
		return nil, nil
	}

	var diskInfos []DiskInfo
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		usePercent, err := strconv.Atoi(strings.TrimSuffix(fields[4], "%"))
		if err != nil {
			return nil, err
		}

		blocks, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, err
		}

		used, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			return nil, err
		}

		available, err := strconv.ParseInt(fields[3], 10, 64)
		if err != nil {
			return nil, err
		}

		diskInfos = append(diskInfos, DiskInfo{
			DevicePath: fields[0],
			Blocks:     blocks,
			Used:       used * 1024, // Convert to bytes from 1k blocks
			Available:  available * 1024,
			UsePercent: usePercent,
			MountedOn:  fields[5],
		})
	}

	return diskInfos, nil
}
