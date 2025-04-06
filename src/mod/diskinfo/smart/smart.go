package smart

/*
	SMART.go

	This script uses the smartctl command to retrieve information about the disk.
	It supports both NVMe and SATA disks on Linux systems only.
*/

import (
	"bufio"
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"imuslab.com/bokofs/bokofsd/mod/diskinfo"
)

// GetDiskType checks if the disk is NVMe or SATA
func GetDiskType(disk string) (DiskType, error) {
	if !strings.HasPrefix(disk, "/dev/") {
		disk = "/dev/" + disk
	}

	//Make sure the target is a disk
	if !diskinfo.DevicePathIsValidDisk(disk) {
		return DiskType_Unknown, errors.New("disk is not a valid disk")
	}

	//Check if the disk is a NVMe or SATA disk
	if strings.HasPrefix(disk, "/dev/nvme") {
		return DiskType_NVMe, nil
	} else if strings.HasPrefix(disk, "/dev/sd") {
		return DiskType_SATA, nil
	}
	return DiskType_Unknown, errors.New("disk is not NVMe or SATA")
}

// GetNVMEInfo retrieves NVMe disk information using smartctl
func GetNVMEInfo(disk string) (*NVMEInfo, error) {
	if !strings.HasPrefix(disk, "/dev/") {
		disk = "/dev/" + disk
	}

	cmd := exec.Command("smartctl", "-i", disk)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	info := &NVMEInfo{}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Model Number:") {
			info.ModelNumber = strings.TrimSpace(strings.TrimPrefix(line, "Model Number:"))
		} else if strings.HasPrefix(line, "Serial Number:") {
			info.SerialNumber = strings.TrimSpace(strings.TrimPrefix(line, "Serial Number:"))
		} else if strings.HasPrefix(line, "Firmware Version:") {
			info.FirmwareVersion = strings.TrimSpace(strings.TrimPrefix(line, "Firmware Version:"))
		} else if strings.HasPrefix(line, "PCI Vendor/Subsystem ID:") {
			info.PCIVendorSubsystemID = strings.TrimSpace(strings.TrimPrefix(line, "PCI Vendor/Subsystem ID:"))
		} else if strings.HasPrefix(line, "IEEE OUI Identifier:") {
			info.IEEEOUIIdentifier = strings.TrimSpace(strings.TrimPrefix(line, "IEEE OUI Identifier:"))
		} else if strings.HasPrefix(line, "Total NVM Capacity:") {
			info.TotalNVMeCapacity = strings.TrimSpace(strings.TrimPrefix(line, "Total NVM Capacity:"))
		} else if strings.HasPrefix(line, "Unallocated NVM Capacity:") {
			info.UnallocatedNVMeCapacity = strings.TrimSpace(strings.TrimPrefix(line, "Unallocated NVM Capacity:"))
		} else if strings.HasPrefix(line, "Controller ID:") {
			info.ControllerID = strings.TrimSpace(strings.TrimPrefix(line, "Controller ID:"))
		} else if strings.HasPrefix(line, "NVMe Version:") {
			info.NVMeVersion = strings.TrimSpace(strings.TrimPrefix(line, "NVMe Version:"))
		} else if strings.HasPrefix(line, "Number of Namespaces:") {
			info.NumberOfNamespaces = strings.TrimSpace(strings.TrimPrefix(line, "Number of Namespaces:"))
		} else if strings.HasPrefix(line, "Namespace 1 Size/Capacity:") {
			info.NamespaceSizeCapacity = strings.TrimSpace(strings.TrimPrefix(line, "Namespace 1 Size/Capacity:"))
		} else if strings.HasPrefix(line, "Namespace 1 Utilization:") {
			info.NamespaceUtilization = strings.TrimSpace(strings.TrimPrefix(line, "Namespace 1 Utilization:"))
		} else if strings.HasPrefix(line, "Namespace 1 Formatted LBA Size:") {
			info.NamespaceFormattedLBASize = strings.TrimSpace(strings.TrimPrefix(line, "Namespace 1 Formatted LBA Size:"))
		} else if strings.HasPrefix(line, "Namespace 1 IEEE EUI-64:") {
			info.NamespaceIEEE_EUI_64 = strings.TrimSpace(strings.TrimPrefix(line, "Namespace 1 IEEE EUI-64:"))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return info, nil
}

