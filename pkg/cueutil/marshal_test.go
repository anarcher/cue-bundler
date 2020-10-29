package cueutil

import "testing"

func TestMarshal(t *testing.T) {
	type Value struct {
		Title string `json:"title"`
	}
	v := &Value{
		Title: "aaa",
	}

	bs, err := Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs))
}
