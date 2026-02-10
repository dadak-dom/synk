// All functions related to sending things over the network

package network

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"synk/config"
	"synk/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/schollz/peerdiscovery"
)

// function that sends a file
// func SendFile(file_name string, ipaddr IP) {

// }

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
				if strings.HasPrefix(ip, "192.168") ||  strings.HasPrefix(ip, "172."){
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
	sharedDirectoryInfo := utils.ScanSharedDirectory(config.GetConfigValue(config.SharedDirectory))
	c.IndentedJSON(http.StatusOK, sharedDirectoryInfo)
}

func UploadFile(c *gin.Context) {
	file, _ := c.FormFile("file")
	dir := c.PostForm("dir")
	// check if the folder exists - if not, create it
	temp := filepath.Dir(config.ConstructCompleteFilePath(dir))
	if _, err := os.Stat(temp); err != nil {
		log.Println("Directory missing, creating now...", temp)
		os.MkdirAll(temp, os.ModePerm)
	}
	log.Println("UPLOAD FILE RECEIVED: ", file.Filename, " IN DIRECTORY: ", dir)

	// FIXME : make this save to the real shared dir
	// c.SaveUploadedFile(file, "C:\\Users\\dadak\\Desktop\\personal-projects\\synk\\test_shared_dir_remote\\"+file.Filename)
	// c.SaveUploadedFile(file, "/home/dominik/synk/test_shared_dir_remote"+file.Filename)
	log.Println("SAVING TO: ", config.ConstructCompleteFilePath(dir))
	c.SaveUploadedFile(file, config.ConstructCompleteFilePath(dir))
	// c.SaveUploadedFile(file, config.GetConfigValue())
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

// Allows peer to download a file, provided they give a file index
func GetFile(c *gin.Context) {
	// i := c.Param("index")
	i, _ := strconv.Atoi(c.Query("index"))
	// get list of files in shared folder
	sharedFiles := utils.ListFilesInSharedDirectory(config.GetConfigValue(config.SharedDirectory))
	fmt.Println("FILE THAT I WOULD SEND: , ", sharedFiles[i], "i = ", i, "param=", c.Query("index"))
	// I need to be able to construct the path to the file in question
	complete_path := config.ConstructCompleteFilePath(sharedFiles[i])
	fmt.Println("============================\nCOMPLETE PATH: ", complete_path)

	// c.String(http.StatusOK, "HELLO")
	c.File(complete_path)
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
