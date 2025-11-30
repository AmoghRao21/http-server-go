package server

import (
	"encoding/base64"
	"strings"
)

var authUser = "admin"
var authPass = "secret"

func checkAuth(req *Req) bool {
	var h string

	if v, ok := req.Hdr["Authorization"]; ok {
		h = v
	} else if v, ok := req.Hdr["authorization"]; ok {
		h = v
	} else if v, ok := req.Hdr["AUTHORIZATION"]; ok {
		h = v
	}

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
