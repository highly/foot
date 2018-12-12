package utils

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
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
