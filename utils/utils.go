package utils

import "time"

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
