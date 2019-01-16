package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"strings"
)

func NewPackages(update bool, pkgs string) (*Packages, error) {
	p := &Packages{
		pkgs:   make(map[string]bool),
		nodes:  make([]*Node, 0),
		names:  make(map[string]bool),
		update: update,
	}

	if !update {
		os.RemoveAll("./vendor")
	}

	pkgs = strings.Trim(pkgs, " \n")
	if "" != pkgs {
		tmp := strings.Split(pkgs, " ")
		if 1 == len(tmp) {
			tmp = strings.Split(pkgs, ",")
		}

		for _, name := range tmp {
			name = strings.Trim(name, " \n")
			if "" == name {
				continue
			}
			p.pkgs[name] = true
		}
	}

	return p, nil
}

func (p *Packages) String() {
	for _, node := range p.nodes {
		key := node.name

		if "" != node.version {
			key = fmt.Sprintf("%s, version:%s", key, node.version)
		}

		if "" != node.repo {
			key = fmt.Sprintf("%s, repo:%s", key, node.repo)
		}
		fmt.Println(key)
	}
}

func (p *Packages) Init() error {
	name := FindYamlFile()
	if "" == name {
		return fmt.Errorf("directory:./ do not found yaml file.")
	}

	file, err := os.Open(name)
	if nil != err {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		data, _, err := reader.ReadLine()
		if io.EOF == err {
			break
		}

		if err != nil {
			return err
		}

		if 0 == len(data) {
			continue
		}

		p.init(name, string(data))
	}

	return nil
}

func (p *Packages) DownloadPkgs() error {
	var wg sync.WaitGroup
	var t1 = time.Now()
	for _, node := range p.nodes {
		if "" == node.name {
			continue
		}

		path := fmt.Sprintf("./vendor/%s", node.name)
		if len(p.pkgs) > 0 {
			if !p.updatePkg(node) {
				continue
			}
			os.RemoveAll(path)
			fmt.Println(fmt.Sprintf("package:%s is removed, spend time:%v.", path, time.Now().Sub(t1)))
		}

		if err := CreatePath(path); nil != err {
			return err
		}

		if ok := IsExist(path); ok && p.update {
			continue
		}

		pos := strings.LastIndex(path, "/")
		path = path[:pos]
		if "" == path || " " == path {
			continue
		}

		diff := Diff(node)

		AlterVersion(node)

		wg.Add(1)
		go p.handle(path, diff, node, &wg)
	}
	wg.Wait()

	fmt.Println(fmt.Sprintf("all packages download spend time:%v.", time.Now().Sub(t1)))

	return nil
}

func (p *Packages) updatePkg(node *Node) bool {
	for pkg, _ := range p.pkgs {
		if node.name == pkg || strings.HasSuffix(node.name, pkg) {
			return true
		}
	}

	return false
}

func (p *Packages) init(name string, value string) error {
	value = strings.Trim(value, " \n\t")
	if "" == value {
		return fmt.Errorf("read %s file is empty.", name)
	}

	if !strings.HasPrefix(value, "- package:") && !strings.HasPrefix(value, "repo:") && !strings.HasPrefix(value, "version:") {
		return nil
	}

	value = strings.Replace(value, " ", "", -1)
	values := strings.SplitN(value, ":", 2)

	switch values[0] {
	case "repo":
		p.nodes[len(p.nodes)-1].repo = values[1]
	case "version":
		p.nodes[len(p.nodes)-1].version = strings.Replace(values[1], "~", "v", 1)
	case "-package":
		_, ok := p.names[values[1]]
		if ok {
			return nil
		}
		p.names[values[1]] = true
		p.nodes = append(p.nodes, &Node{
			name: values[1],
		})
	}

	return nil
}

func (p *Packages) handle(path string, diff bool, node *Node, wg *sync.WaitGroup) {
	defer wg.Done()

	var err error = nil
	var t2 = time.Now()
	for i := 0; i < 5; i++ {
		if err = p.exec(path, node); nil != err {
			if !diff {
				switch i {
				case 0:
					num := len(node.version)
					if num > 2 {
						node.version = node.version[:num-2]
					}
				case 1:
					node.version = strings.Replace(node.version, "v", "", -1)
				case 2:
					node.version = ""
				}
			}
			continue
		}

		if err1 := Rename(node); nil != err1 {
			fmt.Println(fmt.Sprintf("package:%s rename err:%v, path:%s.", node.name, err1, path))
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

func (p *Packages) timeout(ch chan struct{}, name string, cmd *exec.Cmd) {
	for i := 0; i < 600; i++ {
		select {
		case <-ch:
			return
		default:
			time.Sleep(time.Second)
		}
	}

	if nil != cmd.Process {
		kill(cmd.Process.Pid)
	}

	fmt.Println(fmt.Sprintf("package:%s download timeout, please checkout net or authentication.", name))
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
	if nil == handle {
		return fmt.Errorf("exec command handle is nil.")
	}

	stdout, err := handle.StdoutPipe()
	if err != nil {
		return err
	}

	go PipeLine(stdout, handle)

	ch := make(chan struct{})
	defer close(ch)
	go p.timeout(ch, node.name, handle)

	return handle.Run()
}
