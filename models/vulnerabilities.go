package models

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dockerpedia/api/db"
	"github.com/gin-gonic/gin"
)

type Vulnerability struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	NamespaceName string `json:"namespace_name,omitempty"`
	Description   string `json:"description,omitempty"`
	Link          string `json:"link,omitempty"`
	Severity      string `json:"severity,omitempty"`
	Metadata      string `json:"metadata,omitempty"`
	// fixed_by exists when vulnerability is under feature.
	FixedBy string `json:"fixed_by,omitempty"`
	// affected_versions exists when vulnerability is under notification.
	AffectedVersions []*Feature `json:"affected_versions"`
}

func getImageVulnerabilitesSQL(image_id int64, vulns *[]Vulnerability) {
	var vuln Vulnerability

	stmt, err := db.GetDB().Prepare(`select
	 DISTINCT v.name, v.id, v.description,
    v.link, v.severity, v.metadata
    FROM
    tag as image
    JOIN tag_layer as tl
      ON image.id = tl.tag_id
    JOIN layer as l
      ON tl.layer_id = l.id
    JOIN layer_diff_featureversion as ld
      ON ld.layer_id = l.id
    JOIN featureversion as fv
      ON ld.featureversion_id = fv.id
    JOIN feature as f
      ON f.id = fv.feature_id
    JOIN vulnerability_affects_featureversion as vf
      ON fv.id = vf.featureversion_id
    JOIN vulnerability as v
      ON v.id = vf.vulnerability_id
    where image.id=$1;
    `)
	rows, err := stmt.Query(image_id)

	if err != nil {
		fmt.Print(err.Error())
	}

	for rows.Next() {
		err := rows.Scan(
			&vuln.Name,
			&vuln.Id,
			&vuln.Description,
			&vuln.Link,
			&vuln.Severity,
			&vuln.Metadata,
		)

		//vuln.FixedBy = getFixedVulnerabilitySQL(vuln.Id, featureVulnerableId)
		*vulns = append(*vulns, vuln)
		if err != nil {
			fmt.Print(err.Error())
		}
	}
}

func FetchVulnerability(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}

	vuln := getVulnerabilitySQL(int64(id))
	c.JSON(http.StatusOK, vuln)
}

func getVulnerabilitySQL(id int64) (vuln Vulnerability) {
	sqlStatement := `SELECT name,description,
	   link,severity,metadata FROM vulnerability WHERE id=$1 LIMIT 1;`
	row := db.GetDB().QueryRow(sqlStatement, id)

	err := row.Scan(
		&vuln.Name,
		&vuln.Description,
		&vuln.Link,
		&vuln.Severity,
		&vuln.Metadata,
	)
	if err != nil {
		log.Fatal(err)
	}
	return vuln
}

func getFixedVulnerabilitySQL(vulnerability_id, feature_id int64) (featureVersion string) {
	sqlStatement := `SELECT version WHERE vulnerability_id=$1 and feature_id=$2 LIMIT 1;`
	row := db.GetDB().QueryRow(sqlStatement, vulnerability_id, feature_id)
	err := row.Scan(featureVersion)
	if err != nil {
		log.Fatal(err)
	}
	return
}
