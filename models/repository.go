package models

import (
	"fmt"
	"net/http"

	"gopkg.in/guregu/null.v3"

	"github.com/gin-gonic/gin"
	"github.com/sirspock/dockerpedia-api/db"
)

type Repository struct {
	Id              null.Int    `json:"id"`
	Name            null.String `json:"name"`
	Full_name       null.String `json:"full_name"`
	Namespace       null.String `json:"namespace"`
	User            null.String `json:"user"`
	Affiliation     null.String `json:"affilation"`
	Description     null.String `json:"description"`
	Is_automated    null.Bool   `json:"is_automated"`
	Last_updated    null.Time   `json:"last_updated"`
	Pull_count      null.Int    `json:"pull_count"`
	Repository_type null.String `json:"repository_type"`
	Star_count      null.Int    `json:"start_count"`
	Status          null.Bool   `json:"status"`
	Tags_checked    null.Time   `json:"tags_checked"`
	Official        null.Bool   `json:"official"`
	Score           null.Int    `json:"score"`
	Images          []Image     `json:"images"`
}

func SearchRepository(c *gin.Context) {
	pattern := c.DefaultQuery("query", "mysql")

	var (
		repo  Repository
		repos []Repository
	)

	fmt.Println(pattern)

	stmt, err := db.GetDB().Prepare("SELECT * FROM image WHERE name like '%' || $1 || '%' ORDER BY pull_count DESC limit 20")
	rows, err := stmt.Query(pattern)

	if err != nil {
		fmt.Print(err.Error())
	}

	for rows.Next() {
		err = rows.Scan(
			&repo.Id,
			&repo.Name,
			&repo.Full_name,
			&repo.Namespace,
			&repo.User,
			&repo.Affiliation,
			&repo.Description,
			&repo.Is_automated,
			&repo.Last_updated,
			&repo.Pull_count,
			&repo.Repository_type,
			&repo.Star_count,
			&repo.Status,
			&repo.Tags_checked,
			&repo.Official,
			&repo.Score,
		)

		repos = append(repos, repo)
		if err != nil {
			fmt.Print(err.Error())
		}
	}

	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"result": repos,
		"count":  len(repos),
	})
}

func getRepositoriesPattern(repos *[]Repository, pattern string) {
	var repo Repository
	stmt, err := db.GetDB().Prepare("SELECT * FROM image WHERE namespace like '%' || $1 || '%' ORDER BY pull_count DESC limit 10")
	rows, err := stmt.Query(pattern)

	if err != nil {
		fmt.Print(err.Error())
	}

	for rows.Next() {
		err = rows.Scan(
			&repo.Id,
			&repo.Name,
			&repo.Full_name,
			&repo.Namespace,
			&repo.User,
			&repo.Affiliation,
			&repo.Description,
			&repo.Is_automated,
			&repo.Last_updated,
			&repo.Pull_count,
			&repo.Repository_type,
			&repo.Star_count,
			&repo.Status,
			&repo.Tags_checked,
			&repo.Official,
			&repo.Score,
		)
		*repos = append(*repos, repo)
	}
	defer rows.Close()
}

// fetchAllTodo fetch all todos
func FetchRepository(c *gin.Context) {
	var (
		repo Repository
	)

	id := c.Param("id")
	sqlStatement := `SELECT
        id, name, namespace, full_name, user, description, is_automated,
        last_updated, pull_count, star_count, tags_checked, score, official
        FROM image WHERE id=$1 LIMIT 20;`
	row := db.GetDB().QueryRow(sqlStatement, id)

	err := row.Scan(
		&repo.Id,
		&repo.Name,
		&repo.Namespace,
		&repo.Full_name,
		&repo.User,
		&repo.Description,
		&repo.Is_automated,
		&repo.Last_updated,
		&repo.Pull_count,
		&repo.Star_count,
		&repo.Tags_checked,
		&repo.Score,
		&repo.Official,
	)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, repo)

}
