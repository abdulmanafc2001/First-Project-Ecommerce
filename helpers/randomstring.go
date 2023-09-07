package helper

import (
	"math/rand"
	"strings"
	"time"
)

func RandomStringGenerator() string {
	lower := "abcdefghijklmnopqrstuvwxyz"
	upper := strings.ToUpper(lower)

	all := lower + upper
	length := 8

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	str := ""
	for i := 0; i < length; i++ {
		str += string(all[random.Intn(len(all))])
	}
	return str
}
