package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Config struct {
	/* General configuration */
	ChunkRootPath  string `json:"chunk_root_path"`
	FileSystemType string `json:"file_system_type"`
	Mode           string `json:"mode"` // store, index, or hub
	AuthServer     string `json:"auth_server"`

	/* Storage HDD configuration */
	DiskAutoMount  bool   `json:"disk_auto_mount"`
	DiskMountPoint string `json:"disk_mount_point"`
	DiskUUID       string `json:"disk_uuid"`
	DiskDevFile    string `json:"disk_dev_file"`

	/* Cache SSD configuration */
	UseCacheDisk    bool   `json:"use_cache_disk"`
	CacheAutoMount  bool   `json:"cache_auto_mount"`
	CacheMountPoint string `json:"cache_mount_point"`
	CacheUUID       string `json:"cache_uuid"`
	CacheDevFile    string `json:"cache_dev_file"`
}

var configFilePath string
var systemConfig *Config

func init() {
	flag.StringVar(&configFilePath, "config", "config.json", "Path to the configuration file")
	flag.Parse()

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		fmt.Printf("Config file not found: %s\n", configFilePath)
		createDefaultConfig()
		os.Exit(1)
	}

	var err error
	config, err := loadConfig(configFilePath)
	if err != nil {
		fmt.Printf("Error loading config file: %v\n", err)
		os.Exit(1)
	}

	systemConfig = config

}

func createDefaultConfig() {
	defaultConfig := &Config{
		ChunkRootPath:  "/chk",
		FileSystemType: "ext4",
		Mode:           "production",
		AuthServer:     "http://localhost:8080",

		DiskAutoMount:  true,
		DiskMountPoint: "/mnt/disk",
		DiskUUID:       "default-disk-uuid",
		DiskDevFile:    "/dev/sda",

		UseCacheDisk:    true,
		CacheAutoMount:  true,
		CacheMountPoint: "/mnt/cache",
		CacheUUID:       "default-cache-uuid",
		CacheDevFile:    "/dev/sdb",
	}

	err := saveConfig(configFilePath, defaultConfig)
	if err != nil {
		fmt.Printf("Error saving default config file: %v\n", err)
		os.Exit(1)
	}

	systemConfig = defaultConfig
	fmt.Println("Default configuration created.")
}

func saveConfig(filePath string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

func loadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	if systemConfig == nil {

	}
	fmt.Printf("Loaded configuration: %+v\n", systemConfig)
}
