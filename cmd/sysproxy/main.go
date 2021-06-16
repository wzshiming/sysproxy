package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/wzshiming/sysproxy"
)

var name = filepath.Base(os.Args[0])

func init() {
	if len(os.Args) != 2 {
		os.Stderr.WriteString(fmt.Sprintf(`System proxy settings.
	
	Usage: %s <proxy address>
	      %s 127.0.0.1:8080
	`, name, name))
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) == 1 {
		err := sysproxy.OffHTTP()
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		err := sysproxy.OnHTTP(os.Args[1])
		if err != nil {
			log.Println(err)
			return
		}
	}
}
