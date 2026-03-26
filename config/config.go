package config

import (
	"bytes"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
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
	Theme ConfigItem = "theme.txt"
	//TODO: add more as needed
)

var AllConfigItems = []ConfigItem{
	SharedDirectory,
	FileIgnoreList,
	FolderIgnoreList,
	Theme,
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

// return a randomly named file in the config location
func RandomFileName(file_extension string) string {
	return filepath.Join(ConfigLocation, randStringRunes(42)+file_extension)
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
	switch value {
	case SharedDirectory, Theme:
		r := ReadTextFile(ConfigLocation, string(value))
		log.Println("Config value for: ", value, r)
		return r
	default:
		log.Fatal("Missing case in GetConfigValueString")
	}
	return ""
}

func GetConfigValueStringList(value ConfigItem) []string {
	switch value {
	// TODO: if more cases come, add them here
	case FolderIgnoreList, FileIgnoreList:
		log.Println("Config value for: ", value)
		return readJsonLinesFile(value)
	default:
		log.Fatal("Missing case in GetConfigValueString")
	}
	return make([]string, 0)
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
	// if running on windows, reverse the path cleaning
	if runtime.GOOS == "windows" {
		o = strings.Replace(o, "/", "\\", -1)
	}
	return o
}
