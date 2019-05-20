package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"./models"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	ensureConfigDir()
	i := Impl{}
	i.InitDB()
	defer i.DB.Close()
	i.InitSchema()
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	// router, err := rest.MakeRouter(
	// 	rest.Get("/companies",
	// )

	// api.SetApp(rest.AppSimple(func(w rest.ResponseWriter,
	// 	r *rest.Request) {
	// 	w.WriteJson(map[string]string{"Body": "Hello World!"})
	// }))
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
	i.DB.AutoMigrate(&models.Company{})
}
