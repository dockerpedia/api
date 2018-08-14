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
	Repositories Repositories `json:"children"`
}
type Repository struct {
	Id              null.Int    `json:"id",omitempty`
	Name            null.String `json:"name",omitempty`
	Full_name       null.String `json:"full_name",omitempty`
	Namespace       null.String `json:"namespace",omitempty`
	User            null.String `json:"user",omitempty`
	Affiliation     null.String `json:"affilation",omitempty`
	Description     null.String `json:"description",omitempty`
	Is_automated    null.Bool   `json:"is_automated",omitempty`
	Last_updated    null.Time   `json:"last_updated",omitempty`
	Pull_count      null.Int    `json:"pull_count",omitempty`
	Repository_type null.String `json:"repository_type",omitempty`
	Star_count      null.Int    `json:"start_count",omitempty`
	Status          null.Bool   `json:"status",omitempty`
	Tags_checked    null.Time   `json:"tags_checked",omitempty`
	Official        null.Bool   `json:"official",omitempty`
	Score           null.Int    `json:"value",omitempty`
	Images          []Image     `json:"children",omitempty`
	Full_size       null.Int    `json:"full_size",omitempty`
	Analysed        null.Bool   `json:"is_automated",omitempty`
}
type Repositories []Repository

type RepositoryQuery struct {
	Namespace       null.String `json:"namespace",omitempty`
}

func (repositories *Repositories) ModifyImage(images map[int64][]Image){
	for i := 0; i < len(*repositories); i++ {
		(*repositories)[i].Images = images[(*repositories)[i].Id.Int64]
	}

}


func SearchUser(c *gin.Context) {
	pattern := c.DefaultQuery("query", "mysql")

	var (
		repo  RepositoryQuery
		repos []string
	)


	stmt, err := db.GetDB().Prepare(`SELECT DISTINCT namespace FROM image 
	WHERE analysed='t' AND namespace like LOWER($1 || '%')`)

	rows, err := stmt.Query(pattern)

	if err != nil {
		log.Print(err.Error())
	}

	for rows.Next() {
		err = rows.Scan(
			&repo.Namespace,
		)

		repos = append(repos, repo.Namespace.String)
		if err != nil {
			log.Print(err.Error())
		}
	}

	defer rows.Close()

	c.JSON(http.StatusOK,  gin.H{
		"result": repos,
		"count":  len(repos),
	})
}

func SearchRepository(c *gin.Context) {
	pattern := c.DefaultQuery("query", "mysql")

	var (
		repo  Repository
		repos Repositories
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
	WHERE namespace=$1 	AND analysed='t' ORDER BY pull_count
			`)
		if err != nil {
			return nil, err
		}
		rows, err := stmt.Query(search)
		return rows, err
	}

}


func getRepoImages(repos *Repositories, imageRepoIds *[]int, search string, numberImages int, packages bool) {
	var repo Repository
	var image Image
	repoHash := make(map[int64]bool)
	imagesHash := make(map[int64][]Image)

	var rows *sql.Rows
	var err error
	if packages {
		stmt, _ := db.GetDB().Prepare(`
SELECT image.*, tag.* from image JOIN lateral 
(select * from tag where tag.image_id=image.id and tag.analysed ORDER BY score DESC limit $1) tag on true 
WHERE LOWER(image.name) like LOWER('%' || $2 || '%')  ORDER BY pull_count DESC`)
		rows, err = stmt.Query(numberImages, search)

	} else {
		stmt, _ := db.GetDB().Prepare(`SELECT image.*, tag.* from image JOIN lateral 
(select * from tag where tag.image_id=image.id and tag.analysed limit $1) tag on true 
WHERE image.namespace=$2`)
		rows, err = stmt.Query(numberImages, search)
	}

	if err != nil {
		log.Println("error join query", err)
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
			&image.OperatingSystem,
		)
		repoId := repo.Id.Int64

		if _, ok := repoHash[repoId]; !ok {
			*repos = append(*repos, repo)
			repoHash[repoId] = true
		}
		calculateRisk(&image)
		imagesHash[repoId] = append(imagesHash[repoId], image)

	}

	defer rows.Close()
	repos.ModifyImage(imagesHash)

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
