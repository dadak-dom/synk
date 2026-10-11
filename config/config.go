package config

import (
	"bytes"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/simonfrey/jsonl"
)

type ConfigItem string

type ConfigValue interface {
	~string | ~[]string
}

const (
	SharedDirectory  ConfigItem = "shared_directory.txt"
	FileIgnoreList   ConfigItem = "file_ignore.jsonl"
	FolderIgnoreList ConfigItem = "folder_ignore.jsonl"
	Theme            ConfigItem = "theme.txt"
	TrustedNetworks  ConfigItem = "trusted_networks.jsonl"
	TrustedDevices   ConfigItem = "trusted_devices.jsonl"
	EnableAutoSynk   ConfigItem = "enable_auto_synk.txt"
	DeviceName       ConfigItem = "device_name.txt"
	//TODO: add more as needed
)

var AllConfigItems = []ConfigItem{
	SharedDirectory,
	FileIgnoreList,
	FolderIgnoreList,
	Theme,
	TrustedNetworks,
	TrustedDevices,
	EnableAutoSynk,
	DeviceName,
}

var textConfigItems = make([]ConfigItem, 0)
var jsonLinesConfigItems = make([]ConfigItem, 0)

// Every time the app starts up, check to see if the most important config items have been set. If not, set them.
func ConfigInitSetup() {
	for _, config_item := range AllConfigItems {
		if strings.HasSuffix(string(config_item), ".txt") {
			textConfigItems = append(textConfigItems, config_item)
		} else if strings.HasSuffix(string(config_item), ".jsonl") {
			jsonLinesConfigItems = append(jsonLinesConfigItems, config_item)
		}
	}
	if !isConfigItemSet(SharedDirectory) {
		UpdateUserConfigString(SharedDirectory, os.TempDir()) // if the value is missing, set it to the temp dir
	}
	if !isConfigItemSet(DeviceName) {
		UpdateUserConfigString(DeviceName, generateRandomDeviceName())
	}
	log.Println("Current device name: ", GetConfigValueString(DeviceName))
}

// Get the config file location
func configSetup() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("Fatal error when setting up config directory: ", err)
	}
	path := filepath.Join(dir, "synk")
	err = os.MkdirAll(path, 0755)
	if err != nil {
		log.Fatal("Fatal error when setting up config directory: ", err)
	}
	// make all files that need to exist
	for _, f := range AllConfigItems {
		p := filepath.Join(path, string(f))
		_, e := os.Stat(p)
		if os.IsNotExist(e) {
			os.WriteFile(p, make([]byte, 0), 0755)
		}
	}
	return path
}

func isStringlistEmpty(slice []string) bool {
	return len(slice) == 0
}

func isStringEmpty(s string) bool {
	return s == ""
}

func isConfigItemSet(config_item ConfigItem) bool {
	if slices.Contains(textConfigItems, config_item) {
		return !isStringEmpty(GetConfigValueString(config_item))
	} else if slices.Contains(jsonLinesConfigItems, config_item) {
		return !isStringlistEmpty(GetConfigValueStringList(config_item))
	}
	return false
}

func GetConfigItemFileLocation(ci ConfigItem) string {
	dir, _ := os.UserConfigDir()
	path := filepath.Join(dir, "synk")
	return filepath.Join(path, string(ci))
}

var ConfigLocation = configSetup()

func initRand() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// return a randomly named file
func RandomFileName(file_extension string) string {
	return randStringRunes(42) + file_extension
}

func UpdateUserConfigString(updated_item ConfigItem, value string) {
	WriteTextFile(ConfigLocation, string(updated_item), value)
}

func UpdateUserConfigStringList(updated_item ConfigItem, value []string) {
	writeJsonLinesFile(ConfigLocation, string(updated_item), value)
}

