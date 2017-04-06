package main

import (
	"flag"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	auth "auth"
	pb "pb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	host   = flag.String("host", "service-gw.rpc.rokid.com:30000", "Server address")
	lang   = flag.String("lang", "zh", "Language")
	codec  = flag.String("codec", "opu2", "Codec")
	text   = flag.String("text", "今天天气怎么样?", "Tts Text")
	fname  = flag.String("file", "", "Out file")
	cred   = flag.Bool("cred", false, "Need credentials?")
	authit = flag.Bool("auth", false, "Need auth?")
	tls    = flag.Bool("tls", false, "Need tls?")
	count  = flag.Int("count", 1, "Test count")
)

func call_tts(conn *grpc.ClientConn, lang, codec, text, fname string, authit bool) {
	var file *os.File

	if 0 != len(fname) {
		if f, err := os.Create(fname); err != nil {
			log.Fatalf("could not create file %v: %v", fname, err)
		} else {
			file = f
		}
		defer file.Close()
	}

	seed := time.Now().UTC().UnixNano()
	clientrand := rand.New(rand.NewSource(seed))
	id := clientrand.Int31()

	client := pb.NewSpeechClient(conn)
	ctx := context.Background()
	if authit {
		sc := auth.NewSpeechCredential("tts", "1.0")
		if ret, e := sc.Auth(client, ctx); e != nil || ret != 0 {
			log.Fatalf("client.Auth(): %d %v", ret, e)
		}
	}
	in := &pb.TtsRequest{Header: &pb.TtsHeader{Id: id, Declaimer: lang, Codec: codec}, Text: text}
	stream, err := client.Tts(ctx, in)
	if err != nil {
		log.Fatalf("client.Tts(): %v", err)
	}

	for {
		var data *pb.TtsResponse

		if data, err = stream.Recv(); err != nil {
			if err != io.EOF {
				log.Fatalf("stream.Recv(): %v", err)
			}
			break
		}
		log.Printf("Got tts: Text('%s'), Voice(len=%d)", data.Text, len(data.Voice))

		if file != nil {
			file.Write(data.Voice)
		}
	}
}

func do_conn(host string, cred, tls bool) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if cred {
		opts = append(opts, grpc.WithPerRPCCredentials(new(auth.SpeechCredential)))
	}
	if tls {
		var sn string
		creds := credentials.NewClientTLSFromCert(nil, sn)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		log.Fatalf("fail to dial('%s'): %v", host, err)
	}

	return conn, nil
}

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	flag.Parse()

	conn, err := do_conn(*host, *cred, *tls)
	if err != nil {
		return
	}
	defer conn.Close()

	for i := 0; i < *count; i += 1 {
		call_tts(conn, *lang, *codec, *text, *fname, *authit)
	}
}
