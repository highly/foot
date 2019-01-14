package utils

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func BaseDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}

func Ip() string {
	var ip string
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ip = ipnet.IP.String()
			break
		}
	}
	return ip
}

func IsDir(path string) bool {
	if f, err := os.Stat(path); err == nil {
		return f.Mode().IsDir()
	}
	return false
}

func IsFile(path string) bool {
	if f, err := os.Stat(path); err == nil {
		return f.Mode().IsRegular()
	}
	return false
}

func Md5Hash(in string) string {
	h := md5.New()
	io.WriteString(h, in)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Implode(in interface{}, glue string) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(in)), glue), "[]")
}

func IsJson(body string) bool {
	var temp map[string]interface{}
	return json.Unmarshal([]byte(body), &temp) == nil
}
