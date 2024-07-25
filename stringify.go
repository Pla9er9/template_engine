package main

import "encoding/json"

func Stringify(v any) string {
	if v == nil {
		return "nil"
	}

	switch v.(type) {
	case string:
		return v.(string)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		return string(b)
	}
}
