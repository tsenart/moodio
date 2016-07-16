#!/bin/bash

set -eu

go build
rm -f input.wav
ffmpeg -loglevel quiet -f avfoundation -i ":0" input.wav &
tail -n +1 -F input.wav | ./moodio -content.type="audio/wav" \
 -speech2text.creds="$SPEECH2TEXT_CREDS" \
 -toneanalyzer.creds="$TONEANALYSER_CREDS" \
 -spotify.token="$SPOTIFY_TOKEN" | tee output.mp3 | mplayer -msglevel all=0 -cache 1024 -
