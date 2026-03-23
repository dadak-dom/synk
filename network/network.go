// All functions related to sending things over the network

package network

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"synk/config"
	"synk/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/schollz/peerdiscovery"
)

var enableAPI = false

// Check whether or not to allow the API - i.e. am I connected to a trusted network?
func UpdateAPIStatus() {
	// get trusted networks
	tns := config.GetConfigValueStringList(config.TrustedNetworks)
	// check if current in trusted
	log.Println("RUNNING API STATUS UPDATE", tns, GetCurrentNetworkName())
	if slices.Contains(tns, GetCurrentNetworkName()) {
		enableAPI = true
	} else {
		enableAPI = false
	}
	log.Println("What is api status? (is it enabled?): ", enableAPI)
	// TODO : if the API is not enabled, make every API call return an error
}

// GetLocalIP returns the non loopback local IP of the host
// need to find an IP that is on the same subnet as peers
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.To4().String()
				log.Println("Potential IP: ", ip)
				// if on Windows, don't use 192.168.56.X, as that will give the wrong IP for the API
				if runtime.GOOS == "windows" && ip == "192.168.56.1" { // FIXME: This is a bandaid solution.
					continue
				}
				if strings.HasPrefix(ip, "192.168") || strings.HasPrefix(ip, "172.") {
					return ip
				}
				// return ipnet.IP.String()
			}
		}
	}
	// return "192.168.0.238"
	return ""
}

func GetSharedFolderInfo(c *gin.Context) {
	if !enableAPI {
		c.JSON(http.StatusUnauthorized, "")
	} else {
		sharedDirectoryInfo := utils.ScanSharedDirectory(config.GetConfigValueString(config.SharedDirectory))
		c.JSON(http.StatusOK, sharedDirectoryInfo)
	}
}

func GetFolderIgnoreList(c *gin.Context) {
	if !enableAPI {
		c.JSON(http.StatusUnauthorized, "")
	} else {
		fp := filepath.Join(config.ConfigLocation, string(config.FolderIgnoreList))
		c.File(fp)
	}
}

func GetFileIgnoreList(c *gin.Context) {
	if !enableAPI {
		c.JSON(http.StatusUnauthorized, "")
	} else {
		fp := filepath.Join(config.ConfigLocation, string(config.FileIgnoreList))
		c.File(fp)
	}
}

// receives an individual item to add or remove from the ignore list
// FIXME: next two functions can be abstracted later
func UpdateFileIgnoreList(c *gin.Context) {
	if !enableAPI {
		c.JSON(http.StatusUnauthorized, "")
	} else {
		// step 1: receive the config file & save it to a temp config location
		file, _ := c.FormFile("file")
		temp_file := config.RandomFileName(".jsonl")
		c.SaveUploadedFile(file, temp_file)
		// step 2: compare the remote peer's list to your own & add any files that are not in your list already
		peer_file := config.OpenJsonLinesConfigFile(temp_file)
		local_file := config.GetConfigValueStringList(config.FileIgnoreList)
		new_local_file := slices.Clone(local_file)
		for _, peer_item := range peer_file {
			if !slices.Contains(local_file, peer_item) {
				new_local_file = append(new_local_file, peer_item)
			}
		}
		// step 3: update config file & delete temp
		os.Remove(temp_file)
		config.UpdateUserConfigStringList(config.FileIgnoreList, new_local_file)
	}
}

func UpdateFolderIgnoreList(c *gin.Context) {
	if !enableAPI {
		c.JSON(http.StatusUnauthorized, "")
	} else {
		// step 1: receive the config file & save it to a temp config location
		file, _ := c.FormFile("file")
		temp_file := config.RandomFileName(".jsonl")
		c.SaveUploadedFile(file, temp_file)
		// step 2: compare the remote peer's list to your own & add any files that are not in your list already
		peer_file := config.OpenJsonLinesConfigFile(temp_file)
		local_file := config.GetConfigValueStringList(config.FolderIgnoreList)
		new_local_file := slices.Clone(local_file)
		for _, peer_item := range peer_file {
			if !slices.Contains(local_file, peer_item) {
				new_local_file = append(new_local_file, peer_item)
			}
		}
		// step 3: update config file & delete temp
		os.Remove(temp_file)
		config.UpdateUserConfigStringList(config.FolderIgnoreList, new_local_file)
	}
}