// GetSATADiskInfo retrieves SATA disk information using smartctl
func GetSATAInfo(disk string) (*SATADiskInfo, error) {
	if !strings.HasPrefix(disk, "/dev/") {
		disk = "/dev/" + disk
	}

	cmd := exec.Command("smartctl", "-i", disk)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	info := &SATADiskInfo{}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Model Family:") {
			info.ModelFamily = strings.TrimSpace(strings.TrimPrefix(line, "Model Family:"))
		} else if strings.HasPrefix(line, "Device Model:") {
			info.DeviceModel = strings.TrimSpace(strings.TrimPrefix(line, "Device Model:"))
		} else if strings.HasPrefix(line, "Serial Number:") {
			info.SerialNumber = strings.TrimSpace(strings.TrimPrefix(line, "Serial Number:"))
		} else if strings.HasPrefix(line, "Firmware Version:") {
			info.Firmware = strings.TrimSpace(strings.TrimPrefix(line, "Firmware Version:"))
		} else if strings.HasPrefix(line, "User Capacity:") {
			info.UserCapacity = strings.TrimSpace(strings.TrimPrefix(line, "User Capacity:"))
		} else if strings.HasPrefix(line, "Sector Size:") {
			info.SectorSize = strings.TrimSpace(strings.TrimPrefix(line, "Sector Size:"))
		} else if strings.HasPrefix(line, "Rotation Rate:") {
			info.RotationRate = strings.TrimSpace(strings.TrimPrefix(line, "Rotation Rate:"))
		} else if strings.HasPrefix(line, "Form Factor:") {
			info.FormFactor = strings.TrimSpace(strings.TrimPrefix(line, "Form Factor:"))
		} else if strings.HasPrefix(line, "SMART support is:") {
			info.SmartSupport = strings.TrimSpace(strings.TrimPrefix(line, "SMART support is:")) == "Enabled"
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return info, nil
}

// SetSMARTEnableOnDisk enables or disables SMART on the specified disk
func SetSMARTEnableOnDisk(disk string, isEnabled bool) error {
	if !strings.HasPrefix(disk, "/dev/") {
		disk = "/dev/" + disk
	}

	enableCmd := "off"
	if isEnabled {
		enableCmd = "on"
	}

	cmd := exec.Command("smartctl", "-s", enableCmd, disk)
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	if strings.Contains(string(output), "SMART Enabled") {
		return nil
	} else {
		// Print the command output to STDOUT if enabling SMART failed
		println(string(output))
		return errors.New("failed to enable SMART on disk")
	}
}

// GetDiskSMARTCheck retrieves the SMART health status of the specified disk
// Usually only returns "PASSED" or "FAILED"
func GetDiskSMARTCheck(diskname string) (*SMARTTestResult, error) {
	if !strings.HasPrefix(diskname, "/dev/") {
		diskname = "/dev/" + diskname
	}

	cmd := exec.Command("smartctl", "-H", "-A", diskname)
	output, err := cmd.Output()
	if err != nil {
		// Check if the error is due to exit code 32 (non-critical error for some disks)
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 32 {
			// Ignore the error and proceed
		} else {
			// Print the command output to STDOUT if the command fails
			println(string(output))
			return nil, err
		}
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	result := &SMARTTestResult{
		TestResult:         "Unknown",
		MarginalAttributes: make([]SMARTAttribute, 0),
	}
	var inAttributesSection bool = false

	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		// Check for overall health result
		if strings.HasPrefix(line, "SMART overall-health self-assessment test result:") {
			result.TestResult = strings.TrimSpace(strings.TrimPrefix(line, "SMART overall-health self-assessment test result:"))
		}

		// Detect the start of the attributes section
		if strings.HasPrefix(line, "ID# ATTRIBUTE_NAME") {
			inAttributesSection = true
			continue
		}

		// Parse marginal attributes
		if inAttributesSection {
			fields := strings.Fields(line)
			if len(fields) >= 10 {
				id, err := strconv.Atoi(fields[0])
				if err != nil {
					continue
				}
				value, err := strconv.Atoi(fields[3])
				if err != nil {
					continue
				}
				worst, err := strconv.Atoi(fields[4])
				if err != nil {
					continue
				}
				threshold, err := strconv.Atoi(fields[5])
				if err != nil {
					continue
				}

				attribute := SMARTAttribute{
					ID:         id,
					Name:       fields[1],
					Flag:       fields[2],
					Value:      value,
					Worst:      worst,
					Threshold:  threshold,
					Type:       fields[6],
					Updated:    fields[7],
					WhenFailed: fields[8],
					RawValue:   strings.Join(fields[9:], " "),
				}
				result.MarginalAttributes = append(result.MarginalAttributes, attribute)

			}
		}
	}

	if err := scanner.Err(); err != nil {
		// Print the command output to STDOUT if parsing failed
		println(string(output))
		return nil, err
	}

	if result.TestResult == "" {
		return nil, errors.New("unable to determine SMART health status")
	}

	return result, nil
}

