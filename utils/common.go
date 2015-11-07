package utils

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sort"
	"strings"
)

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

func Post(url string, header map[string]string, content string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, strings.NewReader(content))
	if err != nil {
		return "", err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}
