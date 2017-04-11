### 编译
```
git clone git@github.com:Rokid/rokid-openvoice.git
cd rokid-openvoice/sample-code-golang
make
```

### 认证

修改src/auth/credentials.go中认证信息:
```
15     key            = "PLEASE_FIX_IT"
16     device_type_id = "PLEASE_FIX_IT"
17     device_id      = "PLEASE_FIX_IT"
18     secret         = "PLEASE_FIX_IT"
```


### 运行
```
$ ./asrclient -tls -auth -file zhrmghg.pcm
2017/04/06 20:05:06.634021 main.go:84: Read file(9600)
2017/04/06 20:05:06.934236 main.go:84: Read file(9600)
2017/04/06 20:05:07.234433 main.go:84: Read file(9600)
2017/04/06 20:05:07.534698 main.go:84: Read file(9600)
2017/04/06 20:05:07.834861 main.go:84: Read file(9600)
2017/04/06 20:05:08.135399 main.go:84: Read file(9600)
2017/04/06 20:05:08.435646 main.go:84: Read file(9600)
2017/04/06 20:05:08.735873 main.go:84: Read file(9600)
2017/04/06 20:05:09.036131 main.go:84: Read file(3196)
2017/04/06 20:05:09.336334 main.go:90: Read file: EOF
2017/04/06 20:05:09.844997 main.go:66: Got asr: Asr('中华人民共和国')
```

```
$ ./ttsclient -tls -auth
2017/04/06 20:09:43.765871 main.go:70: Got tts: Text(''), Voice(len=1800)
2017/04/06 20:09:43.830424 main.go:70: Got tts: Text(''), Voice(len=1280)
```

```
$ ./nlpclient -tls -auth
2017/04/06 20:09:56.238751 main.go:53: Nlp() = {"cloud":false,"confidence":0,"domain":"ROKID.EXCEPTION.NLP","intent":"NO_NLP","posEnd":0,"posStart":0,"slots":{},"voice":""}
```

```
$ ./speechtclient -tls -auth
2017/04/06 20:10:12.377161 main.go:57: Speecht() = asr(), nlp({"cloud":false,"confidence":0,"domain":"ROKID.EXCEPTION.NLP","intent":"NO_NLP","posEnd":0,"posStart":0,"slots":{},"voice":""}), action({"response":{"action":{"shoudEndSession":true,"type":"NORMAL","version":"2.0.0","voice":{"behaviour":"APPEND","item":{"tts":"未命中语音指令"},"needEventCallback":false}},"domain":"ROKID.EXCEPTION.NLP","resType":"INTENT","shot":"CUT"},"version":"2.0.0"})
```

```
$ ./speechvclient -tls -auth -file zhrmghg.pcm
2017/04/06 20:10:27.462962 main.go:96: Read file(9600)
2017/04/06 20:10:27.763276 main.go:96: Read file(9600)
2017/04/06 20:10:28.063522 main.go:96: Read file(9600)
2017/04/06 20:10:28.363844 main.go:96: Read file(9600)
2017/04/06 20:10:28.664124 main.go:96: Read file(9600)
2017/04/06 20:10:28.964324 main.go:96: Read file(9600)
2017/04/06 20:10:29.264642 main.go:96: Read file(9600)
2017/04/06 20:10:29.564953 main.go:96: Read file(9600)
2017/04/06 20:10:29.929904 main.go:96: Read file(3196)
2017/04/06 20:10:30.230184 main.go:102: Read file: EOF
2017/04/06 20:10:30.753039 main.go:68: Got asr: asr(中华人民共和国), nlp(), action()
2017/04/06 20:10:30.860197 main.go:68: Got asr: asr(中华人民共和国), nlp({"cloud":false,"confidence":0,"domain":"ROKID.EXCEPTION.NLP","intent":"NO_NLP","posEnd":0,"posStart":0,"slots":{},"voice":"中华人民共和国"}), action({"response":{"action":{"shoudEndSession":true,"type":"NORMAL","version":"2.0.0","voice":{"behaviour":"APPEND","item":{"tts":"未命中语音指令"},"needEventCallback":false}},"domain":"ROKID.EXCEPTION.NLP","resType":"INTENT","shot":"CUT"},"version":"2.0.0"})
```

### 测试

