package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/schollz/peerdiscovery"
)

func CompareKeys[K comparable, V1 any, V2 any](
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
func WriteTextFile(dir string, fileName string, content string) {
	err := os.WriteFile(filepath.Join(dir, fileName), []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("File ", fileName, " written successfully.")
}

func ReadTextFile(dir string, fileName string) string {
	content, err := os.ReadFile(filepath.Join(dir, fileName))
	if err != nil {
		log.Println("File reading error ", err)
		return ""
	}

	return string(content)
}

// TODO : test to make sure that this works once I get home
func LANDiscovery() {
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