func writeJsonLinesFile(dir string, fileName string, values []string) {
	buff := bytes.Buffer{}
	w := jsonl.NewWriter(&buff)
	for _, v := range values {
		w.Write(v)
	}
	if err := os.WriteFile(filepath.Join(dir, fileName), buff.Bytes(), 0644); err != nil {
		log.Fatal("Error when writing JSON lines config: ", err)
	}
}

func GetConfigValueString(value ConfigItem) string {
	// switch value {
	// case SharedDirectory, Theme, EnableAutoSynk, DeviceName:
	// 	r := ReadTextFile(ConfigLocation, string(value))
	// 	log.Println("Config value (string) for: ", value, r)
	// 	return r
	// default:
	// 	log.Fatal("Missing case in GetConfigValueString (did you forget to update the case?)")
	// }
	// return ""
	if !slices.Contains(textConfigItems, value) {
		log.Fatal("Fatal error: GetConfigValueString called on config item that does not return a string")
		return ""
	} else {
		return readTextFile(value)
	}
}

func GetConfigValueStringList(value ConfigItem) []string {
	// switch value {
	// // TODO: if more cases come, add them here
	// case FolderIgnoreList, FileIgnoreList, TrustedNetworks:
	// 	log.Println("Config value (stringList) for: ", value)
	// 	return readJsonLinesFile(value)
	// default:
	// 	log.Fatal("Missing case in GetConfigValueStringList (did you forget to update the case?)")
	// }
	if !slices.Contains(jsonLinesConfigItems, value) {
		log.Fatal("Fatal error: GetConfigValueStringList called on config item that does not return a list of strings")
		return make([]string, 0)
	} else {
		return readJsonLinesFile(value)
	}

}

func WriteTextFile(dir string, fileName string, content string) {
	err := os.WriteFile(filepath.Join(dir, fileName), []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("File ", fileName, " written successfully.")
}

func readTextFile(value ConfigItem) string {
	content, err := os.ReadFile(filepath.Join(ConfigLocation, string(value)))
	if err != nil {
		log.Println("File reading error ", err)
		return ""
	}

	return string(content)
}

func readJsonLinesFile(value ConfigItem) []string {
	content, err := os.ReadFile(filepath.Join(ConfigLocation, string(value)))
	if err != nil {
		log.Fatal("Error reading json lines file: ", err)
	}
	r := jsonl.NewReader(strings.NewReader(string(content)))
	o := make([]string, 0)
	line := ""
	r.ReadSingleLine(&line)
	for line != "" {
		o = append(o, line)
		line = ""
		r.ReadSingleLine(&line)
	}
	return o
}

func OpenJsonLinesConfigFile(dir string) []string {
	return readJsonLinesFile(ConfigItem(dir))
}

// given the ending filepath (e.g. SYNK_ROOT_DIRECTORY/test.txt), get the full path
//
//	(e.g. /home/user/test.txt)
func ConstructCompleteFilePath(ending string) string {
	s := GetConfigValueString(SharedDirectory)
	o := strings.Replace(ending, "SYNK_ROOT_DIRECTORY", s, 1)

	return enforceCorrectSlashDirection(o)
}

// sometimes files that are sent over don't have the correct slashes... this will change them to the appropriate direction
func enforceCorrectSlashDirection(s string) string {
	// if running on windows, reverse the path cleaning
	if runtime.GOOS == "windows" {
		s = strings.Replace(s, "/", "\\", -1)
	} else {
		s = strings.Replace(s, "\\", "/", -1)
	}
	return s
}

func generateRandomDeviceName() string {
	possible_words := [20]string{
		"bilbao",
		"aardvark",
		"abandon",
		"studio",
		"elvish",
		"built",
		"prized",
		"technical",
		"werewolf",
		"backhand",
		"stucco",
		"zoo",
		"correct",
		"essay",
		"objective",
		"lance",
		"knight",
		"sunrise",
		"pasture",
		"meadow",
	}
	output := ""
	initRand()
	output = possible_words[rand.Intn(20)] + "-" + possible_words[rand.Intn(20)] + "-" + possible_words[rand.Intn(20)]
	return output
}
