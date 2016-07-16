package main

import (
	"math/rand"
	"time"

	"github.com/tsenart/spotify"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Recommender struct{ spotify.Client }

func (r *Recommender) Tracks(tone string, limit int) ([]*spotify.SimpleTrack, error) {
	genres := make([]string, 5)
	for i, idx := range rand.Perm(len(seedGenres))[:len(genres)] {
		genres[i] = seedGenres[idx]
	}
	attrs := map[string]interface{}{"market": "DE"}
	for attr, val := range toneAttrs[tone] {
		attrs[attr] = val
	}
	return r.Client.GetRecommendations(limit, genres, attrs)
}

var toneAttrs = map[string]map[string]interface{}{
	"disgust": {
		"min_energy": 0.4, "max_energy": 0.6,
		"min_valence": 0.0, "max_valence": 0.2,
	},
	"anger": {
		"min_energy": 0.7, "max_energy": 0.9,
		"min_valence": 0.0, "max_valence": 0.2,
	},
	"fear": {
		"min_energy": 0.8, "max_energy": 1.0,
		"min_valence": 0.4, "max_valence": 0.6,
	},
	"joy": {
		"min_energy": 0.4, "max_energy": 0.6,
		"min_valence": 0.8, "max_valence": 1.0,
	},
	"sadness": {
		"min_energy": 0.2, "max_energy": 0.4,
		"min_valence": 0.0, "max_valence": 0.2,
	},
}

var seedGenres = []string{
	"acoustic",
	"afrobeat",
	"alt-rock",
	"alternative",
	"ambient",
	"anime",
	"black-metal",
	"bluegrass",
	"blues",
	"bossanova",
	"brazil",
	"breakbeat",
	"british",
	"cantopop",
	"chicago-house",
	"children",
	"chill",
	"classical",
	"club",
	"comedy",
	"country",
	"dance",
	"dancehall",
	"death-metal",
	"deep-house",
	"detroit-techno",
	"disco",
	"disney",
	"drum-and-bass",
	"dub",
	"dubstep",
	"edm",
	"electro",
	"electronic",
	"emo",
	"folk",
	"forro",
	"french",
	"funk",
	"garage",
	"german",
	"gospel",
	"goth",
	"grindcore",
	"groove",
	"grunge",
	"guitar",
	"happy",
	"hard-rock",
	"hardcore",
	"hardstyle",
	"heavy-metal",
	"hip-hop",
	"holidays",
	"honky-tonk",
	"house",
	"idm",
	"indian",
	"indie",
	"indie-pop",
	"industrial",
	"iranian",
	"j-dance",
	"j-idol",
	"j-pop",
	"j-rock",
	"jazz",
	"k-pop",
	"kids",
	"latin",
	"latino",
	"malay",
	"mandopop",
	"metal",
	"metal-misc",
	"metalcore",
	"minimal-techno",
	"movies",
	"mpb",
	"new-age",
	"new-release",
	"opera",
	"pagode",
	"party",
	"philippines-opm",
	"piano",
	"pop",
	"pop-film",
	"post-dubstep",
	"power-pop",
	"progressive-house",
	"psych-rock",
	"punk",
	"punk-rock",
	"r-n-b",
	"rainy-day",
	"reggae",
	"reggaeton",
	"road-trip",
	"rock",
	"rock-n-roll",
	"rockabilly",
	"romance",
	"sad",
	"salsa",
	"samba",
	"sertanejo",
	"show-tunes",
	"singer-songwriter",
	"ska",
	"sleep",
	"songwriter",
	"soul",
	"soundtracks",
	"spanish",
	"study",
	"summer",
	"swedish",
	"synth-pop",
	"tango",
	"techno",
	"trance",
	"trip-hop",
	"turkish",
	"work-out",
	"world-music",
}
