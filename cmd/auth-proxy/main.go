package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":9000", "HTTP(s) address to listen on")
	upstream := flag.String("upstream", "https://stream.watsonplatform.net/authorization/api/v1/token?url=https://stream.watsonplatform.net/speech-to-text/api", "Authenticated upstream DataPower Edge Router service URL")

	flag.Parse()

	http.ListenAndServe(*addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s]: %s %s %s", r.RemoteAddr, r.Method, r.Host, r.RequestURI)

		resp, err := http.Get(*upstream)
		if err != nil {
			respond(w, 500, err)
			return
		}

		if resp.StatusCode >= 400 {
			respond(w, resp.StatusCode, errors.New(http.StatusText(resp.StatusCode)))
			return
		}

		defer resp.Body.Close()
		if token, err := ioutil.ReadAll(resp.Body); err != nil {
			respond(w, 500, err)
		} else {
			respond(w, 200, string(token))
		}
	}))
}

func respond(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err, ok := v.(error); ok {
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
	} else if token, ok := v.(string); ok {
		w.Write([]byte(`{"token": "` + token + `"}`))
	}
}
