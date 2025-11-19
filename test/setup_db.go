package test

import (
	"log"
	"path/filepath"
	"runtime"

	"backendgo/database"
	"github.com/joho/godotenv"
)

func init() {

	// detect root project path
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Join(filepath.Dir(b), "..", ".env")

	// load .env
	err := godotenv.Load(basepath)
	if err != nil {
		log.Println("Warning: gagal load .env dari setup_db.go")
	}

	// connect PostgreSQL
	database.ConnectDB()

	// connect MongoDB
	database.ConnectMongoDB()
}
