package Services

import (
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

type SpotifyService struct {
	auth spotify.Authenticator
}

// NewSpotifyService returns an new instance of a SpotifyService
func NewSpotifyService() SpotifyService {
	auth := spotify.NewAuthenticator(os.Getenv("SPOTIFY_REDIRECT_URL"),
		spotify.ScopePlaylistModifyPrivate,
		spotify.ScopePlaylistReadPrivate,
		spotify.ScopeStreaming,
		spotify.ScopeUserReadEmail,
		spotify.ScopeUserModifyPlaybackState,
		spotify.ScopeUserReadCurrentlyPlaying,
		spotify.ScopeUserReadPlaybackState)

	auth.SetAuthInfo(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))

	service := SpotifyService{
		auth: auth,
	}

	return service
}

// GetClient returns a spotify.Client with the given oauth2.Token
func (s SpotifyService) GetClient(token *oauth2.Token) spotify.Client {
	return s.auth.NewClient(token)
}

// GetAuthUrl returns the redirect url for the OAuth2 workflow
func (s SpotifyService) GetAuthUrl() string {
	return s.auth.AuthURL("demoState")
}

// RedeemToken redeems a token and returns an oauth2.Token
func (s SpotifyService) RedeemToken(request *http.Request) (*oauth2.Token, error) {
	return s.auth.Token("demoState", request)
}
