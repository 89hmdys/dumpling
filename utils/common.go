package utils

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"net"
	"sort"
	"strings"
)

const FORMAT_DATETIME_SHORT = "20060102150405"
const FORMAT_DATETIME = "2006-01-02 15:04:05"

func ParaFilter(para map[string]string) {
	for key, value := range para {
		if value == "" || key == "sign" || key == "sign_type" {
			delete(para, key)
		}
	}
}

func Join(para map[string]string) string {
	toSort := []string{}
	for key, value := range para {
		toSort = append(toSort, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(toSort, "&")
}

func SortAndJoin(para map[string]string) string {
	toSort := []string{}
	for key, value := range para {
		toSort = append(toSort, fmt.Sprintf("%s=%s", key, value))
	}
	sort.Strings(toSort)
	return strings.Join(toSort, "&")
}

func UUID() string {
	return strings.Replace(uuid.NewRandom().String(), "-", "", -1)
}

func LocalIP() string {
	info, _ := net.InterfaceAddrs()
	for _, addr := range info {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			return ""
		}
		if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}
	return ""
}
