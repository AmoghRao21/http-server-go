package server

import "strings"

type Handler func(*Req) (int, []byte, string)

type fastNode struct {
	part     string
	isWild   bool
	handler  Handler
	children []*fastNode
}

type Router struct {
	roots map[string]*fastNode
}

func newRouter() *Router {
	return &Router{
		roots: make(map[string]*fastNode),
	}
}

func (r *Router) add(method, path string, h Handler) {
	method = strings.ToUpper(method)
	parts := strings.Split(path, "/")[1:]

	root, ok := r.roots[method]
	if !ok {
		root = &fastNode{}
		r.roots[method] = root
	}
	insertFast(root, parts, h)
}

func insertFast(n *fastNode, parts []string, h Handler) {
	if len(parts) == 0 {
		n.handler = h
		return
	}
	part := parts[0]
	isWild := len(part) > 0 && part[0] == ':'

	for _, c := range n.children {
		if c.part == part || c.isWild {
			insertFast(c, parts[1:], h)
			return
		}
	}

	child := &fastNode{
		part:   part,
		isWild: isWild,
	}
	n.children = append(n.children, child)
	insertFast(child, parts[1:], h)
}

func (r *Router) match(method, path string) Handler {
	method = strings.ToUpper(method)
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	parts := strings.Split(path, "/")[1:]

	n := searchFast(root, parts)
	if n == nil {
		return nil
	}
	return n.handler
}

func searchFast(n *fastNode, parts []string) *fastNode {
	if len(parts) == 0 {
		return n
	}
	part := parts[0]

	for _, c := range n.children {
		if c.part == part || c.isWild {
			out := searchFast(c, parts[1:])
			if out != nil {
				return out
			}
		}
	}
	return nil
}
