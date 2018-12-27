package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func (p *Packages) findYamlFile() string {
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

func (p *Packages) checkPrefix(value string, key string) string {
	if strings.HasPrefix(value, key) {
		value = strings.Replace(value, key, "", -1)
		return strings.Trim(value, " ")
	}

	return value
}

func (p *Packages) alterVersion(node *Node) {
	if "github.com/rifflock/lfshook" == node.name {
		node.version = strings.Trim(node.version, ".0")
	}
}

func (p *Packages) handle(path string, diff bool, node *Node, wg *sync.WaitGroup) {
	defer wg.Done()

	var err error = nil
	var t2 = time.Now()
	for i := 0; i < 4; i++ {
		if err = p.exec(path, node); nil != err {
			if !diff {
				if i > 1 {
					node.version = strings.Replace(node.version, "v", "", -1)
				}
			}
			continue
		}

		if nil != err {
			break
		}

		if err1 := p.rename(node); nil != err1 {
			fmt.Println(fmt.Sprintf("package:%s rename err:%v.", node.name, err1))
		}
		break
	}
	if nil == err {
		fmt.Println(fmt.Sprintf("package:%s is download, spend time:%v.", node.name, time.Now().Sub(t2)))
	} else {
		fmt.Println(fmt.Sprintf("package:%s download err:%v.", node.name, err))
	}

	return
}

func (p *Packages) diff(node *Node) bool {
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

func (p *Packages) rename(node *Node) error {
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
			exist, _ := p.checkFileExist(source)
			if exist {
				return os.Rename(source, fmt.Sprintf("./vendor/%s", node.name))
			}
		}
	}

	return nil
}

func (p *Packages) newFile(path string) (bool, error) {
	exist, err := p.checkFileExist(path)
	if nil != err {
		return exist, err
	}

	if !exist {
		if err = p.newPath(path); nil != err {
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

func (p *Packages) checkFileExist(path string) (bool, error) {
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

func (p *Packages) exec(path string, node *Node) error {
	cmd := fmt.Sprintf("cd %s;git clone", path)
	if "" != node.version {
		cmd += fmt.Sprintf(" -b %s", node.version)
	}

	if "" != node.repo {
		cmd += fmt.Sprintf(" %s", node.repo)
	} else {
		cmd += fmt.Sprintf(" git://%s", node.name)
	}

	// fmt.Println(cmd)

	handle := exec.Command("/bin/bash", "-c", cmd)

	var out bytes.Buffer
	handle.Stdout = &out

	return handle.Run()
}

func (p *Packages) newPath(path string) error {
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

func (p *Packages) readFile(name string) ([]byte, error) {
	if "" == name {
		return nil, fmt.Errorf("path is empty.")
	}

	_, err := os.Stat(name)
	if nil != err {
		return nil, err
	}

	return ioutil.ReadFile(name)
}
