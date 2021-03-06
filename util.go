package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
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
			if IsExist(source) {
				return os.Rename(source, fmt.Sprintf("./vendor/%s", node.name))
			}
		}
	}

	return nil
}

func CreatePath(path string) error {
	if "" == path {
		return fmt.Errorf("path is empty.")
	}

	pos := strings.LastIndex(path, "/")
	if -1 == pos {
		return fmt.Errorf("path:%s is invalide.", path)
	}
	path = path[:pos]
	_, err := os.Stat(path)
	if nil == err {
		return nil
	}

	if os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	return nil
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if nil != err {
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(fmt.Sprintf("path:%s, err:%v.", path, err))
		return false
	}

	files, err := ioutil.ReadDir(path + "/")
	if nil != err {
		fmt.Println(fmt.Sprintf("path:%s, err:%v.", path, err))
		return false
	}

	if len(files) == 0 {
		return false
	}

	return true
}

func PipeLine(r io.Reader, cmd *exec.Cmd) {
	br := bufio.NewReader(r)
	for {
		data, _, err := br.ReadLine()
		if err != nil {
			return
		}

		value := string(data)
		if strings.HasPrefix(value, "package:") {
			fmt.Println(value)
		}

		if strings.Contains(value, "yes/no") {
			if nil != cmd.Process {
				kill(cmd.Process.Pid)
			}
		}
	}
}

func kill(pid int) {
	process, err := os.FindProcess(pid)
	if nil != err {
		return
	}

	process.Signal(syscall.SIGKILL)
}
