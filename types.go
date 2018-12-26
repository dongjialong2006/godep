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
	names  map[string]bool
	update bool
}

func NewPackages(update bool) (*Packages, error) {
	p := &Packages{
		nodes:  make([]*Node, 0),
		names:  make(map[string]bool),
		update: update,
	}

	if !update {
		os.RemoveAll("./vendor")
	}

	return p, nil
}
