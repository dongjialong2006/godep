package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	Branch    = ""
	Version   = ""
	BuildTime = ""
	BuildType = "godep"
)

func PrintVersionWithTime() string {
	if "" == BuildTime {
		BuildTime = time.Now().Format("2006-01-02 15:04:05")
	}

	if "" == Branch {
		Branch = "master"
	}

	if "" == Version {
		Version = "1.0.0"
	}

	return fmt.Sprintf("%s-%v-%v, build time:%v.", BuildType, Branch, Version, BuildTime)
}

func main() {
	var update bool = false
	var version bool = false
	set := flag.NewFlagSet("godep", flag.ContinueOnError)
	set.BoolVar(&update, "update", false, "update packages accord to glide.yaml file")
	set.BoolVar(&version, "version", false, "godep version")
	set.BoolVar(&version, "v", false, "godep version")
	if err := set.Parse(os.Args[1:]); nil != err {
		fmt.Println(err)
		return
	}

	if len(os.Args) > 1 {
		if "-h" == os.Args[1] || "-help" == os.Args[1] {
			return
		}

		os.Args[1] = strings.Trim(os.Args[1], " ")
		if "update" == os.Args[1] || "up" == os.Args[1] || "-update" == os.Args[1] || "-up" == os.Args[1] {
			update = true
		} else if "version" == os.Args[1] || "-version" == os.Args[1] || "v" == os.Args[1] || "-v" == os.Args[1] {
			version = true
		} else {
			fmt.Println(fmt.Sprintf("unknown command:%s.", strings.Join(os.Args, " ")))
			return
		}
	}

	if version {
		fmt.Println(PrintVersionWithTime())
		return
	}

	pkg, err := NewPackages(update)
	if nil != err {
		fmt.Println(err)
		return
	}

	if err = pkg.DownloadPkgs(); nil != err {
		fmt.Println(err)
	}

	return
}
