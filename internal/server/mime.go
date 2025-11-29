package server

import "strings"

func mime(path string) string {
	l := strings.ToLower(path)
	switch {
	case strings.HasSuffix(l, ".html"):
		return "text/html"
	case strings.HasSuffix(l, ".css"):
		return "text/css"
	case strings.HasSuffix(l, ".js"):
		return "application/javascript"
	case strings.HasSuffix(l, ".png"):
		return "image/png"
	case strings.HasSuffix(l, ".jpg"), strings.HasSuffix(l, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(l, ".gif"):
		return "image/gif"
	case strings.HasSuffix(l, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(l, ".json"):
		return "application/json"
	default:
		return "application/octet-stream"
	}
}
