package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type SpotifyToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func getAuthHeader(clientID, clientSecret string) string {
	authString := fmt.Sprintf("%s:%s", clientID, clientSecret)
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authString))
	return "Basic " + encodedAuth
}

func GetSpotifyToken() (*SpotifyToken, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load .env")
	}

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("missing Spotify credentials")
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", getAuthHeader(clientID, clientSecret))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("spotify error: %s", string(body))
	}

	var token SpotifyToken
	if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

type Playlist struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	ExternalUrls ExternalUrls `json:"external_urls"`
	Owner        Owner        `json:"owner"`
	Tracks       TracksObject `json:"tracks"`
}

type TracksObject struct {
	Items []PlaylistItem `json:"items"`
}

type PlaylistItem struct {
	Track Track `json:"track"`
}

type ExternalUrls struct {
	Spotify string `json:"spotify"`
}

type Owner struct {
	DisplayName    string       `json:"display_name"`
	ID             string       `json:"id"`
	ExExternalUrls ExternalUrls `json:"external_urls"`
}

type Artist struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	ExternalUrls ExternalUrls `json:"external_urls"`
}

type Album struct {
	ExternalUrls ExternalUrls `json:"external_urls"`
	Artists      []Artist     `json:"artists"`
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	ReleaseDate  string       `json:"release_date"`
}

type Track struct {
	Name         string       `json:"name"`
	DurationMs   int          `json:"duration_ms"`
	ExternalUrls ExternalUrls `json:"external_urls"`
	Artists      []Artist     `json:"artists"`
	Album        Album        `json:"album"`
}

func GetMusics(ctx *gin.Context) {
	fmt.Print("hello, i got request")

	token, err := GetSpotifyToken()
	if err != nil {
		log.Fatal(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "try later"})
		return
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/0cwPcui7aGHkmfHZiD3Hb9", nil)
	if err != nil {
		log.Fatal(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "try later"})
		return
	}

	authStr := fmt.Sprintf("Bearer %s", token.AccessToken)

	req.Header.Set("Authorization", authStr)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "try later"})
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "try later", "body": string(body)})
		return
	}

	var playlist Playlist
	if err := json.NewDecoder(res.Body).Decode(&playlist); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "try later", "error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success", "body": playlist})
}

func GetPlaylist(ctx *gin.Context) {
	playlistID, ok := ctx.Params.Get("playlistID")
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
		return
	}

	token, err := GetSpotifyToken()
	if err != nil {
		log.Fatal(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "try later"})
		return
	}

	spotifyPlaylistEndpoint := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s", playlistID)

	req, err := http.NewRequest("GET", spotifyPlaylistEndpoint, nil)
	if err != nil {
		log.Fatal(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "try later"})
		return
	}

	authStr := fmt.Sprintf("Bearer %s", token.AccessToken)

	req.Header.Set("Authorization", authStr)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "try later"})
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "try later", "body": string(body)})
		return
	}

	var playlist Playlist
	if err := json.NewDecoder(res.Body).Decode(&playlist); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "try later", "error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success", "body": playlist})
}
