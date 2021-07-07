package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/wzshiming/sysproxy"
)

var name = filepath.Base(os.Args[0])

func init() {
	if len(os.Args) < 2 {
		printDefaults()
	}
}

func printDefaults() {
	os.Stderr.WriteString(fmt.Sprintf(`System proxy settings.
Usage: 
	%s <scheme> <option> <proxy address>
	%s http on 127.0.0.1:8080
	%s http off
	%s http get
	%s https on 127.0.0.1:8080
	%s https off
	%s https get
	%s no_proxy on 127.0.0.1,[::]
	%s no_proxy off
	%s no_proxy get
	%s pac on http://127.0.0.1:8080
	%s pac off
	%s pac get
`, name, name, name, name, name, name, name, name, name, name, name, name, name))
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	args := os.Args[1:]
	switch args[0] {
	default:
		log.Printf("unsupport %q", args[0])
		printDefaults()
	case "http":
		handleHTTP(args[1:])
	case "https":
		handleHTTPS(args[1:])
	case "pac":
		handlePAC(args[1:])
	case "no_proxy":
		handleNoProxy(args[1:])
	}
}

func handleHTTP(args []string) {
	if len(args) == 0 {
		printDefaults()
		return
	}

	switch args[0] {
	default:
		log.Printf("unsupport %q", args[0])
		printDefaults()
	case "on":
		if len(args) == 1 {
			printDefaults()
			return
		}
		err := sysproxy.OnHTTP(args[1])
		if err != nil {
			log.Println(err)
			printDefaults()
			return
		}
	case "off":
		err := sysproxy.OffHTTP()
		if err != nil {
			log.Println(err)
			printDefaults()
			return
		}
	case "get":
		out, err := sysproxy.GetHTTP()
		if err != nil {
			log.Println(err)
			printDefaults()
			return
		}
		fmt.Println(out)
	}
}

func handleHTTPS(args []string) {
	if len(args) == 0 {
		printDefaults()
		return
	}

	switch args[0] {
	default:
		log.Printf("unsupport %q", args[0])
		printDefaults()
	case "on":
		if len(args) == 1 {
			printDefaults()
			return
		}
		err := sysproxy.OnHTTPS(args[1])
		if err != nil {
			log.Println(err)
			printDefaults()
			return
		}
	case "off":
		err := sysproxy.OffHTTPS()
		if err != nil {
			log.Println(err)
			printDefaults()
			return
		}
	case "get":
		out, err := sysproxy.GetHTTPS()
		if err != nil {
			log.Println(err)
			printDefaults()
			return
		}
		fmt.Println(out)
	}
}

func handleNoProxy(args []string) {
	if len(args) == 0 {
		printDefaults()
		return
	}

	switch args[0] {
	default:
		log.Printf("unsupport %q", args[0])
		printDefaults()
	case "on":
		if len(args) == 1 {
			printDefaults()
			return
		}
		err := sysproxy.OnNoProxy(strings.Split(args[1], ","))
		if err != nil {
			log.Println(err)
			printDefaults()
			return
		}
	case "off":
		err := sysproxy.OffNoProxy()
		if err != nil {
			log.Println(err)
			printDefaults()
			return
		}
	case "get":
		out, err := sysproxy.GetNoProxy()
		if err != nil {
			log.Println(err)
			printDefaults()
			return
		}
		fmt.Println(strings.Join(out, ","))
	}
}

func handlePAC(args []string) {
	if len(args) == 0 {
		printDefaults()
		return
	}

	switch args[0] {
	default:
		log.Printf("unsupport %q", args[0])
		printDefaults()
	case "on":
		if len(args) == 1 {
			printDefaults()
			return
		}
		err := sysproxy.OnPAC(args[1])
		if err != nil {
			log.Println(err)
			printDefaults()
			return
		}
	case "off":
		err := sysproxy.OffPAC()
		if err != nil {
			log.Println(err)
			printDefaults()
			return
		}
	case "get":
		out, err := sysproxy.GetPAC()
		if err != nil {
			log.Println(err)
			printDefaults()
			return
		}
		fmt.Println(out)
	}
}
