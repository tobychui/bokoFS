package diskinfo

// Disk represents a disk device with its attributes.
type Disk struct {
	Name       string       `json:"name"`                 // Name of the disk, e.g. sda
	Identifier string       `json:"identifier"`           // Disk identifier, e.g. 0x12345678
	Model      string       `json:"model"`                // Disk model, e.g. Samsung SSD 860 EVO 1TB
	Size       int64        `json:"size"`                 // Size of the disk in bytes
	Used       int64        `json:"used"`                 // Used space in bytes, calculated from partitions
	Free       int64        `json:"free"`                 // Free space in bytes, calculated as Size - Used
	DiskLabel  string       `json:"disklabel"`            // Disk label type, e.g. gpt
	BlockType  string       `json:"blocktype"`            // Type of the block device, e.g. disk
	Partitions []*Partition `json:"partitions,omitempty"` // List of partitions on the disk
}

// Partition represents a partition on a disk with its attributes.
type Partition struct {
	UUID       string `json:"uuid"`                 // UUID of the file system
	PartUUID   string `json:"partuuid"`             // Partition UUID
	PartLabel  string `json:"partlabel"`            // Partition label
	Name       string `json:"name"`                 // Name of the partition, e.g. sda1
	Path       string `json:"path"`                 // Path of the partition, e.g. /dev/sda1
	Size       int64  `json:"size"`                 // Size of the partition in bytes
	Used       int64  `json:"used"`                 // Used space in bytes
	Free       int64  `json:"free"`                 // Free space in bytes
	BlockSize  int    `json:"blocksize"`            // Block size in bytes, e.g. 4096
	BlockType  string `json:"blocktype"`            // Type of the block device, e.g. part
	FsType     string `json:"fstype"`               // File system type, e.g. ext4
	MountPoint string `json:"mountpoint,omitempty"` // Mount point of the partition, e.g. /mnt/data
}

// Block represents a block device with its attributes.
type Block struct {
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	Size       int64  `json:"size"`
	BlockSize  int    `json:"blocksize"`
	BlockType  string `json:"blocktype"`
	FsType     string `json:"fstype"`
	MountPoint string `json:"mountpoint,omitempty"`
}
