package finder

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strings"

	"github.com/zmb3/spotify/v2"
)

type Finder struct {
	client     *spotify.Client
	playlistID string
	r          *http.Request
}

type AudioFeatures struct {
	spotify.AudioFeatures
	ArtistID spotify.ID
}

type Track struct {
	AudioFeatures
	Genres []string
}

func New(client *spotify.Client, r *http.Request, playlistID string) *Finder {
	return &Finder{
		client:     client,
		playlistID: playlistID,
		r:          r,
	}
}

func (f *Finder) getPlaylist() (*spotify.PlaylistItemPage, error) {
	return f.client.GetPlaylistItems(f.r.Context(), spotify.ID(f.playlistID))
}

func (f *Finder) getAudioFeatures(trackIDs []spotify.ID) ([]*spotify.AudioFeatures, error) {
	return f.client.GetAudioFeatures(f.r.Context(), trackIDs...)
}

func (f *Finder) addGenres(tracks []AudioFeatures) ([]Track, error) {
	ids := make(map[spotify.ID]struct{})
	for _, track := range tracks {
		ids[track.ArtistID] = struct{}{}
	}

	var artistIDs []spotify.ID
	for id := range ids {
		artistIDs = append(artistIDs, id)
	}

	artists, err := f.client.GetArtists(f.r.Context(), artistIDs...)
	if err != nil {
		return nil, err
	}

	genres := make(map[spotify.ID][]string)
	for _, artist := range artists {
		genres[artist.ID] = artist.Genres
	}

	res := make([]Track, len(tracks))
	for i, track := range tracks {
		res[i] = Track{
			AudioFeatures: track,
			Genres:        genres[track.ArtistID],
		}
	}

	return res, nil
}

func (f *Finder) getTracks() ([]Track, error) {
	playlist, err := f.getPlaylist()
	if err != nil {
		return nil, err
	}

	trackIDs := make([]spotify.ID, len(playlist.Items))
	for i, item := range playlist.Items {
		trackIDs[i] = item.Track.Track.ID
	}

	audioFeatures, err := f.getAudioFeatures(trackIDs)
	if err != nil {
		return nil, err
	}

	// add genre id to tracks
	audioFeaturesWithAlbumID := make([]AudioFeatures, len(audioFeatures))
	for i, af := range audioFeatures {
		audioFeaturesWithAlbumID[i] = AudioFeatures{
			AudioFeatures: *af,
			ArtistID:      playlist.Items[i].Track.Track.Artists[0].ID,
		}
	}

	return f.addGenres(audioFeaturesWithAlbumID)
}

func (f *Finder) groupTrackFeatures(tracks []Track) map[string][]Track {
	// first we need to group the tracks by their genres
	genres := make(map[string][]Track)
	for _, track := range tracks {
		for _, genre := range track.Genres {
			genres[genre] = append(genres[genre], track)
		}
	}

	getUpperAndLowerBound := func(num float32) (float32, float32) {
		if num > 1 || num < 0 {
			upper := num + float32(math.Mod(float64(num), 10))
			return upper - 10, upper
		}
		upper := float32(math.Ceil(float64(num)*10) / 10)

		lower := upper - 0.1

		return lower, upper + 0.1
	}

	// now group by audio features
	audioFeaturesMap := make(map[string][]Track)
	for genre, tracks := range genres {
		for _, track := range tracks {
			// get the lower and upper bounds for each audio feature
			danceabilityLower, danceabilityUpper := getUpperAndLowerBound(track.Danceability)
			energyLower, energyUpper := getUpperAndLowerBound(track.Energy)
			loudnessLower, loudnessUpper := getUpperAndLowerBound(track.Loudness)
			speechinessLower, speechinessUpper := getUpperAndLowerBound(track.Speechiness)
			acousticnessLower, acousticnessUpper := getUpperAndLowerBound(track.Acousticness)
			instrumentalnessLower, instrumentalnessUpper := getUpperAndLowerBound(track.Instrumentalness)
			livenessLower, livenessUpper := getUpperAndLowerBound(track.Liveness)
			valenceLower, valenceUpper := getUpperAndLowerBound(track.Valence)
			tempoLower, tempoUpper := getUpperAndLowerBound(track.Tempo)

			// create a val for the audio features
			val := fmt.Sprintf(
				"%f_%f_%f_%f_%f_%f_%f_%f_%f_%f_%f_%f_%f_%f_%f_%f_%f_%f_%s",
				danceabilityLower, danceabilityUpper,
				energyLower, energyUpper,
				loudnessLower, loudnessUpper,
				speechinessLower, speechinessUpper,
				acousticnessLower, acousticnessUpper,
				instrumentalnessLower, instrumentalnessUpper,
				livenessLower, livenessUpper,
				valenceLower, valenceUpper,
				tempoLower, tempoUpper,
				genre,
			)

			// add the track to the audio features map
			if slc, ok := audioFeaturesMap[val]; !ok {
				audioFeaturesMap[val] = []Track{track}
			} else {
				audioFeaturesMap[val] = append(slc, track)
			}
		}
	}

	return audioFeaturesMap
}

