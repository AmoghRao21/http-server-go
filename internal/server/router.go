package server

import "strings"

type Handler func(*Req) (int, []byte, string)

type Router struct {
	routes map[string]map[string]Handler
}

func newRouter() *Router {
	return &Router{
		routes: map[string]map[string]Handler{},
	}
}

func (rt *Router) add(method, path string, h Handler) {
	method = strings.ToUpper(method)
	if rt.routes[method] == nil {
		rt.routes[method] = map[string]Handler{}
	}
	rt.routes[method][path] = h
}

func (rt *Router) match(method, path string) Handler {
	method = strings.ToUpper(method)
	if m, ok := rt.routes[method]; ok {
		if h, ok := m[path]; ok {
			return h
		}
		for p, h := range m {
			if strings.HasSuffix(p, "/:id") {
				base := p[:len(p)-4]
				if strings.HasPrefix(path, base) {
					return h
				}
			}
			if strings.HasSuffix(p, "/:file") {
				base := p[:len(p)-6]
				if strings.HasPrefix(path, base) {
					return h
				}
			}
			if strings.HasPrefix(path, "/static/") && p == "/static/:file" {
				return h
			}
		}
	}
	return nil
}
