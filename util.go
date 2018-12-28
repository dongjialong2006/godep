package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func FindYamlFile() string {
	files, err := ioutil.ReadDir("./")
	if nil != err {
		fmt.Println(err)
		return ""
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(file.Name(), ".yaml") {
			return file.Name()
		}
	}

	return ""
}

func AlterVersion(node *Node) {
	if "github.com/rifflock/lfshook" == node.name {
		node.version = strings.Trim(node.version, ".0")
	}
}

func Diff(node *Node) bool {
	if "" != node.repo {
		pos := strings.LastIndex(node.repo, "/")
		value := node.repo[pos+1:]
		pos = strings.Index(value, ".")
		if -1 != pos {
			value = value[:pos]
		}

		if !strings.HasSuffix(node.name, value) {
			return true
		}
	}

	return false
}

func Rename(node *Node) error {
	if "" != node.repo {
		pos := strings.LastIndex(node.repo, "/")
		value := node.repo[pos+1:]
		pos = strings.Index(value, ".")
		if -1 != pos {
			value = value[:pos]
		}

		if !strings.HasSuffix(node.name, value) {
			pos = strings.LastIndex(node.name, "/")
			source := fmt.Sprintf("./vendor/%s%s", node.name[:pos+1], value)
			exist, _ := CheckFileExist(source)
			if exist {
				return os.Rename(source, fmt.Sprintf("./vendor/%s", node.name))
			}
		}
	}

	return nil
}

func OpenFile(path string) (bool, error) {
	exist, err := CheckFileExist(path)
	if nil != err {
		return exist, err
	}

	if !exist {
		if err = CreatePath(path); nil != err {
			return exist, err
		}
	}

	files, err := ioutil.ReadDir(path + "/")
	if nil != err {
		return false, nil
	}

	if len(files) == 0 {
		exist = false
	}

	return exist, nil
}

func CheckFileExist(path string) (bool, error) {
	if "" == path {
		return false, fmt.Errorf("path is empty.")
	}

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func CreatePath(path string) error {
	if "" == path {
		return fmt.Errorf("path is empty.")
	}

	pos := strings.LastIndex(path, "/")
	if -1 == pos {
		return nil
	}
	path = path[:pos]
	_, err := os.Stat(path)
	if nil == err {
		return nil
	}

	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
	}

	return err
}
