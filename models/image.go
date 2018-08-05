package models

import (
	"fmt"
				"github.com/dockerpedia/api/db"
		"gopkg.in/guregu/null.v3"
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
	"log"
)

const MAXVALUE = 1000
const NUMBER_IMAGES = 20

type Request struct {
	User    string `form:"user" json:"user,omitempty"`
	Package string `form:"package" json:"package,omitempty"`
	Images int `form:"images" json:"images,omitempty"`
}


func Max(x int64, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

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
	Critical     null.Int    `json:"vulnerabilities_critical"`
	DefCon1      null.Int    `json:"vulnerabilities_defcon1"`
	High         null.Int    `json:"vulnerabilities_high"`
	Low          null.Int    `json:"vulnerabilities_low"`
	Medium       null.Int    `json:"vulnerabilities_medium"`
	Negligible   null.Int    `json:"vulnerabilities_negligible"`
	Unknown      null.Int    `json:"vulnerabilities_unknown"`
	Score        null.Int    `json:"value"`
	Analysed     null.Bool   `json:"analysed"`
	User         string      `json:"username"`
	Url          string      `json:"url"`
	OperatingSystem string	 `json:"operating_system"`
	Risk        string       `json:"risk"`
}

func getImageRepositorySQL(id int64, images *[]Image, limit int) {
	var image Image

	stmt, err := db.GetDB().Prepare("SELECT * FROM tag WHERE image_id=$1 and analysed ORDER BY SCORE DESC limit $2")
	rows, err := stmt.Query(id, limit)

	if err != nil {
		log.Print(err.Error())
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
			&image.OperatingSystem,
		)
		getOperatingSystem(&image)
		calculateRisk(&image)
		*images = append(*images, image)
		if err != nil {
			log.Print(err.Error())
		}
	}
}

func calculateRisk(image *Image){
	if image.Critical.ValueOrZero() > 0 {
		image.Risk = "critical"
	}else if image.High.ValueOrZero() > 0{
		image.Risk = "high"
	}else if image.Medium.ValueOrZero() > 0{
		image.Risk = "medium"
	} else if image.Low.ValueOrZero() > 0 || image.Unknown.ValueOrZero() > 0 || image.Negligible.ValueOrZero() > 0 {
		image.Risk = "low"
	} else {
		image.Risk = "none"
	}
}
/*
Get operating system
 */