```
$ go test -v asrclient
=== RUN   TestAsr
2017/04/06 20:00:11 Read file(9600)
2017/04/06 20:00:11 Read file(9600)
2017/04/06 20:00:11 Read file(9600)
2017/04/06 20:00:12 Read file(9600)
2017/04/06 20:00:12 Read file(9600)
2017/04/06 20:00:12 Read file(9600)
2017/04/06 20:00:13 Read file(9600)
2017/04/06 20:00:13 Read file(9600)
2017/04/06 20:00:13 Read file(3196)
2017/04/06 20:00:14 Read file: EOF
2017/04/06 20:00:14 Got asr: Asr('中华人民共和国')
--- PASS: TestAsr (3.85s)
PASS
ok      asrclient   3.853s
```

```
$ go test -v ttsclient
=== RUN   TestTts
2017/04/06 20:00:25 Got tts: Text(''), Voice(len=2009)
2017/04/06 20:00:25 Got tts: Text(''), Voice(len=1230)
2017/04/06 20:00:25 Got tts: Text(''), Voice(len=1435)
2017/04/06 20:00:25 Got tts: Text(''), Voice(len=2501)
2017/04/06 20:00:25 Got tts: Text(''), Voice(len=738)
2017/04/06 20:00:26 Got tts: Text(''), Voice(len=1845)
2017/04/06 20:00:26 Got tts: Text(''), Voice(len=1722)
2017/04/06 20:00:26 Got tts: Text(''), Voice(len=2009)
2017/04/06 20:00:26 Got tts: Text(''), Voice(len=1558)
2017/04/06 20:00:26 Got tts: Text(''), Voice(len=1804)
2017/04/06 20:00:26 Got tts: Text(''), Voice(len=697)
--- PASS: TestTts (2.78s)
PASS
ok      ttsclient   2.789s
```

```
$ go test -v nlpclient
=== RUN   TestNlp
2017/04/06 20:00:33 Nlp(我要听张学友的歌) = {"cloud":false,"confidence":0,"domain":"ROKID.EXCEPTION.NLP","intent":"NO_NLP","posEnd":0,"posStart":0,"slots":{},"voice":"我要听张学友的歌"}
--- PASS: TestNlp (0.50s)
PASS
ok      nlpclient   0.504s
```

```
$ go test -v speechtclient
=== RUN   TestSpeecht
2017/04/06 20:00:42 Speecht(我要听张学友的歌) = asr(我要听张学友的歌), nlp({"cloud":false,"confidence":0,"domain":"ROKID.EXCEPTION.NLP","intent":"NO_NLP","posEnd":0,"posStart":0,"slots":{},"voice":"我要听张学友的歌"}), action({"response":{"action":{"shoudEndSession":true,"type":"NORMAL","version":"2.0.0","voice":{"behaviour":"APPEND","item":{"tts":"未命中语音指令"},"needEventCallback":false}},"domain":"ROKID.EXCEPTION.NLP","resType":"INTENT","shot":"CUT"},"version":"2.0.0"})
--- PASS: TestSpeecht (0.46s)
PASS
ok      speechtclient   0.466s
```

```
$ go test -v speechvclient
=== RUN   TestSpeechv
2017/04/06 20:00:56 Read file(9600)
2017/04/06 20:00:56 Read file(9600)
2017/04/06 20:00:56 Read file(9600)
2017/04/06 20:00:57 Read file(9600)
2017/04/06 20:00:57 Read file(9600)
2017/04/06 20:00:57 Read file(9600)
2017/04/06 20:00:57 Read file(9600)
2017/04/06 20:00:58 Read file(9600)
2017/04/06 20:00:58 Read file(3196)
2017/04/06 20:00:58 Read file: EOF
2017/04/06 20:00:59 Got asr: asr(中华人民共和国), nlp(), action()
2017/04/06 20:00:59 Got asr: asr(中华人民共和国), nlp({"cloud":false,"confidence":0,"domain":"ROKID.EXCEPTION.NLP","intent":"NO_NLP","posEnd":0,"posStart":0,"slots":{},"voice":"中华人民共和国"}), action({"response":{"action":{"shoudEndSession":true,"type":"NORMAL","version":"2.0.0","voice":{"behaviour":"APPEND","item":{"tts":"未命中语音指令"},"needEventCallback":false}},"domain":"ROKID.EXCEPTION.NLP","resType":"INTENT","shot":"CUT"},"version":"2.0.0"})
--- PASS: TestSpeechv (3.98s)
PASS
ok      speechvclient   3.983s
```
