package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
    var err error

    dsn := os.Getenv("DB_DSN") 
    if dsn == "" {
        log.Fatal("DB_DSN tidak ditemukan di environment variable")
    }

    DB, err = sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal("Gagal koneksi ke database:", err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatal("Gagal ping database:", err)
    }

    fmt.Println("Berhasil terhubung ke database PostgreSQL 🚀")
}
