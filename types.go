package main

import (
	"os"
)

type Node struct {
	name    string
	repo    string
	version string
}

func NewNode(name string, repo string, version string) *Node {
	return &Node{
		name:    name,
		repo:    repo,
		version: version,
	}
}

type Packages struct {
	nodes  []*Node
	update bool
}

func NewPackages(update bool) (*Packages, error) {
	p := &Packages{
		nodes:  make([]*Node, 0),
		update: update,
	}

	if !update {
		if err := os.RemoveAll("./vendor"); nil != err {
			return nil, err
		}
	}

	return p, nil
}
