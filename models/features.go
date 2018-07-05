package models

import "fmt"
import (
	"github.com/dockerpedia/api/db"
)

type Feature struct {
	Id              int             `json:"id,omitempty"`
	VersionId       int              `json:"version_id,omitempty"`
	Name            string          `json:"name,omitempty"`
	NamespaceName   string          `json:"namespace_name,omitempty"`
	Version         string          `json:"version,omitempty"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	AddedBy         string          `json:"added_by,omitempty"`
	Summary         map[string]int  `json:"summary,omitempty"`
}


func getImagePackagesSQL(image_id int64, features *[]Feature) {
	var feature Feature

	stmt, err := db.GetDB().Prepare(`
	select DISTINCT l.name, f.name, fv.version, f.id, fv.id FROM
      image as repo
      JOIN tag as image
          ON repo.id = image.image_id
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
      where image.id=$1;
    `)
	rows, err := stmt.Query(image_id)

	if err != nil {
		fmt.Print(err.Error())
	}

	for rows.Next() {
		err := rows.Scan(
			&feature.AddedBy,
			&feature.Name,
			&feature.Version,
			&feature.Id,
			&feature.VersionId,
		)
		*features = append(*features, feature)
		if err != nil {
			fmt.Print(err.Error())
		}
	}
}


func getVulnerabilityPackage(featureVersion int,  vulns *[]Vulnerability) {
	var vuln Vulnerability

	stmt, err := db.GetDB().Prepare(`
	select  v.id, v.name, v.description, v.link, v.severity, v.metadata FROM
      featureversion as fv
      JOIN vulnerability_affects_featureversion as vf
        ON vf.featureversion_id = fv.id
      JOIN vulnerability as v
        ON v.id = vf.vulnerability_id
      where fv.id=$1;
    `)
	rows, err := stmt.Query(featureVersion)

	if err != nil {
		fmt.Print(err.Error())
	}

	for rows.Next() {
		err := rows.Scan(
			&vuln.Id,
			&vuln.Name,
			&vuln.Description,
			&vuln.Link,
			&vuln.Severity,
			&vuln.Metadata,
		)
		*vulns = append(*vulns, vuln)
		if err != nil {
			fmt.Print(err.Error())
		}
	}

}