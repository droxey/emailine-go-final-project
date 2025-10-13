# Proposal: A Spotify Music Data Viewer

I want to create a small tool that connects to Spotify and shows personalized music statistics.

I will build a command-line tool that connects to the Spotify Web API using OAuth.  

It will:
- Show my top tracks and artists from Spotify.
- Display my most-played genres.
- Generate a playlist with my current favorite songs.


## Tooling
- Go
- `net/http` for making API calls  
- `encoding/json` for handling responses  
- `golang.org/x/oauth2` for authentication 