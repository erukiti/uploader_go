package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type DummyStorageColumn struct {
	content []byte
	meta    FileMeta
}

type DummyStorage struct {
	store map[int]DummyStorageColumn
	id    int
}

func (this *DummyStorage) Store(content []byte, meta FileMeta) int {
	this.store[this.id] = DummyStorageColumn{content, meta}
	this.id++
	return this.id - 1
}

func (this *DummyStorage) Fetch(id int) ([]byte, FileMeta, error) {
	if stored, ok := this.store[id]; ok {
		return stored.content, stored.meta, nil
	} else {
		return nil, FileMeta{}, errors.New("Not found")
	}
}

func NewDummyStorage() *DummyStorage {
	dummyStorage := new(DummyStorage)
	dummyStorage.store = make(map[int]DummyStorageColumn)
	dummyStorage.id = 0
	return dummyStorage
}

func TestDownload(t *testing.T) {
	dummyStorage := NewDummyStorage()
	now, _ := time.Parse(time.RFC3339, "2014-01-01T00:00:00Z09:00")
	id := dummyStorage.Store([]byte("hoge"), FileMeta{createdAt: now})

	files := Files{storage: dummyStorage, expire: 1 * time.Minute}
	content, _, err := files.Download(id, now)
	assert.Nil(t, err)
	assert.Equal(t, string(content), "hoge")
}

func TestDownload2(t *testing.T) {
	dummyStorage := NewDummyStorage()
	now, _ := time.Parse(time.RFC3339, "2014-01-01T00:00:00Z09:00")

	id := dummyStorage.Store([]byte("hoge"), FileMeta{createdAt: now})

	files := Files{storage: dummyStorage, expire: 1 * time.Minute}
	content, _, err := files.Download(id+1, now)
	assert.NotNil(t, err)
	assert.Nil(t, content)
}

func TestDownload3(t *testing.T) {
	dummyStorage := NewDummyStorage()
	now, _ := time.Parse(time.RFC3339, "2014-01-01T00:00:00Z09:00")

	id := dummyStorage.Store([]byte("hoge"), FileMeta{createdAt: now})

	files := Files{storage: dummyStorage, expire: 1 * time.Minute}
	content, _, err := files.Download(id, now.Add(1*time.Minute))
	assert.NotNil(t, err)
	assert.Nil(t, content)
}

func TestUpload(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2014-01-01T00:00:00Z09:00")
	dummyStorage := NewDummyStorage()
	files := Files{storage: dummyStorage, expire: 1 * time.Minute}
	id := files.Upload([]byte("hoge"), false, "", "image/png", now)
	content, _, _ := dummyStorage.Fetch(id)
	assert.Equal(t, content, []byte("hoge"))
}

func TestUpload2(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2014-01-01T00:00:00Z09:00")
	dummyStorage := NewDummyStorage()
	files := Files{storage: dummyStorage, expire: 1 * time.Minute}
	id1 := files.Upload([]byte("hoge"), false, "", "image/png", now)
	id2 := files.Upload([]byte("fuga"), false, "", "image/png", now)
	assert.NotEqual(t, id1, id2)
}
