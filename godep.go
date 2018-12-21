package main

import (
	"fmt"
)

func main() {
	pkg, err := NewPackages()
	if nil != err {
		fmt.Println(err)
		return
	}

	if err = pkg.DownloadPkgs(); nil != err {
		fmt.Println(err)
	}

	return
}
