package utils

import (
	"math/rand"
	"strings"
)

func GetRandomHost(hosts string) string {
	host := strings.Split(hosts, ",")
	return host[rand.Intn(len(host))]
}
