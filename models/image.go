package models

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirspock/dockerpedia-api/db"
	"gopkg.in/guregu/null.v3"
)

type Image struct {
	Id           null.Int    `json:"id"`
	Name         null.String `json:"name"`
	Last_updated null.Time   `json:"last_updated"`
	Full_size    null.Int    `json:"full_size"`
	Id_docker    null.Int    `json:"id_docker"`
	Image_id     null.Int    `json:"image_id"`
	Last_check   null.Time   `json:"last_check"`
	Status       null.Bool   `json:"status"`
	Last_try     null.Time   `json:"last_try"`
	Packages     null.Int    `json:"packages"`
	Critical     null.Int    `json:"critical"`
	DefCon1      null.Int    `json:"defcon1"`
	High         null.Int    `json:"high"`
	Low          null.Int    `json:"low"`
	Medium       null.Int    `json:"medium"`
	Negligible   null.Int    `json:"negligible"`
	Unknown      null.Int    `json:"unknown"`
	Score        null.Int    `json:"score"`
	Analysed     null.Bool   `json:"analysed"`
}

func getImageRepositorySQL(id int64, images *[]Image) {
	var image Image

	stmt, err := db.GetDB().Prepare("SELECT * FROM tag WHERE image_id=$1 and analysed ORDER BY SCORE ASC limit 5")
	rows, err := stmt.Query(id)

	if err != nil {
		fmt.Print(err.Error())
	}

	for rows.Next() {
		err := rows.Scan(
			&image.Id,
			&image.Name,
			&image.Last_updated,
			&image.Full_size,
			&image.Id_docker,
			&image.Image_id,
			&image.Last_check,
			&image.Status,
			&image.Last_try,
			&image.Packages,
			&image.Critical,
			&image.DefCon1,
			&image.High,
			&image.Low,
			&image.Medium,
			&image.Negligible,
			&image.Unknown,
			&image.Score,
			&image.Analysed,
		)

		*images = append(*images, image)
		if err != nil {
			fmt.Print(err.Error())
		}
	}
}

// fetchAllTodo fetch all todos
func FetchImagesRepository(c *gin.Context) {
	images := []Image{}

	query := c.Param("id")
	id, err := strconv.ParseInt(query, 10, 64)

	if err != nil {
		panic(err)
	}

	getImageRepositorySQL(int64(id), &images)

	c.JSON(http.StatusOK, gin.H{
		"result": images,
		"count":  len(images),
	})
}

func FetchImagesViz(c *gin.Context) {
	pattern := c.DefaultQuery("query", "")
	images := []Image{}
	repos := []Repository{}

	getRepositoriesPattern(&repos, pattern)

	for i := 0; i < len(repos); i++ {
		getImageRepositorySQL(repos[i].Id.ValueOrZero(), &images)
		for _, image := range images {
			repos[i].Images = append(repos[i].Images, image)
		}
	}

	c.JSON(http.StatusOK, repos)

}
