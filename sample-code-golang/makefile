target := ttsclient asrclient nlpclient speechvclient speechtclient

all: pre $(target)

pre: src/pb/speech.pb.go

src/pb/speech.pb.go: src/pb/speech.proto
	cd src/pb; protoc --go_out=plugins=grpc:. *.proto

ttsclient:
	go build $@

asrclient:
	go build $@

nlpclient:
	go build $@

speechvclient:
	go build $@

speechtclient:
	go build $@

clean:
	rm -vf $(target)
