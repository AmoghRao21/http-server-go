package server

func hRoot(req *Req) (int, []byte, string) {
	return 200, []byte("welcome"), "text/plain"
}

func hEcho(req *Req) (int, []byte, string) {
	msg := req.Query["message"]
	if msg == "" {
		msg = req.Query["msg"]
	}
	return 200, []byte(msg), "text/plain"
}
