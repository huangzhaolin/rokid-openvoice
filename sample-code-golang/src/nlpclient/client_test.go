package main

import "testing"

func TestNlp(t *testing.T) {
	conn, err := do_conn("apigw.open.rokid.com:443", true, true)
	if err != nil {
		return
	}
	defer conn.Close()

	call_nlp(conn, "zh", "我要听张学友的歌", "com.rokid.music1", true)
}

func BenchmarkNlp(b *testing.B) {
	conn, err := do_conn("apigw.open.rokid.com:443", true, true)
	if err != nil {
		return
	}
	defer conn.Close()

	for i := 0; i < b.N; i++ {
		call_nlp(conn, "zh", "我要听张学友的歌", "com.rokid.music1", true)
	}
}
