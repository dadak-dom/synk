package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"synk/utils"
	"time"

	nmap "github.com/Ullaakut/nmap/v3"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// remove when done
	var TEMP_TEST_DIRECTORY = "test_shared_dir_local"
	var TEMP_TEST_DIRECTORY_REMOTE = "test_shared_dir_remote"
	// ttestsetestsetetestset
	scanSharedDirectory(TEMP_TEST_DIRECTORY)

	compareSharedDirectories(scanSharedDirectory(TEMP_TEST_DIRECTORY), scanSharedDirectory(TEMP_TEST_DIRECTORY_REMOTE))

	// runScan()
	// fmt.Println("scan done?")

}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
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

// runScan scans the local network
func runScan() {
	// Equivalent to
	// nmap -sV -T4 192.168.0.0/24 with a filter to remove non-RTSP ports.
	scanner, err := nmap.NewScanner(
		context.Background(),
		nmap.WithTargets("192.168.0.0/24"),
		nmap.WithPorts("80", "554", "8554"),
		nmap.WithServiceInfo(),
		nmap.WithTimingTemplate(nmap.TimingAggressive),
		// Filter out ports that are not RTSP
		// nmap.WithFilterPort(func(p nmap.Port) bool {
		// 	return p.Service.Name == "rtsp"
		// }),
		// Filter out hosts that don't have any open ports
		nmap.WithFilterHost(func(h nmap.Host) bool {
			// Filter out hosts with no open ports.
			for idx := range h.Ports {
				if h.Ports[idx].Status() == "open" {
					return true
				}
			}

			return false
		}),
	)
	if err != nil {
		log.Fatalf("unable to create nmap scanner: %v", err)
	}

	result, warnings, err := scanner.Run()
	if len(*warnings) > 0 {
		log.Printf("run finished with warnings: %s\n", *warnings) // Warnings are non-critical errors from nmap.
	}
	if err != nil {
		log.Fatalf("nmap scan failed: %v", err)
	}

	for _, host := range result.Hosts {
		fmt.Printf("Host %s\n", host.Addresses[0])

		for _, port := range host.Ports {
			fmt.Printf("\tPort %d open\n", port.ID)
		}
	}

}
