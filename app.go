package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	folderselector "synk/folder_selector"
	"synk/utils"
	"time"
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

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// TODO: Add a check for a saved value for the shared directory
	utils.LANDiscovery()

	// remove when done
	var TEMP_TEST_DIRECTORY = "test_shared_dir_local"
	var TEMP_TEST_DIRECTORY_REMOTE = "test_shared_dir_remote"
	// ttestsetestsetetestset
	scanSharedDirectory(TEMP_TEST_DIRECTORY)

	compareSharedDirectories(scanSharedDirectory(TEMP_TEST_DIRECTORY), scanSharedDirectory(TEMP_TEST_DIRECTORY_REMOTE))

	// runScan()
	// fmt.Println("scan done?")

}



// Take in the user's command
// Return the new directory, as well as the contents of the directory
// nextFolder : if a folder was selected to be entered, specify which one
func (a *App) FolderSelectorControl (currentDir string, command folderselector.FolderSelectorCommand, nextFolder string) folderselector.FolderSelectorResult {
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
	}

	fmt.Println(output)

	return output
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

var shared_directory = ""
// to be called by the frontend, sets the value of the shared directory in the backend
func (a *App) SetSharedDirectory(dir string) {
	shared_directory = dir
	fmt.Println("INFO: Shared directory has been set to: ", dir)
}

func (a *App) RunSynk() {
	// check to make sure shared_directory is not empty
	if shared_directory == "" {
		fmt.Println("ERROR: Shared Directory has not been specified!")
	} else {
		// TODO: Actually implement the synchronizing over the network.
		RunScan()
	}
}

// Given the directory that the user has chosen as their "shared directory",
// return a dictionary where the key is the file name, and the value is the last modified time
func scanSharedDirectory(dir string) map[string]time.Time {

	output := make(map[string]time.Time)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(info.ModTime())
		output[strings.Replace(path, dir, "SYNK_ROOT_DIRECTORY", 1)] = info.ModTime()
		return nil
	})
	if err != nil {
		panic(err)
	}

	// fmt.Println(output)
	fmt.Println(utils.DirMapToString(output))

	return output
}

// Given: local directory contents, remote directory contents
// Return: Files that need to be requested from remote peer, files that the remote peer needs
//
//	to request from the local peer
func compareSharedDirectories(localDir map[string]time.Time, remoteDir map[string]time.Time) map[string][]string {
	output := make(map[string][]string)
	send_list, receive_list := make([]string, 0), make([]string, 0)

	// get the keys from each input mapping
	// local_file_names, remote_file_names := maps.Keys(localDir), maps.Keys(remoteDir)

	// if a key is only present locally, add it to the send list
	// if a key is only present remotely, add it to the receive list
	// if a key is present in both, compare the value in the mapping
	//		if the local version is newer, add it to the send list
	//		if the remote version is newer, add it to the receive list

	only_local, only_remote, in_both := utils.CompareKeys(localDir, remoteDir)
	fmt.Println("Only local: ", only_local)
	fmt.Println("Only remote: ", only_remote)
	fmt.Println("Both: ", in_both)

	send_list = append(send_list, only_local...)
	receive_list = append(receive_list, only_remote...)

	for _, fileName := range in_both {
		if localDir[fileName].Before(remoteDir[fileName]) {
			send_list = append(send_list, fileName)
		} else if localDir[fileName].After(remoteDir[fileName]) {
			receive_list = append(receive_list, fileName)
		}
	}

	// finally, update the output var
	output["SEND"] = send_list
	output["RECEIVE"] = receive_list

	return output
}

// RunScan scans the local network for other Synk users
func RunScan() {

}
