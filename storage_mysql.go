package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

type storage struct {
	db *sql.DB
}

func NewStorage() (storage, error) {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	scheme := os.Getenv("MYSQL_DB")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, scheme)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return storage{}, err
	} else {
		return storage{db}, nil
	}
}

func (this *storage) Store(content []byte, meta FileMeta) (int64, error) {
	stmt, err := this.db.Prepare("insert into storage(created_at, is_private, content_type, salt, hashed_password, content) values (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return -1, err
	}
	var f int
	if meta.isPrivate {
		f = 1
	} else {
		f = 0
	}
	res, err := stmt.Exec(meta.createdAt, f, meta.contentType, meta.salt, meta.hashedPassword, content)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

func (this *storage) Fetch(id int64) ([]byte, FileMeta, error) {
	var err error
	stmt, err := this.db.Prepare("select created_at, is_private, content_type, salt, hashed_password, content from storage where id = ?")
	if err != nil {
		return []byte(""), FileMeta{}, err
	}
	var content []byte
	var createdAt time.Time
	var isPrivate int
	var contentType string
	var salt string
	var hashedPassword string
	err = stmt.QueryRow(id).Scan(createdAt, isPrivate, contentType, salt, hashedPassword, content)
	if err != nil {
		return []byte(""), FileMeta{}, err
	}
	return content, NewMetaHashedPassword(isPrivate == 1, contentType, createdAt, hashedPassword, salt), nil
}

func (this *storage) FetchMeta(id int64) (FileMeta, error) {
	return FileMeta{}, nil
}

func (this *storage) Delete(id int64) (int64, error) {
	return -1, nil
}
