package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gorilla/websocket"
)

type Speech2Text struct {
	Creds string
	token string
}

func (s2t Speech2Text) Recognize(audio io.Reader, contentType string, texts chan<- string) error {
	const url = "wss://stream.watsonplatform.net/speech-to-text/api/v1/recognize?x-watson-learning-opt-out=1"
	hdr := http.Header{"X-Watson-Authorization-Token": []string{s2t.token}}
	conn, resp, err := websocket.DefaultDialer.Dial(url, hdr)

	switch err {
	case nil:
	case websocket.ErrBadHandshake:
		if resp.StatusCode == http.StatusUnauthorized {
			if s2t.token, err = token(s2t.Creds); err != nil {
				return err
			}
			return s2t.Recognize(audio, contentType, texts)
		}
	default:
		return err
	}

	defer conn.Close()

	start := &Action{
		Action:            "start",
		ContentType:       contentType,
		InactivityTimeout: -1,
		MaxAlternatives:   3,
		Continuous:        true,
		InterimResults:    true,
		WordConfidence:    true,
	}

	if err := conn.WriteJSON(start); err != nil {
		return err
	}

	errch := make(chan error, 2)
	go func() { errch <- s2t.recv(conn, texts) }()
	go func() { errch <- s2t.send(conn, audio, contentType) }()
	return <-errch
}

func (s2t Speech2Text) send(conn *websocket.Conn, audio io.Reader, contentType string) error {
	var buf bytes.Buffer
	for {
		buf.Reset()
		switch n, err := io.CopyN(&buf, audio, 1024*128); err {
		case nil:
			if err = conn.WriteMessage(websocket.BinaryMessage, buf.Bytes()); err != nil {
				return err
			}
			log.Printf("Sent %d audio bytes", n)
		case io.EOF, io.ErrUnexpectedEOF:
			log.Print("No audio bytes available. Sleeping 1s")
			time.Sleep(time.Second)
		default:
			return err
		}
	}
}

func (s2t Speech2Text) recv(conn *websocket.Conn, texts chan<- string) error {
	for {
		_, rc, err := conn.NextReader()
		if err != nil {
			return err
		}

		data, err := ioutil.ReadAll(rc)
		if err != nil {
			return nil
		}

		log.Printf("Received JSON: %s", data)

		var rs Response
		if err := json.Unmarshal(data, &rs); err != nil {
			return err
		}

		log.Printf("Received: %+v", rs)

		if text := rs.Text(); len(text) > 0 {
			texts <- text
		}
	}
}

type Action struct {
	Action            string `json:"action"`
	ContentType       string `json:"content-type"`
	InactivityTimeout int    `json:"inactivity_timeout"`
	MaxAlternatives   int    `json:"max_alternatives"`
	Continuous        bool   `json:"continuous"`
	InterimResults    bool   `json:"interim_results"`
	WordConfidence    bool   `json:"word_confidence"`
}

type Response struct {
	State    string   `json:"state"`
	Warnings []string `json:"warnings"`
	Results  []Result `json:"results"`
}

func (rs Response) Text() string {
	for _, r := range rs.Results {
		if len(r.Alternatives) > 0 {
			return r.Alternatives[0].Transcript
		}
	}
	return ""
}

type Result struct {
	Alternatives []Alternative `json:"alternatives"`
}

type Alternative struct {
	Transcript string `json:"transcript"`
}

func token(creds string) (string, error) {
	url := fmt.Sprintf("https://%s@stream.watsonplatform.net/authorization/api/v1/token?url=https://stream.watsonplatform.net/speech-to-text/api", creds)
	req, _ := http.NewRequest("GET", url, nil)

	dump, err := httputil.DumpRequestOut(req, false)
	if err != nil {
		return "", err
	}

	log.Print(string(dump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", errors.New(http.StatusText(resp.StatusCode))
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	return string(data), err
}
