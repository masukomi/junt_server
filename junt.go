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
	"fmt"
	"masukomi.org/junt/controllers"
	"strconv"
)

const DEFAULT_PORT = 8123

func main() {
	ensureConfigDir()
	i := Impl{}
	i.InitDB()
	defer i.DB.Close()
	i.InitSchema()
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	// api.Use(&rest.CorsMiddleware{
	//     RejectNonCorsRequests: false,
	//     OriginValidator: func(origin string, request *rest.Request) bool {
	//         return origin == "http://my.other.host"
	//     },
	//     AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	//     AllowedHeaders: []string{
	//         "Accept", "Content-Type", "X-Custom-Header", "Origin"},
	//     AccessControlAllowCredentials: true,
	//     AccessControlMaxAge:           3600,
	// })
	crudThings := map[string]controllers.CrudController{
		"companies":      &controllers.CompaniesController{i.DB},
		"people":         &controllers.PeopleController{i.DB},
		"followups":      &controllers.FollowupsController{i.DB},
		"homeworks":      &controllers.HomeworksController{i.DB},
		"interviews":     &controllers.InterviewsController{i.DB},
		"jobs":           &controllers.JobsController{i.DB},
		"offers":         &controllers.OffersController{i.DB},
		"status_changes": &controllers.StatusChangesController{i.DB},
		"thanks":         &controllers.ThanksEmailsController{i.DB},
	}
	ec := controllers.EventsController{i.DB}
	makeRouterArgs := []*rest.Route{}
	for name, controller := range crudThings {
		makeRouterArgs = append(makeRouterArgs, rest.Get("/"+name, controller.ListAll))
		makeRouterArgs = append(makeRouterArgs, rest.Get("/"+name+"/:id", controller.FindById))
		makeRouterArgs = append(makeRouterArgs, rest.Post("/"+name, controller.Create))
		// TODO IMPLEMENT UPDATE
		// makeRouterArgs = append(makeRouterArgs, rest.Delete("/"+name, controller.Update))
	}

	makeRouterArgs = append(makeRouterArgs, rest.Get("/events", ec.ListAll))
	makeRouterArgs = append(makeRouterArgs, rest.Get("/events/job/:id", ec.ListAllForJob))

	router, err := rest.MakeRouter(makeRouterArgs...)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	fmt.Println("Listening at http://localhost:" + strconv.Itoa(DEFAULT_PORT))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(DEFAULT_PORT), api.MakeHandler()))

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
		&models.Job{},
		&models.Person{},
		// events ....
		&models.Followup{},
		&models.Homework{},
		&models.Interview{},
		&models.Offer{},
		&models.StatusChange{},
		&models.ThanksEmail{},
	)
}
