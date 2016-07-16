package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/tsenart/spotify"

	"golang.org/x/oauth2"
)

func main() {
	contentType := flag.String("content.type", "audio/l16; rate=44100; channels=2", "Input audio content type")
	s2tcreds := flag.String("speech2text.creds", "", "Watson Speech to Text service credentials")
	tonecreds := flag.String("toneanalyzer.creds", "", "Watson Tone analyzer service credentials")
	spotifyToken := flag.String("spotify.token", "", "Spotify access token")
	minWords := flag.Int("min.words", 5, "Minimum number of words to send to Tone Analysis")

	flag.Parse()

	err := run(os.Stdin, *s2tcreds, *contentType, *tonecreds, *spotifyToken, *minWords)
	if err != nil {
		log.Fatal(err)
	}
}

func run(audio io.Reader, s2tcreds, contentType, tonecreds, spotifyToken string, minWords int) error {
	s2t := Speech2Text{Creds: s2tcreds}
	ta := ToneAnalyzer{Creds: tonecreds}
	rec := Recommender{
		Client: spotify.NewAuthenticator("http://localhost:9090/").
			NewClient(&oauth2.Token{AccessToken: spotifyToken, TokenType: "Bearer"}),
	}

	texts := make(chan string)
	errch := make(chan error)

	go func() { errch <- s2t.Recognize(audio, contentType, texts) }()

	for {
		select {
		case text := <-texts:
			if words := strings.Split(text, " "); len(words) < minWords {
				log.Printf("%q has less than %d words. Ignoring", text, minWords)
				continue
			}

			tone, err := ta.Analyze(text)
			if err != nil {
				log.Printf("Failed analysing tone for %q: %v", text, err)
				continue
			}

			log.Printf("Tone for %q is %q", text, tone)

			tracks, err := rec.Tracks(tone, 1)
			if err != nil || len(tracks) != 1 {
				log.Printf("Failed getting reccomendations for %q: %v", tone, err)
				continue
			}

			track := tracks[0]
			log.Printf("Recommended track: %+v", track)

			resp, err := http.Get(track.PreviewURL)
			if err != nil {
				log.Printf("Failed getting track data: %v", err)
				continue
			}

			if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
				log.Printf("Failed playing track: %v", err)
			}
			resp.Body.Close()

		case err := <-errch:
			return err
		}
	}
}
