package config

import (
	"fmt"
	"synk/utils"
)

type ConfigItem string

const (
	SharedDirectory ConfigItem = "shared_directory.txt"
	//TODO: add more as needed
)

// describes where to find the files for each config item
var ConfigLocation = map[ConfigItem]string{
	SharedDirectory: "./config",
}

func UpdateUserConfig(updated_item ConfigItem, value string) {
	// save shared directory
	//TODO: Add more as needed...
	utils.WriteTextFile(ConfigLocation[updated_item], string(SharedDirectory), value)
}

func GetConfigValue(value ConfigItem) string {
	r := utils.ReadTextFile(ConfigLocation[value], string(value))
	fmt.Println("Config value for: ", value, r)
	return r
}