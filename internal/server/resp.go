package server

import (
	"bytes"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var LoadTest = false

var rt = newRouter()

func init() {
	rt.add("GET", "/", hRoot)
	rt.add("GET", "/echo", hEcho)

	rt.add("POST", "/data", chain(hDataPost, authMw))
	rt.add("GET", "/data", chain(hDataGetAll, authMw))
	rt.add("GET", "/data/:id", chain(hDataGetOne, authMw))
	rt.add("PUT", "/data/:id", chain(hDataPut, authMw))
	rt.add("DELETE", "/data/:id", chain(hDataDelete, authMw))
	rt.add("PATCH", "/data/:id", chain(hDataPatch, authMw))

	rt.add("GET", "/static/:file", hStatic)
}

func wr(conn net.Conn, req *Req) bool {
	start := time.Now()

	if req.Method == "OPTIONS" {
		writeResp(conn, 200, []byte(""), "text/plain", false)
		if !LoadTest {
			log.Println(req.Method, req.Path, 200, time.Since(start).Milliseconds())
		}
		return false
	}

	keep := shouldKeepAlive(req)
	method := req.Method

	if req.Method == "GET" && req.Path == "/" {
		body := []byte("hello")
		writeResp(conn, 200, body, "text/plain", keep)
		if !LoadTest {
			log.Println(req.Method, req.Path, 200, time.Since(start).Milliseconds())
		}
		return keep
	}

	h := rt.match(method, cleanPath(req.Path))
	if h == nil {
		writeResp(conn, 404, []byte("not found"), "text/plain", false)
		if !LoadTest {
			log.Println(req.Method, req.Path, 404, time.Since(start).Milliseconds())
		}
		return false
	}

	code, body, ctype := h(req)

	if req.Method == "HEAD" {
		body = []byte{}
	}

	writeResp(conn, code, body, ctype, keep)
	if !LoadTest {
		log.Println(req.Method, req.Path, code, time.Since(start).Milliseconds())
	}

	return keep
}

func shouldKeepAlive(req *Req) bool {
	connVal := strings.ToLower(req.Hdr["Connection"])
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
	var buf bytes.Buffer
	buf.Grow(len(body) + 256)
	buf.WriteString("HTTP/1.1 ")
	buf.WriteString(strconv.Itoa(code))
	buf.WriteByte(' ')
	buf.WriteString(statusText(code))
	buf.WriteString("\r\n")
	buf.WriteString("Content-Type: ")
	buf.WriteString(ctype)
	buf.WriteString("\r\n")
	buf.WriteString("Content-Length: ")
	buf.WriteString(strconv.Itoa(len(body)))
	buf.WriteString("\r\n")
	buf.WriteString("Date: ")
	buf.WriteString(time.Now().UTC().Format(time.RFC1123))
	buf.WriteString("\r\n")

	if keep {
		buf.WriteString("Connection: keep-alive\r\n")
	} else {
		buf.WriteString("Connection: close\r\n")
	}

	buf.WriteString("Access-Control-Allow-Origin: *\r\n")
	buf.WriteString("Access-Control-Allow-Methods: GET, POST, PUT, DELETE, PATCH, OPTIONS\r\n")
	buf.WriteString("Access-Control-Allow-Headers: Content-Type, Authorization\r\n")

	if code == 401 {
		buf.WriteString(`WWW-Authenticate: Basic realm="Secure Area"` + "\r\n")
	}

	buf.WriteString("\r\n")

	buf.Write(body)

	_, _ = conn.Write(buf.Bytes())
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