func (f *Finder) Find() ([]spotify.SimpleTrack, error) {
	tracks, err := f.getTracks()
	if err != nil {
		return nil, err
	}

	audioFeaturesMap := f.groupTrackFeatures(tracks)

	// fmt.Println(audioFeaturesMap)

	// now randomly select 2 different items

	// gather total number of items in map
	totalItems := 0
	slc := make([][]string, 0, len(audioFeaturesMap))
	for g, tracks := range audioFeaturesMap {
		totalItems += len(tracks)
		for _, track := range tracks {
			slc = append(slc, []string{g, track.ID.String()})
		}
	}

	// randomly select 2 different items
	getRand := func() int {
		return rand.Intn(totalItems)
	}

	firstIndex := getRand()
	secondIndex := getRand()
	for firstIndex == secondIndex {
		secondIndex = getRand()
	}

	getTrackAttributes := func(val string) (*spotify.TrackAttributes, string) {
		vals := strings.Split(val, "_")
		// strToFloat := func(str string) float64 {
		// 	f, _ := strconv.ParseFloat(str, 32)
		// 	return float64(f)
		// }

		trackAttributes := spotify.NewTrackAttributes()
		// MinDanceability(strToFloat(vals[0])).
		// MaxDanceability(strToFloat(vals[1])).
		// MinEnergy(strToFloat(vals[2])).
		// MaxEnergy(strToFloat(vals[3])).
		// MinLoudness(strToFloat(vals[4])).
		// MaxLoudness(strToFloat(vals[5])).
		// MinSpeechiness(strToFloat(vals[6])).
		// MaxSpeechiness(strToFloat(vals[7])).
		// MinAcousticness(strToFloat(vals[8])).
		// MaxAcousticness(strToFloat(vals[9])).
		// MinInstrumentalness(strToFloat(vals[10])).
		// MaxInstrumentalness(strToFloat(vals[11])).
		// MinLiveness(strToFloat(vals[12])).
		// MaxLiveness(strToFloat(vals[13])).
		// MinValence(strToFloat(vals[14])).
		// MaxValence(strToFloat(vals[15])).
		// MinTempo(strToFloat(vals[16])).
		// MaxTempo(strToFloat(vals[17]))

		return trackAttributes, vals[18]
	}

	search := func(val, id string) (*spotify.Recommendations, error) {
		trackAttributes, genre := getTrackAttributes(val)
		seeds := spotify.Seeds{
			Tracks: []spotify.ID{spotify.ID(id)},
			Genres: []string{genre},
		}

		return f.client.GetRecommendations(f.r.Context(), seeds, trackAttributes, spotify.Limit(5))
	}

	first, second := slc[firstIndex], slc[secondIndex]
	firstCh, secondCh := make(chan *spotify.Recommendations), make(chan *spotify.Recommendations)

	go func() {
		recs, err := search(first[0], first[1])
		if err != nil {
			firstCh <- nil
			return
		}

		firstCh <- recs
	}()

	go func() {
		recs, err := search(second[0], second[1])
		if err != nil {
			secondCh <- nil
			return
		}

		secondCh <- recs
	}()

	combined := make([]spotify.SimpleTrack, 0, 10)

	for i := 0; i < 2; i++ {
		select {
		case recs := <-firstCh:
			if recs != nil {
				combined = append(combined, recs.Tracks...)
			}
		case recs := <-secondCh:
			if recs != nil {
				combined = append(combined, recs.Tracks...)
			}
		}
	}

	return combined, nil
}
