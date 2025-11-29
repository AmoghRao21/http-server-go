package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var rt = newRouter()

func init() {
	rt.add("GET", "/", hRoot)
	rt.add("GET", "/echo", hEcho)
	rt.add("POST", "/data", hDataPost)
	rt.add("GET", "/data", hDataGetAll)
	rt.add("GET", "/data/:id", hDataGetOne)
	rt.add("GET", "/static/:file", hStatic)
}

func wr(conn net.Conn, req *Req) bool {
	start := time.Now()

	if req.Method == "OPTIONS" {
		writeResp(conn, 200, []byte(""), "text/plain", false)
		log.Println(req.Method, req.Path, 200, time.Since(start).Milliseconds())
		return false
	}

	method := req.Method
	if method == "HEAD" {
		method = "GET"
	}

	h := rt.match(method, cleanPath(req.Path))
	if h == nil {
		writeResp(conn, 404, []byte("not found"), "text/plain", false)
		log.Println(req.Method, req.Path, 404, time.Since(start).Milliseconds())
		return false
	}

	code, body, ctype := h(req)

	keep := shouldKeepAlive(req)
	if req.Method == "HEAD" {
		body = []byte{}
	}

	writeResp(conn, code, body, ctype, keep)
	log.Println(req.Method, req.Path, code, time.Since(start).Milliseconds())

	return keep
}

func shouldKeepAlive(req *Req) bool {
	connVal := strings.ToLower(req.Hdr["connection"])
	if connVal == "close" {
		return false
	}
	if req.Ver == "HTTP/1.1" {
		if connVal == "" || connVal == "keep-alive" {
			return true
		}
	}
	if req.Ver == "HTTP/1.0" && connVal == "keep-alive" {
		return true
	}
	return false
}

func writeResp(conn net.Conn, code int, body []byte, ctype string, keep bool) {
	hdr := map[string]string{}
	hdr["Content-Type"] = ctype
	hdr["Content-Length"] = strconv.Itoa(len(body))
	hdr["Date"] = time.Now().UTC().Format(time.RFC1123)
	if keep {
		hdr["Connection"] = "keep-alive"
	} else {
		hdr["Connection"] = "close"
	}
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

func statusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	case 413:
		return "Payload Too Large"
	case 500:
		return "Internal Server Error"
	default:
		return "Status"
	}
}
