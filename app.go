package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"synk/config"
	folderselector "synk/folder_selector"
	"synk/network"
	"synk/utils"
	"time"

	"github.com/fsnotify/fsnotify"
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
	a.ctx = ctx

	// TODO: implement a file watcher
	// In order to update the actual file watcher, we need to watch the shared_directory config item for any changes
	// Setting up file watcher:
	network.UpdateAPIStatus()
	sharedDirectoryUpdates := make(chan string)
	autoSynkUpdates := make(chan string)
	go watchSharedDirConfig(sharedDirectoryUpdates)
	go watchSharedDirContents(sharedDirectoryUpdates, autoSynkUpdates)
	go listenForPeers()
	go startAPI()
	// TODO: implement the below:
	go listenForAutoSynk(autoSynkUpdates)
}

func listenForAutoSynk(fileUpdates <-chan string) {
	// TODO: handle any updates to files in watch directory
	for fu := range fileUpdates {
		autoSynkEnabled := config.GetConfigValueString(config.EnableAutoSynk)
		if autoSynkEnabled == "true" {
			for _, peer := range peerList {
				connection := "http://" + peer + ":8080"
				file_content, errReading := os.Open(fu)

				if errReading != nil {
					log.Fatal("Could not open file: ", errReading)
					// return false
				}

				var requestBody bytes.Buffer
				writer := multipart.NewWriter(&requestBody)

				defer file_content.Close()

				part, err := writer.CreateFormFile("file", filepath.Base(fu))
				if err != nil {
					log.Fatal("Error creating form file: ", err)
					// return false
				}

				_, err = io.Copy(part, file_content)
				if err != nil {
					log.Fatal("Error copying file data: ", err)
					// return false
				}

				err = writer.WriteField("dir", utils.GetStandardizedFileName(fu))
				if err != nil {
					log.Fatal("Error writing form field: ", err)
					// return false
				}

				err = writer.Close()
				if err != nil {
					log.Fatal("Error closing writer: ", err)
					// return false
				}
				url := connection + "/uploadFile"
				req, err := http.NewRequest("POST", url, &requestBody)
				if err != nil {
					log.Fatal("Error creating request: ", err)
					// return false
				}

				req.Header.Set("Content-Type", writer.FormDataContentType())
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					log.Fatal("Error sending request: ", err)
					// return false
				}
				defer resp.Body.Close()

				log.Println("Server responded with status: ", resp.Status)
			}
		} else {
			log.Println("Auto synk disabled, ignoring write...")
		}
	}
}

