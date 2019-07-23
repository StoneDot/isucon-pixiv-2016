package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"strconv"
)

type Post struct {
	ID      int    `db:"id"`
	Imgdata []byte `db:"imgdata"`
	Mime    string `db:"mime"`
}

func main() {
	stderr := log.New(os.Stderr, "", 0)
	host := os.Getenv("ISUCONP_DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("ISUCONP_DB_PORT")
	if port == "" {
		port = "3306"
	}
	_, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Failed to read DB port number from an environment variable ISUCONP_DB_PORT.\nError: %s", err.Error())
	}
	user := os.Getenv("ISUCONP_DB_USER")
	if user == "" {
		user = "root"
	}
	password := os.Getenv("ISUCONP_DB_PASSWORD")
	dbname := os.Getenv("ISUCONP_DB_NAME")
	if dbname == "" {
		dbname = "isuconp"
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		user,
		password,
		host,
		port,
		dbname,
	)
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	posts := []Post{}
	err = db.Select(&posts, "SELECT * FROM `posts`")
	if err != nil {
		panic(err)
	}

	dir := "images/"
	err = os.MkdirAll(dir, 775)
	for _, p := range posts {
		ext := ""
		switch p.Mime {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		}
		filename := fmt.Sprintf("%s.%s", strconv.Itoa(p.ID), ext)
		path := dir + filename
		file, err := os.Open(path)
		if err != nil {
			stderr.Println("Failed to open file: " + filename)
			continue
		}
		_, err = file.Write(p.Imgdata)
		if err != nil {
			stderr.Println("Failed to write file: " + filename)
		}
	}
}
