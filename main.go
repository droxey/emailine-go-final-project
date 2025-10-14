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

func main() {
	client, err := getSpotifyClient()
	if err != nil {
		log.Fatal("Failed to get Spotify client:", err)
	}

	ctx := context.Background()

	fmt.Println("Fetching your Spotify data...")

	topTracks, err := client.CurrentUsersTopTracks(ctx, spotify.Limit(10))
	if err != nil {
		log.Fatal("Error fetching top tracks:", err)
	}

	topArtists, err := client.CurrentUsersTopArtists(ctx, spotify.Limit(10))
	if err != nil {
		log.Fatal("Error fetching top artists:", err)
	}

	fmt.Println("\nYour Spotify Music Data")
	fmt.Println("=" + strings.Repeat("=", 30))

	fmt.Println("\nTop Tracks:")
	for i, track := range topTracks.Tracks {
		artists := make([]string, len(track.Artists))
		for j, artist := range track.Artists {
			artists[j] = artist.Name
		}
		fmt.Printf("%d. %s - %s\n", i+1, track.Name, strings.Join(artists, ", "))
	}

	fmt.Println("\nTop Artists:")
	for i, artist := range topArtists.Artists {
		fmt.Printf("%d. %s\n", i+1, artist.Name)
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
