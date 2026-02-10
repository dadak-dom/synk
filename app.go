package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"synk/config"
	folderselector "synk/folder_selector"
	"synk/network"
	"synk/utils"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

var selectedFolder, err = os.UserHomeDir()

var peerList = make([]string, 0)

func updatePeerList(p []string) {
	log.Println("Updating peer list: ", p)
	peerList = p
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	// msg := <-peers
	// log.Println("PEERS:", msg)

	// log.Fatal("DONE")
	a.ctx = ctx
	router := gin.Default()
	// router.SetTrustedProxies([]string{})
	updates := make(chan []string)

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			updates <- network.LANDiscovery()
		}
	}()

	go func() {
		for peers := range updates {
			log.Println("Updated peers:", peers)
			updatePeerList(peers)
		}
	}()

	// FIXME: Consider setting up the API when starting a transfer, and then shutting it down when it's done
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST"},
	}))

	log.Println("LOCAL IP INFO: ", network.GetLocalIP())
	myLocalIP := network.GetLocalIP()
	APIport := ":8080"
	if myLocalIP == "" {
		log.Fatal("Error: could not find local IP address (192.168 address)")
	}

	// Send information about the shared folder to the "active" peer
	router.GET("/getSharedFolder", network.GetSharedFolderInfo)
	router.GET("/getFile", network.GetFile)
	router.POST("/uploadFile", network.UploadFile)
	router.Run(myLocalIP + APIport)
	// TODO: Use the following guide to figure out how to send files back and forth:
	// https://gin-gonic.com/en/docs/examples/upload-file/single-file/

	// TODO: Add a check for a saved value for the shared directory

}

func (a *App) TestLANDiscovery() {
	network.LANDiscovery()
}

