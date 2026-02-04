package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func compareKeys[K comparable, V1 any, V2 any](
	m1 map[K]V1,
	m2 map[K]V2,
) (onlyInM1, onlyInM2, inBoth []K) {

	seen := make(map[K]struct{}, len(m1))

	// Walk m1
	for k := range m1 {
		if _, ok := m2[k]; ok {
			inBoth = append(inBoth, k)
		} else {
			onlyInM1 = append(onlyInM1, k)
		}
		seen[k] = struct{}{}
	}

	// Walk m2
	for k := range m2 {
		if _, ok := seen[k]; !ok {
			onlyInM2 = append(onlyInM2, k)
		}
	}

	return
}

func DirMapToString(m map[string]time.Time) string {
	output := "Filename -> DateMod mapping: \n"

	for k, v := range m {
		output = output + "\tName: " + k + " Date modified: " + v.String() + "\n"
	}

	return output
}

// Given the directory that the user has chosen as their "shared directory",
// return a dictionary where the key is the file name, and the value is the last modified time
func ScanSharedDirectory(dir string) map[string]time.Time {

	output := make(map[string]time.Time)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// make sure that folders aren't being tracked
		if info.IsDir() {
			return nil
		}
		fmt.Println(info.ModTime())
		output[strings.Replace(path, dir, "SYNK_ROOT_DIRECTORY", 1)] = info.ModTime()
		return nil
	})
	if err != nil {
		panic(err)
	}

	delete(output, "SYNK_ROOT_DIRECTORY") // Root folder info doesn't matter
	fmt.Println(DirMapToString(output))

	return output
}

// given a slice and a value, return the index of the first match
func IndexOf[K ~[]E, E comparable](s K, v E) int {
	for i := range s {
		if v == s[i] {
			return i
		}
	}
	return -1
}

// Convert the map of file->data to 2 sorted lists
func ConvertSharedDirectoryMapToLists(m map[string]time.Time) ([]string, []time.Time) {
	fileNames, fileMetadata := make([]string, 0), make([]time.Time, 0)
	for k := range(cleanFileNames(m)) {
		fileNames = append(fileNames, k)
		fileMetadata = append(fileMetadata, m[k])
	}
	// sort the filenames and metadata in parallel - whatever I change in the names, change in the metadata
	var min_index int
	var temp_filenames string
	var temp_meta time.Time
	size := len(fileNames)
	for i := 0; i < size; i++ {
		min_index = i
  		for j := i + 1; j < size; j++ {
			if fileNames[j] < fileNames[min_index] {
				min_index = j
			}
     	 }
      temp_filenames = fileNames[i]
	  temp_meta = fileMetadata[i]
      fileNames[i] = fileNames[min_index]
	  fileMetadata[i] = fileMetadata[min_index]
      fileNames[min_index] = temp_filenames
	  fileMetadata[min_index] = temp_meta
	}

	return fileNames, fileMetadata
}

// To remove confusion across OS's, replace \'s with /'s
func cleanFileNames(m map[string]time.Time) map[string]time.Time{
	o := make(map[string]time.Time)
	for fileName := range m {
		new_key := strings.Replace(fileName, "\\", "/", -1)
		o[new_key] = m[fileName]

	}
	return o
}

func ListFilesInSharedDirectory(dir string) []string {
	o, _ := ConvertSharedDirectoryMapToLists(ScanSharedDirectory(dir))
	return o
}


// Given: local directory contents, remote directory contents
// Return: Files that need to be requested from remote peer, files that the remote peer needs
//
//	to request from the local peer
// FIXME: Make sure that this works across OS's (i.e. swap out the \'s for /'s)
func CompareSharedDirectories(localDir map[string]time.Time, remoteDir map[string]time.Time) map[string][]string {
	output := make(map[string][]string)
	send_list, receive_list := make([]string, 0), make([]string, 0)

	// get the keys from each input mapping
	// local_file_names, remote_file_names := maps.Keys(localDir), maps.Keys(remoteDir)

	// if a key is only present locally, add it to the send list
	// if a key is only present remotely, add it to the receive list
	// if a key is present in both, compare the value in the mapping
	//		if the local version is newer, add it to the send list
	//		if the remote version is newer, add it to the receive list

	localDirCleaned := cleanFileNames(localDir)
	remoteDirCleaned := cleanFileNames(remoteDir)

	only_local, only_remote, in_both := compareKeys(localDirCleaned, remoteDirCleaned)
	fmt.Println("Only local: ", only_local)
	fmt.Println("Only remote: ", only_remote)
	fmt.Println("Both: ", in_both)

	send_list = append(send_list, only_local...)
	receive_list = append(receive_list, only_remote...)

	for _, fileName := range in_both {
		if localDirCleaned[fileName].Before(remoteDirCleaned[fileName]) {
			receive_list = append(receive_list, fileName)
		} else if localDirCleaned[fileName].After(remoteDirCleaned[fileName]) {
			send_list = append(send_list, fileName)
		}
	}

	// finally, update the output var
	output["SEND"] = send_list
	output["RECEIVE"] = receive_list

	return output
}
