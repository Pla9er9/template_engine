package main

import "encoding/json"

func Stringify(v any) string {
	if v == nil {
		return ""
	}
	switch v.(type) {
	case string:
		return v.(string)
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}
