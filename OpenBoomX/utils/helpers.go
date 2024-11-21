package utils

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"sort"
)

func Must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}

func IsValidHex(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil
}

func SortedKeysByValue(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return m[keys[i]] < m[keys[j]]
	})
	return keys
}

func NrgbaToHex(c color.NRGBA) string {
	return fmt.Sprintf("%02x%02x%02x", c.R, c.G, c.B)
}
