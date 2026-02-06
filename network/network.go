// All functions related to sending things over the network

package network

import (
	"fmt"
	"log"
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

func GetSharedFolderInfo(c *gin.Context) {
	sharedDirectoryInfo := utils.ScanSharedDirectory(config.GetConfigValue(config.SharedDirectory))
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
	c.SaveUploadedFile(file, "C:\\Users\\dadak\\Desktop\\personal-projects\\synk\\test_shared_dir_remote\\"+file.Filename)
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
func LANDiscovery() {
	log.Println("Running test of LAN Discovery...")
	// discover peers
	discoveries, err := peerdiscovery.Discover(peerdiscovery.Settings{
		Limit:     -1,
		Payload:   []byte("test"),
		Delay:     100 * time.Millisecond,
		TimeLimit: 30 * time.Second,
		Notify: func(d peerdiscovery.Discovered) {
			log.Println(d)
		},
		MulticastAddress: "239.255.255.250",
	})

	// print out results
	if err != nil {
		log.Fatal(err)
	} else {
		if len(discoveries) > 0 {
			fmt.Printf("Found %d other computers\n", len(discoveries))
			for i, d := range discoveries {
				fmt.Printf("%d) '%s' with payload '%s'\n", i, d.Address, d.Payload)
			}
		} else {
			fmt.Println("Found no devices. You need to run this on another computer at the same time.")
		}
	}
}