func GetDiskSMARTHealthSummary(diskname string) (*DriveHealthInfo, error) {
	smartCheck, err := GetDiskSMARTCheck(diskname)
	if err != nil {
		return nil, err
	}

	healthInfo := &DriveHealthInfo{
		DeviceName: diskname,
		IsHealthy:  strings.ToUpper(smartCheck.TestResult) == "PASSED",
	}

	//Populate the device model and serial number from SMARTInfo
	dt, err := GetDiskType(diskname)
	if err != nil {
		return nil, err
	}

	if dt == DiskType_SATA {
		sataInfo, err := GetSATAInfo(diskname)
		if err != nil {
			return nil, err
		}
		healthInfo.DeviceModel = sataInfo.DeviceModel
		healthInfo.SerialNumber = sataInfo.SerialNumber
		healthInfo.IsSSD = strings.Contains(sataInfo.RotationRate, "Solid State")
	} else if dt == DiskType_NVMe {
		nvmeInfo, err := GetNVMEInfo(diskname)
		if err != nil {
			return nil, err
		}
		healthInfo.DeviceModel = nvmeInfo.ModelNumber
		healthInfo.SerialNumber = nvmeInfo.SerialNumber
		healthInfo.IsNVMe = true
	} else {
		return nil, errors.New("unsupported disk type")
	}

	for _, attr := range smartCheck.MarginalAttributes {
		switch attr.Name {
		case "Power_On_Hours":
			if value, err := strconv.ParseUint(attr.RawValue, 10, 64); err == nil {
				healthInfo.PowerOnHours = value
			}
		case "Power_Cycle_Count":
			if value, err := strconv.ParseUint(attr.RawValue, 10, 64); err == nil {
				healthInfo.PowerCycleCount = value
			}
		case "Reallocated_Sector_Ct":
			if value, err := strconv.ParseUint(attr.RawValue, 10, 64); err == nil {
				healthInfo.ReallocatedSectors = value
			}
		case "Wear_Leveling_Count":
			if value, err := strconv.ParseUint(attr.RawValue, 10, 64); err == nil {
				healthInfo.WearLevelingCount = value
			}
		case "Uncorrectable_Error_Cnt":
			if value, err := strconv.ParseUint(attr.RawValue, 10, 64); err == nil {
				healthInfo.UncorrectableErrors = value
			}
		case "Current_Pending_Sector":
			if value, err := strconv.ParseUint(attr.RawValue, 10, 64); err == nil {
				healthInfo.PendingSectors = value
			}
		case "ECC_Recovered":
			if value, err := strconv.ParseUint(attr.RawValue, 10, 64); err == nil {
				healthInfo.ECCRecovered = value
			}
		case "UDMA_CRC_Error_Count":
			if value, err := strconv.ParseUint(attr.RawValue, 10, 64); err == nil {
				healthInfo.UDMACRCErrors = value
			}
		case "Total_LBAs_Written":
			if value, err := strconv.ParseUint(attr.RawValue, 10, 64); err == nil {
				healthInfo.TotalLBAWritten = value
			}
		case "Total_LBAs_Read":
			if value, err := strconv.ParseUint(attr.RawValue, 10, 64); err == nil {
				healthInfo.TotalLBARead = value
			}
		}
	}

	return healthInfo, nil
}
