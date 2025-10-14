package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const (
	redirectURI = "http://localhost:8080/callback"
)

var (
	auth  = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserTopRead))
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

type Track struct {
	ID       string
	Name     string
	Artists  []string
	Album    string
	Duration int
}

type Artist struct {
	ID         string
	Name       string
	Genres     []string
	Popularity int
	Followers  int
}

type SpotifyClient struct {
	client *spotify.Client
	ctx    context.Context
}

func main() {
	client, err := getSpotifyClient()
	if err != nil {
		log.Fatal("Failed to get Spotify client:", err)
	}

	spotifyClient := &SpotifyClient{
		client: client,
		ctx:    context.Background(),
	}

	fmt.Println("Fetching your Spotify data...")

	tracks := spotifyClient.fetchTopTracks()
	artists := spotifyClient.fetchTopArtists()

	displayData(tracks, artists)
}

func (sc *SpotifyClient) fetchTopTracks() []Track {
	topTracks, err := sc.client.CurrentUsersTopTracks(sc.ctx, spotify.Limit(10))
	if err != nil {
		log.Fatal("Error fetching top tracks:", err)
	}

	var tracks []Track
	for _, track := range topTracks.Tracks {
		artists := make([]string, len(track.Artists))
		for i, artist := range track.Artists {
			artists[i] = artist.Name
		}

		tracks = append(tracks, Track{
			ID:       string(track.ID),
			Name:     track.Name,
			Artists:  artists,
			Album:    track.Album.Name,
			Duration: track.Duration,
		})
	}

	return tracks
}

func (sc *SpotifyClient) fetchTopArtists() []Artist {
	topArtists, err := sc.client.CurrentUsersTopArtists(sc.ctx, spotify.Limit(10))
	if err != nil {
		log.Fatal("Error fetching top artists:", err)
	}

	var artists []Artist
	for _, artist := range topArtists.Artists {
		artists = append(artists, Artist{
			ID:         string(artist.ID),
			Name:       artist.Name,
			Genres:     artist.Genres,
			Popularity: artist.Popularity,
			Followers:  artist.Followers.Count,
		})
	}

	return artists
}

func displayData(tracks []Track, artists []Artist) {
	fmt.Println("\nYour Spotify Music Data")
	fmt.Println("=" + strings.Repeat("=", 30))

	fmt.Println("\nTop Tracks:")
	for i, track := range tracks {
		fmt.Printf("%d. %s - %s\n", i+1, track.Name, strings.Join(track.Artists, ", "))
	}

	fmt.Println("\nTop Artists:")
	for i, artist := range artists {
		fmt.Printf("%d. %s (Followers: %d)\n", i+1, artist.Name, artist.Followers)
	}
}

func getSpotifyClient() (*spotify.Client, error) {
	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:")
	fmt.Println(url)

	go func() {
		http.HandleFunc("/callback", completeAuth)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	client := <-ch
	return client, nil
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}
