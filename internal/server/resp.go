package server

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"time"
)

var rt = newRouter()

func init() {
	rt.add("GET", "/", hRoot)
	rt.add("GET", "/echo", hEcho)
	rt.add("POST", "/data", hDataPost)
	rt.add("GET", "/data", hDataGetAll)
	rt.add("GET", "/data/:id", hDataGetOne)
}

func wr(conn net.Conn, req *Req) {
	h := rt.match(req.Method, cleanPath(req.Path))
	if h == nil {
		writeResp(conn, 404, []byte("not found"), "text/plain")
		return
	}

	code, body, ctype := h(req)
	writeResp(conn, code, body, ctype)
}

func writeResp(conn net.Conn, code int, body []byte, ctype string) {
	hdr := map[string]string{}
	hdr["Content-Type"] = ctype
	hdr["Content-Length"] = strconv.Itoa(len(body))
	hdr["Date"] = time.Now().UTC().Format(time.RFC1123)
	hdr["Connection"] = "close"

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "HTTP/1.1 %d %s\r\n", code, statusText(code))
	for k, v := range hdr {
		fmt.Fprintf(&buf, "%s: %s\r\n", k, v)
	}
	buf.WriteString("\r\n")
	buf.Write(body)

	conn.Write(buf.Bytes())
}

func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		return "/" + p
	}
	return p
}