func getOperatingSystem(image *Image){
	stmt, err := db.GetDB().Prepare(`select n.name FROM image as repo
	JOIN tag as image
	ON repo.id = image.image_id
	JOIN tag_layer as tl
	ON image.id = tl.tag_id
	JOIN layer as l
	ON tl.layer_id = l.id
	JOIN namespace as n
	ON n.id = l.namespace_id
	where image.id=$1 limit 1`)
	row := stmt.QueryRow(image.Id)


	if err != nil {
		log.Println(err.Error())
	}

	err = row.Scan(
		&image.OperatingSystem,
	)

	if err != nil {
		image.OperatingSystem = "unknown"
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

	getImageRepositorySQL(int64(id), &images, NUMBER_IMAGES)

	c.JSON(http.StatusOK, gin.H{
		"result": images,
		"count":  len(images),
	})
}

func FetchImagesVizPost(c *gin.Context) {
	var question Request
	var result RepositorySearchResult
	if err := c.ShouldBindJSON(&question); err == nil {
		repos := []Repository{}
		var best_image_score, best_image_size, maxSize int64
		best_image_score, best_image_size, maxSize = 0, 0, 0
		var numberImages int

		if question.Images > 0 {
			numberImages = question.Images
		} else {
			numberImages = 100
		}

		if question.Package != "" {
			getRepositoriesPattern(&repos, question.Package, true)
		} else if question.User != "" {
			getRepositoriesPattern(&repos, question.User, false)
		}

		for i := 0; i < len(repos); i++ {
			images := []Image{}
			getImageRepositorySQL(repos[i].Id.ValueOrZero(), &images, numberImages)
			for j := len(images) - 1; j >= 0; j-- {
				repos[i].Images = append(repos[i].Images, images[j])
				best_image_score = images[j].Score.Int64
				best_image_size = images[j].Full_size.Int64
				maxSize = Max(maxSize, images[j].Full_size.Int64)
			}
			repos[i].Score.SetValid(best_image_score)
			repos[i].Full_size.SetValid(best_image_size)
		}
		result.Repositories = repos
		result.Name = null.StringFrom(question.Package)

		c.JSON(http.StatusOK, gin.H{
			"count":  len(repos),
			"result": result,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func FetchImagesVizPostv2(c *gin.Context) {
	var question Request
	var result RepositorySearchResult
	var imageRepoIds []int
	var numberImages int
	if err := c.ShouldBindJSON(&question); err == nil {
		repos := []Repository{}

		if question.Images > 0 {
			numberImages = question.Images
		} else {
			numberImages = 100
		}


		getRepoImages(&repos, &imageRepoIds, question.User, numberImages)
		result.Repositories = repos

		c.JSON(http.StatusOK, gin.H{
			"count":  len(repos),
			"result": result,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func FetchImagesViz(c *gin.Context) {
	pattern := c.DefaultQuery("query", "")
	repos := []Repository{}
	var result RepositorySearchResult
	var best_image_score int64
	var best_image_size int64
	var maxSize int64

	getRepositoriesPattern(&repos, pattern, true)
	for i := 0; i < len(repos); i++ {
		images := []Image{}

		best_image_score = 0
		best_image_size = 0

		getImageRepositorySQL(repos[i].Id.ValueOrZero(), &images, NUMBER_IMAGES)
		for j := len(images) - 1; j >= 0; j-- {
			repos[i].Images = append(repos[i].Images, images[j])
			best_image_score = images[j].Score.Int64
			best_image_size = images[j].Full_size.Int64
			maxSize = Max(maxSize, images[j].Full_size.Int64)
		}
		repos[i].Score.SetValid(best_image_score)
		repos[i].Full_size.SetValid(best_image_size)
	}
	result.Repositories = repos
	result.Name = null.StringFrom(pattern)

	c.JSON(http.StatusOK, result)
}

func FetchImage(c *gin.Context) {
	var image Image
	id := c.Param("id")
	sqlStatement := `SELECT *, (select full_name from image where id=tag.image_id) as user
	FROM tag WHERE id=$1 LIMIT 1;`
	row := db.GetDB().QueryRow(sqlStatement, id)

	err := row.Scan(
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
		&image.User,
		&image.OperatingSystem,
	)
	image.Url = fmt.Sprintf("https://hub.docker.com/r/%s/%s", image.User, image.Name.String)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, image)
}

func vulnsBy(sev string, store map[string][]Vulnerability) []Vulnerability {
	items, found := store[sev]
	if !found {
		items = make([]Vulnerability, 0)
		store[sev] = items
	}
	return items
}


func groupBySeverity(vs []Vulnerability) map[string][]Vulnerability {
	var store = make(map[string][]Vulnerability)
	for _, v := range vs {
		sevRow := vulnsBy(v.Severity, store)
		store[v.Severity] = append(sevRow, v)
	}
	return store
}

func countBySeverity(vs []Vulnerability) map[string]int {
	var store = make(map[string][]Vulnerability)
	var count = make(map[string]int)
	for _, v := range vs {
		sevRow := vulnsBy(v.Severity, store)
		store[v.Severity] = append(sevRow, v)
	}

	for sev, vulns := range store{
		count[sev] = len(vulns)
	}
	return count
}



func FetchImagesVulns(c *gin.Context) {
	vulns := []Vulnerability{}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("err")
	}
	getImageVulnerabilitesSQL(int64(id), &vulns)
	groups := groupBySeverity(vulns)
	c.JSON(http.StatusOK, groups)
}


func FetchImagesPackages(c *gin.Context) {
	packages := []Feature{}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("err")
	}
	getImagePackagesSQL(int64(id), &packages)
	for i := 0; i < len(packages); i++{
		vulns := []Vulnerability{}
		getVulnerabilityPackage(packages[i].VersionId, &vulns)
		packages[i].Vulnerabilities = vulns
		packages[i].Summary = countBySeverity(vulns)
	}


	c.JSON(http.StatusOK, packages)
}
