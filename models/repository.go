package models

import (
	"database/sql"
		"log"
	"net/http"

	"gopkg.in/guregu/null.v3"

	"github.com/dockerpedia/api/db"
	"github.com/gin-gonic/gin"
)

type RepositorySearchResult struct {
	Name         null.String  `json:"name"`
	Repositories []Repository `json:"children"`
}

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
	Score           null.Int    `json:"value"`
	Images          []Image     `json:"children"`
	Full_size       null.Int    `json:"full_size"`
	Analysed        null.Bool   `json:"is_automated"`
}

func SearchRepository(c *gin.Context) {
	pattern := c.DefaultQuery("query", "mysql")

	var (
		repo  Repository
		repos []Repository
	)


	stmt, err := db.GetDB().Prepare(`SELECT * FROM image 
	WHERE LOWER(name) like LOWER('%' || $1 || '%')  ORDER BY pull_count DESC limit 2`)

	rows, err := stmt.Query(pattern)

	if err != nil {
		log.Print(err.Error())
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
			&repo.Analysed,
		)

		repos = append(repos, repo)
		if err != nil {
			log.Print(err.Error())
		}
	}

	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"result": repos,
		"count":  len(repos),
	})
}

func getRepositoryPatternQuery(search string, pattern bool) (*sql.Rows, error) {
	if pattern {
		stmt, err := db.GetDB().Prepare(`
	SELECT * FROM image
	WHERE LOWER(name) like LOWER('%' || $1 || '%')
	AND analysed='t' ORDER BY pull_count DESC
		`)
		if err != nil {
			return nil, err
		}
		rows, err := stmt.Query(search)
		return rows, err

	} else {
		stmt, err := db.GetDB().Prepare(`
	SELECT * FROM image
	WHERE namespace=$1 ORDER BY pull_count
			`)
		if err != nil {
			return nil, err
		}
		rows, err := stmt.Query(search)
		return rows, err
	}

}

func getRepositoriesPattern(repos *[]Repository, search string, packages bool) {
	var repo Repository
	rows, err := getRepositoryPatternQuery(search, packages)
	if err != nil {
		log.Printf("error: %s", err)
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
			&repo.Analysed,
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
        FROM image WHERE id=$1 LIMIT 2;`
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
