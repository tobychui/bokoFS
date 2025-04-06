package raid

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"imuslab.com/bokofs/bokofsd/mod/disktool/diskfs"
)

// Get the next avaible RAID array name
func GetNextAvailableMDDevice() (string, error) {
	for i := 0; i < 100; i++ {
		mdDevice := fmt.Sprintf("/dev/md%d", i)
		if _, err := os.Stat(mdDevice); os.IsNotExist(err) {
			return mdDevice, nil
		}
	}

	return "", fmt.Errorf("no available /dev/mdX devices found")
}

// Check if the given device is safe to remove from the array without losing data
func (m *Manager) IsSafeToRemove(mdDev string, sdXDev string) bool {
	targetRAIDVol, err := m.GetRAIDDeviceByDevicePath(mdDev)
	if err != nil {
		return false
	}

	//Trim off the /dev/ part if exists
	sdXDev = filepath.Base(sdXDev)

	//Check how many members left if this is removed
	remainingMemebers := 0
	for _, member := range targetRAIDVol.Members {
		if member.Name != sdXDev {
			remainingMemebers++
		}
	}

	//Check if removal of sdX will cause data loss
	if strings.EqualFold(targetRAIDVol.Level, "raid0") {
		return false
	} else if strings.EqualFold(targetRAIDVol.Level, "raid1") {
		//In raid1, you need at least 1 disk to hold data
		return remainingMemebers >= 1
	} else if strings.EqualFold(targetRAIDVol.Level, "raid5") {
		//In raid 5, at least 2 disk is needed before data loss
		return remainingMemebers >= 2
	} else if strings.EqualFold(targetRAIDVol.Level, "raid6") {
		//In raid 6, you need 6 disks with max loss = 2 disks
		return remainingMemebers >= 2
	}

	return true
}

// Check if the given disk (sdX) is currently used in any volume
func (m *Manager) DiskIsUsedInAnotherRAIDVol(sdXDev string) (bool, error) {
	raidPools, err := m.GetRAIDDevicesFromProcMDStat()
	if err != nil {
		return false, errors.New("unable to access RAID controller state")
	}

	for _, md := range raidPools {
		for _, member := range md.Members {
			if member.Name == filepath.Base(sdXDev) {
				return true, nil
			}
		}
	}

	return false, nil
}

// Check if the given disk (sdX) is root drive (the disk that install the OS, aka /)
func (m *Manager) DiskIsRoot(sdXDev string) (bool, error) {
	bdMeta, err := diskfs.GetBlockDeviceMeta(sdXDev)
	if err != nil {
		return false, err
	}

	for _, partition := range bdMeta.Children {
		if partition.Mountpoint == "/" {
			//Root
			return true, nil
		}
	}
	return false, nil
}

// ClearSuperblock clears the superblock of the specified disk so it can be used safely
func (m *Manager) ClearSuperblock(devicePath string) error {
	isMounted, err := diskfs.DeviceIsMounted(devicePath)
	if err != nil {
		return errors.New("unable to validate if the device is unmounted: " + err.Error())
	}
	if isMounted {
		return errors.New("target device is mounted. Make sure it is unmounted before clearing")
	}

	//Make sure there are /dev/ in front of the device path
	if !strings.HasPrefix(devicePath, "/dev/") {
		devicePath = filepath.Join("/dev/", devicePath)
	}
	cmd := exec.Command("sudo", "mdadm", "--zero-superblock", devicePath)

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error clearing superblock: %v", err)
	}

	return nil
}

// Use to restart any not-removed RAID device
func (m *Manager) RestartRAIDService() error {
	cmd := exec.Command("sudo", "mdadm", "--assemble", "--scan")

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		if string(output) == "" {
			//Nothing updated in config.
			return nil
		}
		return fmt.Errorf("error restarting RAID device: %s", strings.TrimSpace(string(output)))
	}

	return nil
}

// Stop RAID device with given path
func (m *Manager) StopRAIDDevice(devicePath string) error {
	cmd := exec.Command("sudo", "mdadm", "--stop", devicePath)

	// Run the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error stopping RAID device: %v", err)
	}

	return nil
}

