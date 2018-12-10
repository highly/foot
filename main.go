package main

import (
	"fmt"
	"github.com/highly/foot/config"
	"os"
	"path/filepath"
)

func main() {
	config.New().ConfigPath(filepath.Join(BaseDir(), "config")).Load("cnf")
	// effective immediately
	//config.New().ConfigPath(filepath.Join(BaseDir(), "config")).Load("cnf").Watching()
	fmt.Println(config.GetString("Users.default.Shawn"))

	//
}

func BaseDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}
