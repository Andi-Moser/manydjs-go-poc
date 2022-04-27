package main

import (
	"github.com/joho/godotenv"
	"manydjs-poc/Http"
	"manydjs-poc/Services"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	spotifyService := Services.NewSpotifyService()
	http := Http.NewHttp(spotifyService)

	err = http.Run()
	if err != nil {
		panic(err)
	}
}