func watchSharedDirConfig(updates chan<- string) {
	log.Println("Creating watcher for shared directory config item...")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	err = watcher.Add(config.GetConfigItemFileLocation(config.SharedDirectory))
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			log.Println("event:", event)
			if event.Has(fsnotify.Write) {
				log.Println("modified file:", event.Name)
				updates <- config.GetConfigValueString(config.SharedDirectory) // add the event to the shared channel
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func watchSharedDirContents(sharedDirChange <-chan string, updatedFiles chan<- string) {
	log.Println("Creating watcher for shared directory...")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	// debounce state
	const debounceDelay = 250 * time.Millisecond
	timers := make(map[string]*time.Timer)
	var mu sync.Mutex

	sf := config.GetConfigValueString(config.SharedDirectory)
	err = watcher.Add(sf)
	filepath.WalkDir(sf, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			watcher.Add(path)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("WATCHING THE FOLLOWING: ", watcher.WatchList())
	for {
		select {
		case newDir := <-sharedDirChange:
			// sf in this case is the old shared folder
			// remove everything on watchlist
			for _, watched := range watcher.WatchList() {
				err = watcher.Remove(watched)
				if err != nil {
					log.Fatal("Error when removing sub-directories of old watchlist: ", err)
				}
			}

			err = watcher.Add(newDir)
			filepath.WalkDir(newDir, func(path string, d fs.DirEntry, err error) error {
				if d.IsDir() {
					watcher.Add(path)
				}
				return nil
			})
			log.Println("Watchlist has been updated; now watching the following: ", watcher.WatchList())
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			log.Println("event:", event)
			// FIXME: For now , I will only be handling file modifications, not creations/deletions
			if event.Has(fsnotify.Write) {
				log.Println("modified file:", event.Name, "string version: ", event.String())
				// Debouncing:  if timer exists for this file, cancel it
				mu.Lock()
				if t, exists := timers[event.Name]; exists {
					t.Stop()
				}
				filename := event.Name

				// Create a new timer that fires after no changes have occurred
				timers[filename] = time.AfterFunc(debounceDelay, func() {
					log.Println("debounced update:", filename)

					updatedFiles <- filename

					mu.Lock()
					delete(timers, filename)
					mu.Unlock()
				})

				mu.Unlock()
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func listenForPeers() {
	peerUpdates := make(chan []string)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			peerUpdates <- network.LANDiscovery()
			//update whether the API is enabled or not
			// FIXME: there should be a better way to de-couple the API status updates from the landiscovery updates.
			network.UpdateAPIStatus()
		}
	}()

	go func() {
		for peers := range peerUpdates {
			log.Println("Updated peers:", peers)
			updatePeerList(peers)
		}
	}()
}

func startAPI() {
	router := gin.Default()
	// FIXME: Consider setting up the API when starting a transfer, and then shutting it down when it's done
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST"},
	}))

	log.Println("LOCAL IP INFO: ", network.GetLocalIP())
	myLocalIP := network.GetLocalIP()
	APIport := ":8080"
	if myLocalIP == "" {
		log.Fatal("Error: could not find local IP address (192.168 OR 172 address)")
	}

	// Send information about the shared folder to the "active" peer
	router.GET("/getSharedFolder", network.GetSharedFolderInfo)
	router.GET("/getFile", network.GetFile)
	router.POST("/uploadFile", network.UploadFile)
	router.POST("/updateFolderIgnoreList", network.UpdateFolderIgnoreList)
	router.POST("/updateFileIgnoreList", network.UpdateFileIgnoreList)
	router.GET("/resetIgnoreList", network.ResetIgnoreList)
	router.GET("/getFolderIgnoreList", network.GetFolderIgnoreList)
	router.GET("/getFileIgnoreList", network.GetFileIgnoreList)
	router.Run(myLocalIP + APIport)
}

func (a *App) GetLocalIP() string {
	return network.GetLocalIP()
}

func (a *App) GetSharedDirectory() string {
	t := config.GetConfigValueString(config.SharedDirectory)
	return t
}

func (a *App) SetConfigItemString(item config.ConfigItem, value string) {
	log.Println("Updating config item: ", item, " with value ", value)
	config.UpdateUserConfigString(item, value)
}

func (a *App) SetConfigItemStringList(item config.ConfigItem, value []string) {
	log.Println("Updating config item: ", item, " with value ", value)
	config.UpdateUserConfigStringList(item, value)
}

// FIXME: Add documentation here stating the workflow for sending/receiving information to the GO backend via the JS frontend
// In other words: if someone had never seen this codebase before, how could I explain to them how to do these common tasks? (or to myself, if I forget (: )
func (a *App) GetConfigValueString(item config.ConfigItem) string {
	return config.GetConfigValueString(item)
}

func (a *App) GetConfigValueStringList(item config.ConfigItem) []string {
	return config.GetConfigValueStringList(item)
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
	case folderselector.Cancel:
		output = folderselector.CancelDir()
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
	return peerList
}

// Return the current theme in the config
func (a *App) GetTheme() string {
	return config.GetConfigValueString(config.Theme)
}

func ignoreListSynkHelper(peer string) error {
	f1 := filepath.Join(config.ConfigLocation, string(config.FolderIgnoreList))
	f2 := filepath.Join(config.ConfigLocation, string(config.FileIgnoreList))
	// do the same thing for both files and folders
	for i, f := range [2]string{f1, f2} {
		of, err := os.Open(f)
		if err != nil {
			return err
		}

		var requestBody bytes.Buffer
		writer := multipart.NewWriter(&requestBody)

		defer of.Close()

		part, err := writer.CreateFormFile("file", filepath.Base(f))
		if err != nil {
			return err
		}

		_, err = io.Copy(part, of)
		if err != nil {
			return err
		}

		if err != nil {
			return err
		}

		err = writer.Close()
		if err != nil {
			return err
		}
		var url string
		if i == 0 {
			url = peer + "/updateFolderIgnoreList"
		} else {
			url = peer + "/updateFileIgnoreList"
		}

		req, err := http.NewRequest("POST", url, &requestBody)
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		log.Println("Server responded with status: ", resp.Status)
	}

	return nil
}

func (a *App) GetCurrentNetworkName() string {
	return network.GetCurrentNetworkName()
}

func (a *App) RunSynkOnPeer(connection string, peerFileInfo map[string]time.Time) bool {
	// Before comparing any files, make sure that the peer has the same ignore lists
	// Get the peer's ignore lists; add your files to them, send them back
	// peer_folder_ignore_list := http.NewRequest("GET", connection + "/getFolderIgnoreList")
	err := ignoreListSynkHelper(connection)
	if err != nil {
		log.Fatal("Error when propogating ignore lists: ", err)
	}

	local_shared_folder := config.GetConfigValueString(config.SharedDirectory)
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
		fmt.Println("Trying to write to: ", config.ConstructCompleteFilePath(f))
		errWrite := os.WriteFile(config.ConstructCompleteFilePath(f), body, 0644)

		if errWrite != nil {
			fmt.Println("Error when writing:")
			log.Fatal(errWrite)
		}
	}

	log.Println("===========================\nRECEIVING DONE, NOW SENDING\n==========================")
	log.Println("Files to send: ", filesToSend)
	// TODO : upload files to peer
	for _, f := range filesToSend {
		// fmt.Println("Send file: ", f, " with index: ", utils.IndexOf(, f))
		// fileIndex := utils.IndexOf(remote_names, f)
		// resp, err := http.Get(connection + "/getFile?index=" + strconv.Itoa(fileIndex))

		// get file that needs to be uploaded to peer
		// FIXME: Can make this a shared function (also used in auto-synk)
		file_content, errReading := os.Open(config.ConstructCompleteFilePath(f))

		if errReading != nil {
			log.Fatal("Could not open file: ", errReading)
			return false
		}

		var requestBody bytes.Buffer
		writer := multipart.NewWriter(&requestBody)

		defer file_content.Close()

		part, err := writer.CreateFormFile("file", filepath.Base(f))
		if err != nil {
			log.Fatal("Error creating form file: ", err)
			return false
		}

		_, err = io.Copy(part, file_content)
		if err != nil {
			log.Fatal("Error copying file data: ", err)
			return false
		}

		err = writer.WriteField("dir", f)
		if err != nil {
			log.Fatal("Error writing form field: ", err)
			return false
		}

		err = writer.Close()
		if err != nil {
			log.Fatal("Error closing writer: ", err)
			return false
		}
		url := connection + "/uploadFile"
		req, err := http.NewRequest("POST", url, &requestBody)
		if err != nil {
			log.Fatal("Error creating request: ", err)
			return false
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Error sending request: ", err)
			return false
		}
		defer resp.Body.Close()

		log.Println("Server responded with status: ", resp.Status)
	}
	// everything went well, return true
	return true

}
