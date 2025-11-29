package server

import "encoding/json"

func toJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

func parseJSON(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}
