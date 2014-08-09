package main

import (
	"errors"
	"time"
)

type Storage interface {
	Store(content []byte, meta FileMeta) int
	Fetch(id int) ([]byte, FileMeta, error)
	FetchMeta(id int) (FileMeta, error)
	Delete(id int) (int, error)
}

type FileMeta struct {
	isPrivate   bool
	contentType string
	createdAt   time.Time
	password    string
}

func (this *FileMeta) IsAuth(password string) bool {
	return this.password == password
}

type Files struct {
	storage Storage
	expire  time.Duration
}

func (this *Files) isExpired(meta FileMeta, now time.Time) bool {
	return now.Equal(meta.createdAt.Add(this.expire)) || now.After(meta.createdAt.Add(this.expire))
}

func (this *Files) Download(id int, now time.Time) ([]byte, FileMeta, error) {
	content, meta, err := this.storage.Fetch(id)
	if this.isExpired(meta, now) {
		return nil, FileMeta{}, errors.New("expired")
	} else {
		return content, meta, err
	}
}

func (this *Files) Upload(content []byte, isPrivate bool, password string, contentType string, now time.Time) int {
	return this.storage.Store(content, FileMeta{isPrivate, contentType, now, password})
}

func (this *Files) Delete(id int, password string, now time.Time) (int, error) {
	meta, err := this.storage.FetchMeta(id)
	if err != nil {
		return -1, err
	} else if this.isExpired(meta, now) {
		return -1, errors.New("expired")
	} else if meta.IsAuth(password) {
		return this.storage.Delete(id)
	} else {
		return -1, errors.New("auth failed")
	}
}
