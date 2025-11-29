package server

import "strconv"

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

func hDataPost(req *Req) (int, []byte, string) {
	var v interface{}
	err := parseJSON(req.Body, &v)
	if err != nil {
		return 400, []byte("bad json"), "text/plain"
	}
	it := st.add(v)
	return 200, toJSON(it), "application/json"
}

func hDataGetAll(req *Req) (int, []byte, string) {
	all := st.all()
	return 200, toJSON(all), "application/json"
}

func hDataGetOne(req *Req) (int, []byte, string) {
	idStr := req.Path[len("/data/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 400, []byte("bad id"), "text/plain"
	}
	it, ok := st.get(id)
	if !ok {
		return 404, []byte("not found"), "text/plain"
	}
	return 200, toJSON(it), "application/json"
}
