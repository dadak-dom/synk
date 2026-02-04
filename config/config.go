package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
	WriteTextFile(ConfigLocation[updated_item], string(SharedDirectory), value)
}

func GetConfigValue(value ConfigItem) string {
	r := ReadTextFile(ConfigLocation[value], string(value))
	fmt.Println("Config value for: ", value, r)
	return r
}

func WriteTextFile(dir string, fileName string, content string) {
	err := os.WriteFile(filepath.Join(dir, fileName), []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("File ", fileName, " written successfully.")
}

func ReadTextFile(dir string, fileName string) string {
	content, err := os.ReadFile(filepath.Join(dir, fileName))
	if err != nil {
		log.Println("File reading error ", err)
		return ""
	}

	return string(content)
}

// given the ending filepath (e.g. SYNK_ROOT_DIRECTORY/test.txt), get the full path
// 	(e.g. /home/user/test.txt)
func ConstructCompleteFilePath(ending string) string {
	o := strings.Replace(ending, "SYNK_ROOT_DIRECTORY", GetConfigValue(SharedDirectory), 1)
	// if running on windows, reverse the path cleaning
	if runtime.GOOS == "windows" {
		o = strings.Replace(o, "/", "\\", -1)
	}
	return o
}
