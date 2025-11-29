package server

type Middleware func(Handler) Handler

func chain(h Handler, m ...Middleware) Handler {
	if len(m) == 0 {
		return h
	}
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}

func authMw(h Handler) Handler {
	return func(req *Req) (int, []byte, string) {
		if !checkAuth(req) {
			return 401, []byte("unauthorized"), "text/plain"
		}
		return h(req)
	}
}
