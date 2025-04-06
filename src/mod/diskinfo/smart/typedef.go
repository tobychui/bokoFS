package smart

type SATADiskInfo struct {
	ModelFamily  string
	DeviceModel  string
	SerialNumber string
	Firmware     string
	UserCapacity string
	SectorSize   string
	RotationRate string
	FormFactor   string
	SmartSupport bool
}

type NVMEInfo struct {
	ModelNumber               string
	SerialNumber              string
	FirmwareVersion           string
	PCIVendorSubsystemID      string
	IEEEOUIIdentifier         string
	TotalNVMeCapacity         string
	UnallocatedNVMeCapacity   string
	ControllerID              string
	NVMeVersion               string
	NumberOfNamespaces        string
	NamespaceSizeCapacity     string
	NamespaceUtilization      string
	NamespaceFormattedLBASize string
	NamespaceIEEE_EUI_64      string
}

type DiskType int

const (
	DiskType_Unknown DiskType = iota
	DiskType_NVMe
	DiskType_SATA
)

type SMARTTestResult struct {
	TestResult         string
	MarginalAttributes []SMARTAttribute
}

type SMARTAttribute struct {
	ID         int
	Name       string
	Flag       string
	Value      int
	Worst      int
	Threshold  int
	Type       string
	Updated    string
	WhenFailed string
	RawValue   string
}

type DriveHealthInfo struct {
	DeviceName           string // e.g., sda
	DeviceModel          string
	SerialNumber         string
	PowerOnHours         uint64
	PowerCycleCount      uint64
	ReallocatedSectors   uint64 // HDD
	ReallocateNANDBlocks uint64 // SSD
	WearLevelingCount    uint64 // SSD/NVMe
	UncorrectableErrors  uint64
	PendingSectors       uint64 // HDD
	ECCRecovered         uint64
	UDMACRCErrors        uint64
	TotalLBAWritten      uint64
	TotalLBARead         uint64
	IsSSD                bool
	IsNVMe               bool
	IsHealthy            bool //true if the test Passed
}
