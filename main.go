package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirspock/dockerpedia-api/db"
	"github.com/sirspock/dockerpedia-api/models"
)

func main() {
	db.Init()

	router := gin.Default()

	v1 := router.Group("/api/v1/")
	{
		v1.GET("/repositories/", models.FetchImagesRepository)
		v1.GET("/repositories/:id/", models.FetchRepository)
		v1.GET("/repositories/:id/images", models.FetchImagesRepository)
	}

	router.Run(":8080")
}
