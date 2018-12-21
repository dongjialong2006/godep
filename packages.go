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

type Node struct {
	name    string
	repo    string
	version string
}

type Packages struct {
	nodes []*Node
}

func NewPackages() (*Packages, error) {
	pkg := &Packages{
		nodes: make([]*Node, 0),
	}

	os.RemoveAll("./vendor")

	if err := pkg.init(); nil != err {
		return nil, err
	}

	return pkg, nil
}

func (p *Packages) init() error {
	data, err := p.readFile("./glide.yaml")
	if nil != err {
		return err
	}

	value := string(data)
	values := strings.Split(value, "import:\n")
	if 2 != len(values) {
		return fmt.Errorf("invalide glide.yaml file.")
	}

	value = values[1]
	pos := strings.Index(value, "ignore:\n")
	if -1 != pos {
		pos1 := strings.Index(value[pos:], "testImport:\n")
		if -1 == pos1 {
			value = value[:pos]
		} else {
			value = strings.Replace(value, value[pos:pos1], "", -1)
		}
	}

	value = strings.Replace(value, "testImport:\n", "", -1)
	values = strings.Split(value, "- package:")

	var repo, version string
	for _, value = range values {
		value = strings.Trim(value, " ")
		tmp := strings.Split(value, "\n")

		if len(tmp) > 0 {
			tmp[0] = strings.Trim(tmp[0], "\n")
			tmp[0] = strings.Trim(tmp[0], " ")
		}

		if "" == tmp[0] {
			continue
		}

		repo = ""
		version = ""

		switch len(tmp) {
		case 0:
			continue
		case 1:
			p.nodes = append(p.nodes, &Node{
				name: tmp[0],
			})
			continue
		default:
			for _, value = range tmp[1:] {
				value = strings.Trim(value, "\n")
				value = strings.Trim(value, " ")
				if strings.HasPrefix(value, "repo:") {
					repo = p.check(value, "repo:")
				}

				if strings.HasPrefix(value, "version:") {
					version = p.check(value, "version:")
					version = strings.Replace(version, "~", "v", 1)
				}
			}
		}

		p.nodes = append(p.nodes, &Node{
			name:    tmp[0],
			repo:    repo,
			version: version,
		})
	}

	return nil
}

func (p *Packages) check(value string, key string) string {
	if strings.HasPrefix(value, key) {
		value = strings.Replace(value, key, "", -1)
		return strings.Trim(value, " ")
	}

	return value
}

func (p *Packages) filter(value string, key string, before bool) string {
	if strings.Contains(value, key) {
		values := strings.Split(value, key)
		if before {
			return values[0]
		}
		return values[1]
	}

	return value
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

func (p *Packages) alterVersion(node *Node) {
	if "github.com/rifflock/lfshook" == node.name {
		node.version = strings.Trim(node.version, ".0")
	}
}

func (p *Packages) DownloadPkgs() error {
	var err error = nil
	var wg sync.WaitGroup

	t1 := time.Now()
	for _, node := range p.nodes {
		if "" == node.name {
			continue
		}

		path := fmt.Sprintf("./vendor/%s", node.name)
		if err = p.newPath(path); nil != err {
			return err
		}

		pos := strings.LastIndex(path, "/")
		path = path[:pos]
		if "" == path || " " == path {
			continue
		}

		p.alterVersion(node)

		wg.Add(1)
		go func(path string, node *Node) {
			var err error = nil
			var t2 = time.Now()
			for i := 0; i < 2; i++ {
				if err = p.exec(path, node); nil != err {
					node.version = strings.Replace(node.version, "v", "", -1)
					continue
				}
				break
			}
			if nil == err {
				fmt.Println(fmt.Sprintf("package:%s is download, spend time:%v.", node.name, time.Now().Sub(t2)))
			} else {
				fmt.Println(err)
			}
			wg.Done()
		}(path, node)
	}
	wg.Wait()
	fmt.Println(fmt.Sprintf("all packages download spend time:%v.", time.Now().Sub(t1)))
	return nil
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
