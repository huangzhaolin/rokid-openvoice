package main

import (
	"flag"
	"log"
	"math/rand"
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
	vt      = flag.String("vt", "", "Voice trigger")
	text    = flag.String("text", "", "Text")
	cred    = flag.Bool("cred", false, "Need credentials?")
	authit  = flag.Bool("auth", false, "Need auth?")
	cdomain = flag.String("cdomain", "", "Current domain")
	device  = flag.String("device", "", "Device info")
	tls     = flag.Bool("tls", false, "Need tls?")
	count   = flag.Int("count", 1, "Test count")
)

func call_speecht(conn *grpc.ClientConn, lang, vt, cdomain, device, text string, authit bool) {
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

	response, err := client.Speecht(ctx, &pb.TextSpeechRequest{
		Header: &pb.SpeechHeader{
			Id:      id,
			Lang:    lang,
			Vt:      vt,
			Cdomain: cdomain,
			Device:  device,
		},
		Asr: text,
	})
	if err != nil {
		log.Fatalf("client.Speecht(): %v", err)
	}
	log.Printf("Speecht(%s) = asr(%s), nlp(%s), action(%s)", text, response.Asr, response.Nlp, response.Action)
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
		call_speecht(conn, *lang, *vt, *cdomain, *device, *text, *authit)
	}
}
