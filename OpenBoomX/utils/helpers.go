package utils

import (
	"encoding/hex"
	"sort"
	"strconv"
	"strings"
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

// Str2ba converts MAC address string representation to little-endian byte array
func Str2ba(addr string) [6]byte {
	a := strings.Split(addr, ":")
	var b [6]byte
	for i, tmp := range a {
		u, _ := strconv.ParseUint(tmp, 16, 8)
		b[len(b)-1-i] = byte(u)
	}
	return b
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
