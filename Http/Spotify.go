package Http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// addSpotifyRoutes registers the routes for the spotify api
func (h Http) addSpotifyRoutes() {
	group := h.router.Group("/api")

	group.GET("/login", h.loginHandler)
	group.GET("/redeem-token", h.redeemTokenHandler)

	authorized := group.Group("/")
	authorized.Use(h.jwtService.AuthorizeJWT())
	authorized.GET("/me", h.meHandler)
}

// meHandler shows information about the spotify session of the current user
func (h Http) meHandler(c *gin.Context) {
	token := h.jwtService.GetPayload(c)
	client := h.spotifyService.GetClient(token)

	options, err := client.PlayerState()
	if err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{
		"playerState": options,
	})
}

// redeemTokenHandler redeems a OAuth2 token and redirects the client with it in the url query string
func (h Http) redeemTokenHandler(c *gin.Context) {
	oAuthToken, err := h.spotifyService.RedeemToken(c.Request)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token := h.jwtService.GenerateToken(oAuthToken)

	c.Redirect(http.StatusFound, fmt.Sprintf("/?token=%s", token))
}

// loginHandler generates the login url which the user is redirected to
func (h Http) loginHandler(c *gin.Context) {
	url := h.spotifyService.GetAuthUrl()

	c.JSON(200, gin.H{
		"url": url,
	})
}
