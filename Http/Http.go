package Http

import (
	"github.com/gin-gonic/gin"
	"manydjs-poc/Services"
)
import "github.com/gin-gonic/contrib/static"

type Http struct {
	router         *gin.Engine
	spotifyService Services.SpotifyService
	jwtService     Services.JwtService
}

// NewHttp returns an instance of the HTTP struct
func NewHttp(spotifyService Services.SpotifyService) Http {
	h := Http{
		router:         gin.Default(),
		spotifyService: spotifyService,
		jwtService:     Services.NewJwtService(),
	}

	h.router.Use(static.Serve("/", static.LocalFile("./web", true)))

	h.addSpotifyRoutes()

	return h
}

// Run starts the HTTP server
func (h Http) Run() error {
	err := h.router.Run(":8080")
	if err != nil {
		return err
	}

	return nil
}
