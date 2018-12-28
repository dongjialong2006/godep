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

func printVersionWithTime() string {
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

func parse() (bool, bool, error) {
	var err error = nil
	var update bool = false
	var version bool = false
	set := flag.NewFlagSet("godep", flag.ContinueOnError)
	set.BoolVar(&update, "update", false, "update packages accord to glide.yaml file")
	set.BoolVar(&version, "version", false, "godep version")
	if err = set.Parse(os.Args[1:]); nil != err {
		return update, version, err
	}

	if len(os.Args) > 1 {
		if "update" == os.Args[1] || "up" == os.Args[1] {
			update = true
		} else if "version" == os.Args[1] || "v" == os.Args[1] {
			version = true
		} else {
			if !update && !version {
				err = fmt.Errorf("unknown command:%s.", strings.Join(os.Args, " "))
			}
		}
	}

	return update, version, err
}

func main() {
	update, version, err := parse()
	if nil != err {
		fmt.Println(err)
		return
	}

	if version {
		fmt.Println(printVersionWithTime())
		return
	}

	pkg, err := NewPackages(update)
	if nil != err {
		fmt.Println(err)
		return
	}

	if err = pkg.Init(); nil != err {
		fmt.Println(err)
		return
	}

	// pkg.String()

	// fmt.Println("\n\n\n\n\n")

	if err = pkg.DownloadPkgs(); nil != err {
		fmt.Println(err)
	}

	return
}
