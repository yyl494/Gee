package gee

import (
	"strings"
)

type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
}

// find the first child match the pattern
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// find all children match the pattern
func (n *node) matchChildren(part string) []*node {
	ret := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			ret = append(ret, child)
		}
	}
	return ret
}

func (n *node) insert(pattern string, parts []string, height int) {
	if height == len(parts) {
		n.pattern = pattern
		return
	}

	child := n.matchChild(parts[height])
	if child == nil {
		child = &node{part: parts[height], isWild: parts[height][0] == ':' || parts[height][0] == '*'}
		n.children = append(n.children, child)
	}

	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if height == len(parts) || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		} else {
			return n
		}
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
