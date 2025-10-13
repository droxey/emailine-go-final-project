# Proposal: A Spotify Music Data Viewer

I want to create a small command line tool that connects to Spotify Web API using OAuth.

It will:
- Show my top tracks and artists from Spotify.
- Display my most-played genres.
- Generate a playlist with my current favorite songs.


## Tooling
- Go
- `net/http` for making API calls  
- `encoding/json` for handling responses  
- `golang.org/x/oauth2` for authentication 