package main

import (
	"log"
	"net/http"
	"wp-media-core/database"
	gcstorage "wp-media-core/gc-storage"
	"wp-media-core/media"

	_ "github.com/lib/pq"
)

const apiBasePath = "/api"

func main() {
	database.IntDB()
	gcstorage.IntGCS()

	media.SetupRoutes(apiBasePath)

	log.Fatalln(http.ListenAndServe(":3050", nil))
	log.Printf("server running on port %d...", 3050)
}
