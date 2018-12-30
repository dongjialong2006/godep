package main

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
	pkgs   map[string]bool
	nodes  []*Node
	names  map[string]bool
	update bool
}
