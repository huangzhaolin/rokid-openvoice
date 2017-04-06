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
	host    = flag.String("host", "service-gw.rpc.rokid.com:30000", "Server address")
	lang    = flag.String("lang", "zh", "Language")
	codec   = flag.String("codec", "pcm", "Codec")
	fname   = flag.String("file", "", "Audio file")
	cred    = flag.Bool("cred", false, "Need credentials?")
	authit  = flag.Bool("auth", false, "Need auth?")
	cdomain = flag.String("cdomain", "", "Current domain")
	device  = flag.String("device", "", "Device info")
	tls     = flag.Bool("tls", false, "Need tls?")
	count   = flag.Int("count", 1, "Test count")
)

func call_speechv(conn *grpc.ClientConn, lang, codec, cdomain, device, fname string, authit bool) {
	var file *os.File

	if f, err := os.Open(fname); err != nil {
		log.Fatalf("could not open file %v: %v", fname, err)
	} else {
		file = f
	}
	defer file.Close()

	client := pb.NewSpeechClient(conn)
	ctx := context.Background()
	if authit {
		sc := auth.NewSpeechCredential("tts", "1.0")
		if ret, e := sc.Auth(client, ctx); e != nil || ret != 0 {
			log.Fatalf("client.Auth(): %d %v", ret, e)
		}
	}
	stream, err := client.Speechv(ctx)
	if err != nil {
		log.Fatalf("client.Speechv(): %v", err)
	}

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					// read done.
					close(waitc)
					break
				} else {
					log.Fatalf("stream.Recv(): %v", err)
				}
			}
			log.Printf("Got asr: asr(%s), nlp(%s), action(%s)", in.Asr, in.Nlp, in.Action)
		}
	}()

	seed := time.Now().UTC().UnixNano()
	clientrand := rand.New(rand.NewSource(seed))
	id := clientrand.Int31()

	req := &pb.VoiceSpeechRequest{
		RequestContent: &pb.VoiceSpeechRequest_Header{
			Header: &pb.SpeechHeader{
				Id:      id,
				Lang:    lang,
				Codec:   codec,
				Cdomain: cdomain,
				Device:  device,
			},
		},
	}
	if err := stream.Send(req); err != nil {
		log.Fatalf("Failed to send header: %v", err)
	}

	voice := make([]byte, 320*30)
	for {
		time.Sleep(300 * time.Millisecond)

		if n, err := file.Read(voice[:]); err == nil {
			log.Printf("Read file(%d)", n)
			req := &pb.VoiceSpeechRequest{RequestContent: &pb.VoiceSpeechRequest_Voice{Voice: voice[:n]}}
			if err = stream.Send(req); err != nil {
				log.Fatalf("Failed to send voice: %v", err)
			}
		} else {
			log.Printf("Read file: %v", err)
			break
		}
	}
	stream.CloseSend()

	<-waitc
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
		call_speechv(conn, *lang, *codec, *cdomain, *device, *fname, *authit)
	}
}