// RemoveRAIDDevice removes the specified RAID device member (disk).
func (m *Manager) RemoveRAIDMember(devicePath string) error {
	// Construct the mdadm command to remove the RAID device
	cmd := exec.Command("sudo", "mdadm", "--remove", devicePath)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If there was an error, return the combined output and the error message
		return fmt.Errorf("error removing RAID device: %s", strings.TrimSpace(string(output)))
	}

	return nil
}

// IsValidRAIDLevel checks if the given RAID level is valid.
func IsValidRAIDLevel(level string) bool {
	// List of valid RAID levels
	validLevels := []string{"raid1", "raid0", "raid6", "raid5", "raid4", "raid10"}

	// Convert the RAID level to lowercase and remove any surrounding whitespace
	level = strings.TrimSpace(strings.ToLower(level))

	// Check if the level exists in the list of valid levels
	for _, validLevel := range validLevels {
		if level == validLevel {
			return true
		}
	}

	// Return false if the level is not found in the list of valid levels
	return false
}

// Get RAID device info from device path
func (m *Manager) GetRAIDDeviceByDevicePath(devicePath string) (*RAIDDevice, error) {
	//Strip the /dev/ part if it was accidentally passed in
	devicePath = filepath.Base(devicePath)

	//Get all the raid devices
	rdevs, err := m.GetRAIDDevicesFromProcMDStat()
	if err != nil {
		return nil, err
	}

	//Check for match
	for _, rdev := range rdevs {
		if rdev.Name == devicePath {
			return &rdev, nil
		}
	}

	return nil, errors.New("target RAID device not found")
}

// Check if a RAID device exists, e.g. md0
func (m *Manager) RAIDDeviceExists(devicePath string) bool {
	_, err := m.GetRAIDDeviceByDevicePath(devicePath)
	return err == nil
}

// Check if a RAID contain disk that failed or degraded given the devicePath, e.g. md0 or /dev/md0
func (m *Manager) RAIDArrayContainsFailedDisks(devicePath string) (bool, error) {
	raidDeviceInfo, err := m.GetRAIDInfo(devicePath)
	if err != nil {
		return false, err
	}
	return strings.Contains(raidDeviceInfo.State, "degraded") || strings.Contains(raidDeviceInfo.State, "faulty"), nil
}

// GetRAIDPartitionSize returns the size of the RAID partition in bytes as an int64
func GetRAIDPartitionSize(devicePath string) (int64, error) {
	// Ensure devicePath is formatted correctly
	if !strings.HasPrefix(devicePath, "/dev/") {
		devicePath = "/dev/" + devicePath
	}

	// Execute the df command with the device path
	cmd := exec.Command("df", "--block-size=1", devicePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("failed to execute df command: %v", err)
	}

	// Parse the output to find the size
	lines := strings.Split(out.String(), "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("unexpected df output: %s", out.String())
	}

	// The second line should contain the relevant information
	fields := strings.Fields(lines[1])
	if len(fields) < 2 {
		return 0, fmt.Errorf("unexpected df output: %s", lines[1])
	}

	// The second field should be the size in bytes
	size, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse size: %v", err)
	}

	return size, nil
}

// GetRAIDUsedSize returns the used size of the RAID array in bytes as an int64
func GetRAIDUsedSize(devicePath string) (int64, error) {
	// Ensure devicePath is formatted correctly
	if !strings.HasPrefix(devicePath, "/dev/") {
		devicePath = "/dev/" + devicePath
	}

	// Execute the df command with the device path
	cmd := exec.Command("df", "--block-size=1", devicePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("failed to execute df command: %v", err)
	}

	// Parse the output to find the used size
	lines := strings.Split(out.String(), "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("unexpected df output: %s", out.String())
	}

	// The second line should contain the relevant information
	fields := strings.Fields(lines[1])
	if len(fields) < 3 {
		return 0, fmt.Errorf("unexpected df output: %s", lines[1])
	}

	// The third field should be the used size in bytes
	usedSize, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse used size: %v", err)
	}

	return usedSize, nil
}
