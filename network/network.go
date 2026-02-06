// All functions related to sending things over the network

package network

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
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
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func GetSharedFolderInfo(c *gin.Context) {
	sharedDirectoryInfo := utils.ScanSharedDirectory(config.GetConfigValue(config.SharedDirectory))
	// sharedDirectoryInfo := utils.ScanSharedDirectory("/home/dominik/synk/test_shared_dir_remote")
	// names, data := utils.ConvertSharedDirectoryMapToLists(sharedDirectoryInfo)
	// fmt.Println("TESTING SORT:")
	// fmt.Println("\tNAMES: ", names)
	// fmt.Println("\tDATA: ", data)
	c.IndentedJSON(http.StatusOK, sharedDirectoryInfo)
}

func UploadFile(c *gin.Context) {
	file, _ := c.FormFile("file")
	log.Println("UPLOAD FILE RECEIVED: ", file.Filename)

	// FIXME : make this save to the real shared dir
	// c.SaveUploadedFile(file, "C:\\Users\\dadak\\Desktop\\personal-projects\\synk\\test_shared_dir_remote\\"+file.Filename)
	// c.SaveUploadedFile(file, "/home/dominik/synk/test_shared_dir_remote"+file.Filename)
	c.SaveUploadedFile(file, config.ConstructCompleteFilePath(file.Filename))
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
	log.Println("Running test of LAN Discovery...")
	peers := make([]string, 0)
	// discover peers
	discoveries, err := peerdiscovery.Discover(peerdiscovery.Settings{
		Limit:     -1,
		Payload:   []byte("test"),
		Delay:     100 * time.Millisecond,
		TimeLimit: 3 * time.Second,
		Notify: func(d peerdiscovery.Discovered) {
			log.Println(d)

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
