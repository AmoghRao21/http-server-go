package server

func hRoot(req *Req) (int, []byte, string) {
	return 200, []byte("welcome"), "text/plain"
}
