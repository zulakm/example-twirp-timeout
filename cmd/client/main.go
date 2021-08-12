package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/twitchtv/twirp"
	"github.com/zulakm/example-twirp-timeout/rpc/haberdasher"
)

func main() {
	flag.Parse()
	serverURL := flag.Arg(0)

	client := haberdasher.NewHaberdasherProtobufClient(serverURL, &http.Client{})
	ctx := context.Background()

	log.Println("calling MakeHat with 1s timeout")
	hat, err := client.MakeHat(WithTwirpTimeout(ctx, 1*time.Second), &haberdasher.Size{Inches: 12})
	if err != nil {
		log.Printf("error: %v", err)
	} else {
		log.Printf("got hat: %v", hat.String())
	}

	log.Println("calling MakeHat with 10ms timeout")
	hat, err = client.MakeHat(WithTwirpTimeout(ctx, 10*time.Millisecond), &haberdasher.Size{Inches: 6})
	if err != nil {
		log.Printf("error: %v", err)
	} else {
		log.Printf("got hat: %v", hat.String())
	}
}

func WithTwirpTimeout(ctx context.Context, timeout time.Duration) context.Context {
	h := make(http.Header)
	h.Add("Twirp-Timeout", timeout.String())
	ctx, _ = twirp.WithHTTPRequestHeaders(ctx, h)
	return ctx
}