// Take in the user's command
// Return the new directory, as well as the contents of the directory
// nextFolder : if a folder was selected to be entered, specify which one
func (a *App) FolderSelectorControl(currentDir string, command folderselector.FolderSelectorCommand, nextFolder string) folderselector.FolderSelectorResult {
	var output folderselector.FolderSelectorResult

	switch command {
	case folderselector.Init:
		//TODO: Add a feature that makes the program remember the user's selection
		output = folderselector.InitializeFolderSelector()
		log.Println("test", output)
		// selectedFolder = output.Directory
	case folderselector.GoHome:
		output = folderselector.GoToHomeDir()
	case folderselector.MoveUp:
		output = folderselector.MoveUpDir(currentDir)
		// selectedFolder = output.Directory
	case folderselector.MoveDown:
		output = folderselector.MoveDownDir(currentDir, nextFolder)
		// selectedFolder = output.Directory
	case folderselector.Select:
		folderselector.SelectSharedFolder(currentDir)
		output = folderselector.FolderSelectorResult{Directory: "", Files: make([]string, 0)}
	}

	fmt.Println(output)

	return output
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Returns the URL of the remote client that we are connecting to
func (a *App) GetPeerList() []string {
	// TODO: Actually implement this. Will probably need some global
	// variable that tracks what devices we're connected with
	// return []string{"192.168.0.235:8080"}
	// return network.LANDiscovery()
	return peerList
}

func (a *App) RunSynkOnPeer(connection string, peerFileInfo map[string]time.Time) {
	// fmt.Println("Running RunSynk on: ", connection, filesToReceive, filesToSend)
	// fmt.Println("Called RunSynk")
	// fmt.Println("Received: ", connection, peerFileInfo)

	// compare the shared directories
	// FIXME: The code below this comment is correct. Uncomment once peer discovery works
	local_shared_folder := config.GetConfigValue(config.SharedDirectory)
	// local_shared_folder := "test_shared_dir_local"
	comparison := utils.CompareSharedDirectories(utils.ScanSharedDirectory(local_shared_folder), peerFileInfo)
	filesToSend, filesToReceive := comparison["SEND"], comparison["RECEIVE"]
	fmt.Println("Files to send: ", filesToSend)
	fmt.Println("Files to receive", filesToReceive)
	log.Println("\n=========================\nCOMPARISON RESULTS\n========================\nLOCAL:")
	log.Println(utils.DirMapToString(utils.ScanSharedDirectory(local_shared_folder)))
	log.Println("\n=========================\nCOMPARISON RESULTS\n========================\nREMOTE:")
	log.Println(utils.DirMapToString(peerFileInfo))
	log.Println("===========================================")

	// Before downloading/sending anything, make sure that both computers have the necessary folders
	// Use peerFileInfo to extract any directories you need locally, and send directories the remote computer needs to add
	// ==============
	peerFolders := make([]string, 0)
	for f := range peerFileInfo {
		// log.Println(config.ConstructCompleteFilePath(f))
		temp := strings.Replace(f, "SYNK_ROOT_DIRECTORY/", "", 1)
		s := strings.Split(temp, "/")
		if len(s) > 1 {
			peerFolders = append(peerFolders, filepath.Dir(temp)) // if there is a folder in the filename, add it to the list
		}
		log.Println(temp)
	}
	log.Println("Peerfolders: ", peerFolders)
	for _, f := range peerFolders {
		temp := config.ConstructCompleteFilePath("SYNK_ROOT_DIRECTORY/" + f)
		_, err := os.Stat(temp)
		if err != nil { // if the folder doesn't exist, create it
			// log.Println("HERE!", err.Error())
			if err := os.MkdirAll(temp, os.ModePerm); err != nil {
				log.Fatal("Error creating new directories: ", err)
			}

		}
	}

	// Now, local folder structure should be synked with the remote computer.
	// Next, tell the remote computer what folders it needs

	//===============
	// download the files
	// how this works: since this computer (the "active" one) has already received information re: remote files via
	// 	peerFileInfo, convert peerFileInfo into a list and ask for the index of the file you need.
	// 	This works because the remote "passive" peer will also have the same sorted list
	remote_names, _ := utils.ConvertSharedDirectoryMapToLists(peerFileInfo)
	for _, f := range filesToReceive {
		fmt.Println("Receive file: ", f, " with index: ", utils.IndexOf(remote_names, f))
		fileIndex := utils.IndexOf(remote_names, f)
		resp, err := http.Get(connection + "/getFile?index=" + strconv.Itoa(fileIndex))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Status Code: %d\n", resp.StatusCode)
		// fmt.Printf("Body: %s\n", body)
		// FIXME: actually save the file to the shared directory
		// will need to copy files to temp folder, complete operation, then delete the temp folder
		// FIXME: For now, just save the file to the test directory
		fmt.Println("Trying to write to: ", config.ConstructCompleteFilePath(f))
		errWrite := os.WriteFile(config.ConstructCompleteFilePath(f), body, 0644)
		if errWrite != nil {

			fmt.Println("Error when writing:")
			log.Fatal(errWrite)
		}
		// log.Println("File ", filepath.Join("/home/dominik/synk/test_shared_dir_local", filepath.Base(f)), " written successfully.")
		// os.WriteFile(filepath.Join("C:\\Users\\dadak\\Desktop\\personal-projects\\synk\\test_shared_dir_local", filepath.Base(f)), body, 0644)
	}
	// log.Fatal("DONE")

	log.Println("===========================\nRECEIVING DONE, NOW SENDING\n==========================")
	log.Println("Files to send: ", filesToSend)
	// TODO : upload files to peer
	for _, f := range filesToSend {
		// fmt.Println("Send file: ", f, " with index: ", utils.IndexOf(, f))
		// fileIndex := utils.IndexOf(remote_names, f)
		// resp, err := http.Get(connection + "/getFile?index=" + strconv.Itoa(fileIndex))

		// get file that needs to be uploaded to peer
		// FIXME: The line below is correct. Uncomment once done prototyping
		file_content, errReading := os.Open(config.ConstructCompleteFilePath(f))
		// file_content, errReading := os.Open(filepath.Join("test_shared_dir_local", filepath.Base(f)))
		if errReading != nil {
			log.Fatal("Could not open file: ", errReading)
		}

		var requestBody bytes.Buffer
		writer := multipart.NewWriter(&requestBody)

		defer file_content.Close()

		part, err := writer.CreateFormFile("file", filepath.Base(f))
		if err != nil {
			log.Fatal("Error creating form file: ", err)
		}

		_, err = io.Copy(part, file_content)
		if err != nil {
			log.Fatal("Error copying file data: ", err)
		}

		err = writer.WriteField("dir", f)
		if err != nil {
			log.Fatal("Error writing form field: ", err)
		}

		err = writer.Close()
		if err != nil {
			log.Fatal("Error closing writer: ", err)
		}
		url := connection + "/uploadFile"
		req, err := http.NewRequest("POST", url, &requestBody)
		if err != nil {
			log.Fatal("Error creating request: ", err)
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Error sending request: ", err)
		}
		defer resp.Body.Close()

		log.Println("Server responded with status: ", resp.Status)

		// var b bytes.Buffer
		// w := multipart.NewWriter(&b)

		// values := map[string]io.Reader{
		// 	"file": file_content,
		// }

		// log.Println("file_content: ", file_content)
		// FIXME there is a bug somewhere here. The API has been verified to work via curl
		// for key, r := range values {
		// 	var fw io.Writer
		// 	if x, ok := r.(io.Closer); ok {
		// 		defer x.Close()
		// 	}
		// 	if x, ok := r.(*os.File); ok {
		// 		if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
		// 			// return
		// 		} else {
		// 			if fw, err = w.CreateFormField(key); err != nil {
		// 				// return
		// 			}
		// 		}
		// 	}
		// 	if _, err = io.Copy(fw, r); err != nil {
		// 		// return err
		// 		log.Fatal("Error when creating multipart form: ", err)
		// 	}
		// }
		// w.Close()
		// // send the request
		// req, err := http.NewRequest("POST", url, &b)
		// if err != nil {
		// 	log.Fatal("Error when creating POST request", err)
		// }
		// req.Header.Set("Content-Type", w.FormDataContentType())
		// client := &http.Client{}
		// resp, err := client.Do(req)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer resp.Body.Close()
		// body, err := io.ReadAll(resp.Body)
		// fmt.Printf("Status Code: %d\n", resp.StatusCode)
		// fmt.Printf("Response Body: %s\n", body)
	}

}