// TODO: functions to add for the "trusted_networks" feature
// Get list of trusted networks
// Add new network to list of trusted
// Remove network from list of trusted
// Return string of the current network - in other words, the SSID / "wifi name"
func GetCurrentNetworkName() string {
	output := ""
	if runtime.GOOS == "windows" {
		out, err := exec.Command("netsh", "wlan", "show", "interfaces").Output()
		if err != nil {
			log.Fatal("Failed to execute command: (in GetCurrentNetworkName)", err)
		}
		re := regexp.MustCompile(`SSID[\s]*: (.*)[\s]*BSSID`)
		firstMatch := re.FindStringSubmatch(string(out))
		if len(firstMatch) == 0 {
			log.Println("Warning: No network name found!")
			output = ""
		} else {
			output = firstMatch[0]
		}
		// post-process the regex - for some reason, the above regex doesn't remove the SSID and BSSID from the string
		temp := strings.Replace(output, "BSSID", "", 1)
		temp2 := strings.Replace(temp, "SSID", "", 1)
		output = strings.Replace(temp2, ":", "", 1) // no idea why, but if I don't split these lines up,the replace function breaks.
		output = strings.Replace(output, " ", "", -1)
		output = strings.Replace(output, "\n", "", -1)
		output = strings.Replace(output, "\r", "", -1)
	}
	// TODO: Add Linux and Mac cases!

	return output
}

func ResetIgnoreList(c *gin.Context) {
	if !enableAPI {
		c.JSON(http.StatusUnauthorized, "")
	} else {
		log.Println("Received request to reset ignore lists...")
		config.UpdateUserConfigStringList(config.FileIgnoreList, make([]string, 0))
		config.UpdateUserConfigStringList(config.FolderIgnoreList, make([]string, 0))
		c.String(http.StatusOK, fmt.Sprintf("Ignore lists reset"))
	}
}

func UploadFile(c *gin.Context) {
	if !enableAPI {
		c.JSON(http.StatusUnauthorized, "")
	} else {
		file, _ := c.FormFile("file")
		dir := c.PostForm("dir")
		// check if the folder exists - if not, create it
		temp := filepath.Dir(config.ConstructCompleteFilePath(dir))
		if _, err := os.Stat(temp); err != nil {
			log.Println("Directory missing, creating now...", temp)
			os.MkdirAll(temp, os.ModePerm)
		}
		log.Println("UPLOAD FILE RECEIVED: ", file.Filename, " IN DIRECTORY: ", dir)

		log.Println("SAVING TO: ", config.ConstructCompleteFilePath(dir))
		c.SaveUploadedFile(file, config.ConstructCompleteFilePath(dir))

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	}
}

// Allows peer to download a file, provided they give a file index
func GetFile(c *gin.Context) {
	if !enableAPI {
		c.JSON(http.StatusUnauthorized, "")
	} else {
		i, _ := strconv.Atoi(c.Query("index"))
		// get list of files in shared folder
		sharedFiles := utils.ListFilesInSharedDirectory(config.GetConfigValueString(config.SharedDirectory))
		fmt.Println("FILE THAT I WOULD SEND: , ", sharedFiles[i], "i = ", i, "param=", c.Query("index"))
		// I need to be able to construct the path to the file in question
		complete_path := config.ConstructCompleteFilePath(sharedFiles[i])
		fmt.Println("============================\nCOMPLETE PATH: ", complete_path)
		c.File(complete_path)
	}
}

// TODO : test to make sure that this works once I get home
func LANDiscovery() []string {
	log.Println("Running LAN Discovery...")
	peers := make([]string, 0)
	// discover peers
	discoveries, err := peerdiscovery.Discover(peerdiscovery.Settings{
		Limit:     -1,
		Payload:   []byte("test"),
		Delay:     100 * time.Millisecond,
		TimeLimit: 3 * time.Second,
		Notify: func(d peerdiscovery.Discovered) {
			// log.Println(d)

		},
		MulticastAddress: "224.0.0.2",
	})

	// print out results
	if err != nil {
		log.Fatal(err)
	} else {
		if len(discoveries) > 0 {
			fmt.Printf("Found %d other computers\n", len(discoveries))
			for i, d := range discoveries {
				fmt.Printf("%d) '%s' with payload '%s'\n", i, d.Address, d.Payload)
				peers = append(peers, d.Address)
			}
		} else {
			fmt.Println("Found no devices. You need to run this on another computer at the same time.")
		}
	}
	return peers
}
