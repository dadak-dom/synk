package folderselector

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"synk/config"
)

type FolderSelectorCommand int

const (
	MoveUp FolderSelectorCommand = iota
	MoveDown
	GoHome
	Init
	Select
)

type FolderSelectorResult struct {
	Directory string
	Files []string
	Folders []string
}

func ListFilesInDirectory(dir string) []string {
	output := make([]string, 0)
	if dir == "" {
		return output
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("FATAL ERROR in ListFilesInDirectory: ", err, "\n\tTried opening: ", dir)
	}
	for _, file := range files {
		if !file.IsDir() {
			output = append(output, file.Name())
		}
	}
	sort.Slice(output, func(i, j int) bool { return strings.ToLower(output[i]) < strings.ToLower(output[j]) })
	return output
}

func ListFoldersInDirectory(dir string) []string {
	output := make([]string, 0)
	if dir == "" {
		return output
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("FATAL ERROR in ListFoldersInDirectory: ", err, "\n\tTried opening: ", dir)
	}
	for _, file := range files {
		if file.IsDir() {
			output = append(output, file.Name())
		}
	}
	sort.Slice(output, func(i, j int) bool { return strings.ToLower(output[i]) < strings.ToLower(output[j]) })
	return output
}

func InitializeFolderSelector() FolderSelectorResult {
	// startDir, err := os.UserHomeDir()
	startDir := ""
	// if err != nil {
	// 	fmt.Println("FATAL ERROR", err)
	// }
	// return FolderSelectorResult{Directory: startDir, Files: ListFilesInDirectory(startDir)}
	c := config.GetConfigValue(config.SharedDirectory)
	if c != "" {
		startDir = c
	}

	return FolderSelectorResult{Directory: startDir, Files: ListFilesInDirectory(startDir), Folders: ListFoldersInDirectory(startDir)}
}

// TODO: These functions (MoveUp, down, etc.) can probably be abstracted, or at least have some components abstracted to reduce repetition

func MoveUpDir(currentDir string) FolderSelectorResult {
	newPath := filepath.Join(currentDir, "..")
	// newFiles, err := os.ReadDir(newPath)
	// if err != nil {
	// 	fmt.Println("FATAL ERROR: ", err)
	// }
	// newFilesString := make([]string, 0)

	// for _, file := range newFiles {
	// 	if file.IsDir() {
	// 		newFilesString = append(newFilesString, file.Name())
	// 	}
	// }

	output := FolderSelectorResult{Directory: newPath, Files: ListFilesInDirectory(newPath), Folders: ListFoldersInDirectory(newPath)}
	return output
}

func MoveDownDir(currentDir string, newFolder string) FolderSelectorResult {
	newPath := filepath.Join(currentDir, newFolder)
	// newFiles, err := os.ReadDir(newPath)
	// 	if err != nil {
	// 	fmt.Println("FATAL ERROR: ", err)
	// }

	// newFilesString := make([]string, 0)

	// for _, file := range newFiles {
	// 	if file.IsDir() {
	// 		newFilesString = append(newFilesString, file.Name())
	// 	}
	// }

	output := FolderSelectorResult{Directory: newPath, Files: ListFilesInDirectory(newPath), Folders: ListFoldersInDirectory(newPath)}
	return output
}

func GoToHomeDir() FolderSelectorResult {
	startDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("FATAL ERROR", err)
	}
	return FolderSelectorResult{Directory: startDir, Files: ListFilesInDirectory(startDir), Folders: ListFoldersInDirectory(startDir)}
}

// Once a folder is selected, it's saved to the config folder
func SelectSharedFolder(dir string) {
	config.UpdateUserConfig(config.SharedDirectory, dir)
}
