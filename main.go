package main

import (
    "database/sql"
    "fmt"
    "os"
    "time"
    "net/http"


    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
)


var db *sql.DB

const (
    dbhost = "DBHOST"
    dbport = "DBPORT"
    dbuser = "DBUSER"
    dbpass = "DBPASS"
    dbname = "DBNAME"
)


type Repo struct {
    

    Id              int
    Name            string
    Full_name       string
    Namespace       string
    User            string
    Affiliation     string
    Description     string
    Is_automated    bool
    Last_updated    time.Time
    Pull_count      int
    Repository_type string
    Star_count      int
    Status          bool
    Tags_checked    time.Time
    Official        bool

}



func main() {
    initDb()
    defer db.Close()


    router := gin.Default()

    v1 := router.Group("/api/v1/images")
    {
        v1.GET("/", fetchAllImage)
        v1.GET("/test", fetchAllImagePattern)
        //v1.GET("/:id", fetchAllImage)

    }

    router.Run(":8080")

}

func fetchAllImagePattern(c *gin.Context){
    pattern := c.DefaultQuery("query", "mysql")
    fmt.Println(pattern)
    var (
        repo  Repo
        repos []Repo
    )

    stmt, err := db.Prepare("SELECT * FROM image WHERE name like '%' || $1 || '%' ORDER BY pull_count DESC limit 10")
    rows, err := stmt.Query(pattern)


    if err != nil {
        fmt.Print(err.Error())
    }

    for rows.Next() {
        err = rows.Scan(
            & repo.Id,
            & repo.Name, 
            & repo.Full_name,
            & repo.Namespace,
            & repo.User, 
            & repo.Affiliation,
            & repo.Description,
            & repo.Is_automated,
            & repo.Last_updated,
            & repo.Pull_count,
            & repo.Repository_type,
            & repo.Star_count, 
            & repo.Status,
            & repo.Tags_checked,
            & repo.Official )

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

// fetchAllTodo fetch all todos
func fetchAllImage(c *gin.Context) {
    var (
        repo  Repo
        repos []Repo
    )

    rows, err := db.Query(`
        SELECT *  FROM image limit 10`)

    if err != nil {
        fmt.Print(err.Error())
    }

    for rows.Next() {
        err = rows.Scan(
            & repo.Id,
            & repo.Name, 
            & repo.Full_name,
            & repo.Namespace,
            & repo.User, 
            & repo.Affiliation,
            & repo.Description,
            & repo.Is_automated,
            & repo.Last_updated,
            & repo.Pull_count,
            & repo.Repository_type,
            & repo.Star_count, 
            & repo.Status,
            & repo.Tags_checked,
            & repo.Official )

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



func initDb() {
    config := dbConfig()
    var err error
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
        "password=%s dbname=%s sslmode=disable",
        config[dbhost], config[dbport],
        config[dbuser], config[dbpass], config[dbname])

    db, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        panic(err)
    }
    err = db.Ping()
    if err != nil {
        panic(err)
    }
    fmt.Println("Successfully connected!")
}


func dbConfig() map[string]string {
    conf := make(map[string]string)
    host, ok := os.LookupEnv(dbhost)
    if !ok {
        panic("DBHOST environment variable required but not set")
    }
    port, ok := os.LookupEnv(dbport)
    if !ok {
        panic("DBPORT environment variable required but not set")
    }
    user, ok := os.LookupEnv(dbuser)
    if !ok {
        panic("DBUSER environment variable required but not set")
    }
    password, ok := os.LookupEnv(dbpass)
    if !ok {
        panic("DBPASS environment variable required but not set")
    }
    name, ok := os.LookupEnv(dbname)
    if !ok {
        panic("DBNAME environment variable required but not set")
    }
    conf[dbhost] = host
    conf[dbport] = port
    conf[dbuser] = user
    conf[dbpass] = password
    conf[dbname] = name
    return conf
}
