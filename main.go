package main

import (
	"github.com/dockerpedia/api/db"
	"github.com/dockerpedia/api/models"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db.Init()

	router := gin.Default()

	v1 := router.Group("/api/v1/")
	{
		v1.GET("/repositories/", models.SearchRepository)
		v1.GET("/repositories/:id/", models.FetchRepository)
		v1.GET("/repositories/:id/images", models.FetchImagesRepository)
		v1.GET("/images/:id", models.FetchImage)
		v1.GET("/images/:id/vulnerabilities", models.FetchImagesVulns)
		v1.GET("/vulnerability/:id", models.FetchVulnerability)
		v1.GET("/images/:id/packages", models.FetchImagesPackages)

		v1.GET("/viz", models.FetchImagesViz)
		v1.POST("/viz", models.FetchImagesVizPost)
		v1.POST("/viz2", models.FetchImagesVizPostv2)
	}

	router.Run(":8080")
}
