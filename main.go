package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/antchfx/xmlquery"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

var usage = "usage: go-coverage [coverage.xml]"

var dbHost = getEnv("ANALYTICS_DATABASE_HOST", "localhost")
var dbPort = getEnv("ANALYTICS_DATABASE_PORT", "5433")
var dbUser = getEnv("ANALYTICS_DATABASE_USERNAME", "postgres")
var dbPassword = getEnv("ANALYTICS_DATABASE_PASSWORD", "")
var dbName = getEnv("ANALYTICS_DATABASE_NAME", "analytics")

type CoverageReport struct {
	ID           uint   `gorm:"primaryKey"`
	Repository   string `gorm:"uniqueIndex:idx_component_version"`
	Component    string `gorm:"uniqueIndex:idx_component_version"`
	Version      string `gorm:"uniqueIndex:idx_component_version"`
	LineRate     string
	Timestamp    time.Time
	LinesCovered int
	LinesValid   int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal(usage)
	}
	report := parseCoverage(os.Args[1])
	log.Printf("coverage report: %v\n", report)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	autoMigrate(db)
	insertReport(db, report)
}

// XML report example
//
// <coverage line-rate="0.5242089" branch-rate="0" version="" timestamp="1669360021206" lines-covered="2750" lines-valid="5246" branches-covered="0" branches-valid="0" complexity="0">
func parseCoverage(filename string) *CoverageReport {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err.Error())
	}
	doc, err := xmlquery.Parse(f)
	if err != nil {
		panic(err)
	}
	c := xmlquery.FindOne(doc, "/coverage")

	repository := os.Getenv("DRONE_REPO")
	component := os.Getenv("ANALYTICS_COMPONENT")
	version := getVersion()

	if repository == "" {
		log.Fatal("Missing DRONE_REPO variable")
	}

	if component == "" {
		log.Fatal("Missing ANALYTICS_COMPONENT variable")
	}

	if version == "" {
		log.Fatal("Can't determine the version")
	}

	lineRateStr := c.SelectAttr("line-rate")
	timestampStr := c.SelectAttr("timestamp")
	timestampInt, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		log.Fatal("Failed to parse timestamp integer", err.Error())
	}
	timestamp := time.Unix(timestampInt/1000, timestampInt-timestampInt%1000)

	linesCoveredStr := c.SelectAttr("lines-covered")
	linesCovered, err := strconv.Atoi(linesCoveredStr)
	if err != nil {
		log.Fatal("Failed to parse lines-covered integer", err.Error())
	}

	linesValidStr := c.SelectAttr("lines-valid")
	linesValid, err := strconv.Atoi(linesValidStr)
	if err != nil {
		log.Fatal("Failed to parse lines-valid integer", err.Error())
	}

	return &CoverageReport{
		Repository:   repository,
		Component:    component,
		Version:      version,
		LineRate:     lineRateStr,
		Timestamp:    timestamp,
		LinesCovered: linesCovered,
		LinesValid:   linesValid,
	}
}

func getVersion() string {
	// 1.0.0
	v := os.Getenv("DRONE_TAG")
	if v != "" {
		return v
	}

	// bcdd4bf0245c82c060407b3b24b9b87301d15ac1
	v = os.Getenv("DRONE_COMMIT_SHA")
	if len(v) > 8 {
		return v[:8]
	}
	return ""
}

func insertReport(db *gorm.DB, report *CoverageReport) {
	tx := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "repository"}, {Name: "component"}, {Name: "version"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"line_rate",
			"timestamp",
			"lines_covered",
			"lines_valid",
			"updated_at",
		}),
	}).Create(report)

	if tx.Error != nil {
		log.Fatalln("Upsert failed", tx.Error.Error())
	}
}

func getEnv(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}

func autoMigrate(db *gorm.DB) {
	log.Println("Migrating database")
	err := db.AutoMigrate(
		&CoverageReport{},
	)
	if err != nil {
		log.Fatal(err.Error())
	}
}
