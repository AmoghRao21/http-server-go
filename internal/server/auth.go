package server

import (
	"encoding/base64"
	"strings"
)

var authUser = "admin"
var authPass = "123123123"

func checkAuth(req *Req) bool {
	h := req.Hdr["authorization"]
	if h == "" {
		return false
	}
	if !strings.HasPrefix(h, "Basic ") {
		return false
	}
	raw := strings.TrimPrefix(h, "Basic ")
	dec, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return false
	}
	parts := strings.SplitN(string(dec), ":", 2)
	if len(parts) != 2 {
		return false
	}
	return parts[0] == authUser && parts[1] == authPass
}
