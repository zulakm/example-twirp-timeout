package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/twitchtv/twirp"
	"github.com/zulakm/example-twirp-timeout/rpc/haberdasher"
)

func main() {
	server := &haberdasherServer{}
	twirpHandler := haberdasher.NewHaberdasherServer(server)
	http.ListenAndServe(":8080", WithTwirpTimeout(twirpHandler))
}

type haberdasherServer struct{}

func (s *haberdasherServer) MakeHat(ctx context.Context, size *haberdasher.Size) (*haberdasher.Hat, error) {
	select {
	case <-ctx.Done():
		return nil, twirp.InternalError("deadline exceeded")
	case <-time.After(100 * time.Millisecond):

		if size.Inches <= 0 {
			return nil, twirp.InvalidArgumentError("inches", "I can't make a hat that small!")
		}
		return &haberdasher.Hat{
			Inches: size.Inches,
			Color:  []string{"white", "black", "brown", "red", "blue"}[rand.Intn(4)],
			Name:   []string{"bowler", "baseball cap", "top hat", "derby"}[rand.Intn(3)],
		}, nil
	}
}

func WithTwirpTimeout(base http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if timeoutStr := r.Header.Get("Twirp-Timeout"); timeoutStr != "" {
			if timeout, err := time.ParseDuration(timeoutStr); err == nil {
				ctx, cancel := context.WithTimeout(r.Context(), timeout)
				defer cancel()
				r = r.WithContext(ctx)
				log.Printf("debug: set request timeout (%v)", timeout)
			}
		}

		base.ServeHTTP(w, r)
	})
}
