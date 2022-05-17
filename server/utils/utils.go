package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func Random(hash string) string {
	strings := []string{}
	for i := range hash {
		strings = append(strings, string(hash[i]))
	}

	rand.Seed(time.Now().UnixNano())

	for i := len(strings) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		strings[i], strings[num] = strings[num], strings[i]
	}

	str := ""
	for i := 0; i < len(strings); i++ {
		str += strings[i]
	}
	return str
}
