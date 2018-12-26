package main

import (
	"fmt"
	"sync"
	"time"

	"strings"
)

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
	name := p.findYamlFile()
	if "" == name {
		return fmt.Errorf("directory:./ do not found yaml file.")
	}

	data, err := p.readFile(name)
	if nil != err {
		return err
	}

	value := string(data)
	if "" == value {
		return fmt.Errorf("read %s file is empty.", name)
	}

	values := strings.Split(value, "import:")
	if 2 != len(values) {
		return fmt.Errorf("invalide %s file.", name)
	}

	value = values[1]
	pos := strings.Index(value, "ignore:")
	if -1 != pos {
		pos1 := strings.Index(value[pos:], "testImport:")
		if -1 == pos1 {
			value = value[:pos]
		} else {
			value = strings.Replace(value, value[pos:pos1], "", -1)
		}
	}

	value = strings.Replace(value, "testImport:", "", -1)
	values = strings.Split(value, "- package:")

	var repo, version string
	for _, value = range values {
		value = strings.Trim(value, " \n")
		tmp := strings.Split(value, "\n")

		if len(tmp) > 0 {
			tmp[0] = strings.Trim(tmp[0], " \n")
		}

		if "" == tmp[0] {
			continue
		}

		_, ok := p.names[tmp[0]]
		if ok {
			continue
		}
		p.names[tmp[0]] = true

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
				value = strings.Trim(value, " \n")
				if strings.HasPrefix(value, "repo:") {
					repo = p.checkPrefix(value, "repo:")
				}

				if strings.HasPrefix(value, "version:") {
					version = p.checkPrefix(value, "version:")
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

func (p *Packages) DownloadPkgs() error {
	var wg sync.WaitGroup
	var t1 = time.Now()
	for _, node := range p.nodes {
		if "" == node.name {
			continue
		}

		path := fmt.Sprintf("./vendor/%s", node.name)
		exist, err := p.newFile(path)
		if nil != err {
			return err
		}

		if exist && p.update {
			continue
		}

		pos := strings.LastIndex(path, "/")
		path = path[:pos]
		if "" == path || " " == path {
			continue
		}

		diff := p.diff(node)

		p.alterVersion(node)

		wg.Add(1)
		go p.handle(path, diff, node, &wg)
	}
	wg.Wait()

	fmt.Println(fmt.Sprintf("all packages download spend time:%v.", time.Now().Sub(t1)))

	return nil
}
