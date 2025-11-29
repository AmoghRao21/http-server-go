package server

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"time"
)

func wr(conn net.Conn, req *Req) {
	code := 200
	body := []byte("welcome")
	contentType := "text/plain; charset=utf-8"

	if req.Method == "GET" && req.Path == "/" {
		code = 200
		body = []byte("welcome")
	} else {
		code = 404
		body = []byte("not found")
	}

	headers := map[string]string{}
	headers["Content-Type"] = contentType
	headers["Content-Length"] = strconv.Itoa(len(body))
	headers["Date"] = time.Now().UTC().Format(time.RFC1123)
	headers["Connection"] = "close"

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "HTTP/1.1 %d %s\r\n", code, statusText(code))
	for key, val := range headers {
		fmt.Fprintf(&buf, "%s: %s\r\n", key, val)
	}
	buf.WriteString("\r\n")
	buf.Write(body)

	conn.Write(buf.Bytes())
}

func statusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return "Status"
	}
}
