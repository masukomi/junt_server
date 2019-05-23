package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"masukomi.org/junt/models"
	// if this wasn't a go module
	// we would import with "./models"
	"masukomi.org/junt/controllers"
)

func main() {
	ensureConfigDir()
	i := Impl{}
	i.InitDB()
	defer i.DB.Close()
	i.InitSchema()
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	cc := controllers.CompaniesController{i.DB}
	pc := controllers.PeopleController{i.DB}
	fc := controllers.FollowupsController{i.DB}
	hc := controllers.HomeworksController{i.DB}
	ic := controllers.InterviewsController{i.DB}
	jc := controllers.JobsController{i.DB}
	oc := controllers.OffersController{i.DB}
	router, err := rest.MakeRouter(
		rest.Get("/companies", cc.ListAll),
		rest.Get("/companies/:id", cc.FindById),
		rest.Post("/companies", cc.Create),

		rest.Get("/followups", fc.ListAll),
		rest.Get("/followups/:id", fc.FindById),
		rest.Post("/followups", fc.Create),

		rest.Get("/homeworks", hc.ListAll),
		rest.Get("/homeworks/:id", hc.FindById),
		rest.Post("/homeworks", hc.Create),

		rest.Get("/interviews", ic.ListAll),
		rest.Get("/interviews/:id", ic.FindById),
		rest.Post("/interviews", ic.Create),

		rest.Get("/jobs", jc.ListAll),
		rest.Get("/jobs/:id", jc.FindById),
		rest.Post("/jobs", jc.Create),

		rest.Get("/offers", oc.ListAll),
		rest.Get("/offers/:id", oc.FindById),
		rest.Post("/offers", oc.Create),

		rest.Get("/people", pc.ListAll),
		rest.Get("/people/:id", pc.FindById),
		rest.Post("/people", pc.Create),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8123", api.MakeHandler()))
}

func configDirPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "junt")
}
func ensureConfigDir() {
	path := configDirPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
}

///// Internal stuff
type Impl struct {
	DB *gorm.DB
}

func (i *Impl) InitDB() {
	var err error
	dbPath := filepath.Join(configDirPath(), "junt.db")
	i.DB, err = gorm.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Got error when connecting to database: '%v'", err)
	}
	i.DB.LogMode(true)
}

func (i *Impl) InitSchema() {
	// WARNING: will only create tables, missing columns, and missing indexes.
	// will NOT change existing column's type or delete unused columns
	i.DB.AutoMigrate(
		&models.Company{},
		&models.Followup{},
		&models.Homework{},
		&models.Interview{},
		&models.Job{},
		&models.Offer{},
		&models.Person{},
		&models.StatusChange{},
		&models.ThanksEmail{},
	)
}
