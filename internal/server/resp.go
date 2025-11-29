package server

import (
	"bytes"
	"fmt"
	"log"
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
	start := time.Now()

	if req.Method == "OPTIONS" {
		writeResp(conn, 200, []byte(""), "text/plain")
		log.Println(req.Method, req.Path, 200, time.Since(start).Milliseconds())
		return
	}

	method := req.Method
	if method == "HEAD" {
		method = "GET"
	}

	h := rt.match(method, cleanPath(req.Path))
	if h == nil {
		writeResp(conn, 404, []byte("not found"), "text/plain")
		log.Println(req.Method, req.Path, 404, time.Since(start).Milliseconds())
		return
	}

	code, body, ctype := h(req)

	if req.Method == "HEAD" {
		body = []byte{}
	}

	writeResp(conn, code, body, ctype)

	log.Println(req.Method, req.Path, code, time.Since(start).Milliseconds())
}

func writeResp(conn net.Conn, code int, body []byte, ctype string) {
	hdr := map[string]string{}
	hdr["Content-Type"] = ctype
	hdr["Content-Length"] = strconv.Itoa(len(body))
	hdr["Date"] = time.Now().UTC().Format(time.RFC1123)
	hdr["Connection"] = "close"
	hdr["Access-Control-Allow-Origin"] = "*"
	hdr["Access-Control-Allow-Methods"] = "GET, POST, PUT, DELETE, PATCH, OPTIONS"
	hdr["Access-Control-Allow-Headers"] = "Content-Type, Authorization"

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
