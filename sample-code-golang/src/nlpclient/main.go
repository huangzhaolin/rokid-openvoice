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
	host    = flag.String("host", "apigw.open.rokid.com:443", "Server address")
	lang    = flag.String("lang", "zh", "Language")
	text    = flag.String("text", "", "Text")
	cdomain = flag.String("cdomain", "", "Current domain")
	cred    = flag.Bool("cred", false, "Need credentials?")
	authit  = flag.Bool("auth", false, "Need auth?")
	tls     = flag.Bool("tls", false, "Need tls?")
	count   = flag.Int("count", 1, "Test count")
)

func call_nlp(conn *grpc.ClientConn, lang, text, cdomain string, authit bool) {
	seed := time.Now().UTC().UnixNano()
	clientrand := rand.New(rand.NewSource(seed))
	id := clientrand.Int31()

	client := pb.NewSpeechClient(conn)
	ctx := context.Background()
	if authit {
		sc := auth.NewSpeechCredential("nlp", "1.0")
		if ret, e := sc.Auth(client, ctx); e != nil || ret != 0 {
			log.Fatalf("client.Auth(): %d %v", ret, e)
		}
	}

	response, err := client.Nlp(ctx, &pb.NlpRequest{
		Header: &pb.NlpHeader{
			Id:      id,
			Lang:    lang,
			Cdomain: cdomain,
		},
		Asr: text,
	})
	if err != nil {
		log.Fatalf("client.Nlp(): %v", err)
	}
	log.Printf("Nlp(%s) = %s", text, response.Nlp)
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
		call_nlp(conn, *lang, *text, *cdomain, *authit)
	}
}
