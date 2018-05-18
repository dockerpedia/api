package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirspock/dockerpedia-api/db"
	null "gopkg.in/guregu/null.v3"
)

type Image struct {
	Id           null.Int    `json:"id"`
	Name         null.String `json:"name"`
	Last_updated null.Time   `json:"last_updated"`
	Full_size    null.Int    `json:"full_size"`
	Id_docker    null.Int    `json:"id_docker"`
	Image_id     null.Int    `json:"image_id"`
	Last_check   time.Time   `json:"last_check"`
	Status       null.Bool   `json:"status"`
	Last_try     time.Time   `json:"last_try"`
	Packages     null.Int    `json:"packages"`
	Critical     null.Int    `json:"critical"`
	DefCon1      null.Int    `json:"defcon1"`
	High         null.Int    `json:"high"`
	Low          null.Int    `json:"low"`
	Medium       null.Int    `json:"medium"`
	Negligible   null.Int    `json:"negligible"`
	Unknown      null.Int    `json:"unknown"`
}

// fetchAllTodo fetch all todos
func FetchImagesRepository(c *gin.Context) {
	var (
		image  Image
		images []Image
	)

	id := c.Param("id")
	fmt.Println(id)

	stmt, err := db.GetDB().Prepare("SELECT * FROM tag WHERE image_id=$1 limit 10")
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
		)

		images = append(images, image)
		if err != nil {
			fmt.Print(err.Error())
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"result": images,
		"count":  len(images),
	})
}
