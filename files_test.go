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
	store map[int64]DummyStorageColumn
	id    int64
}

func (this *DummyStorage) Store(content []byte, meta FileMeta) (int64, error) {
	this.store[this.id] = DummyStorageColumn{content, meta}
	this.id++
	return this.id - 1, nil
}

func (this *DummyStorage) Fetch(id int64) ([]byte, FileMeta, error) {
	if stored, ok := this.store[id]; ok {
		return stored.content, stored.meta, nil
	} else {
		return nil, FileMeta{}, errors.New("Not found")
	}
}

func (this *DummyStorage) FetchMeta(id int64) (FileMeta, error) {
	if stored, ok := this.store[id]; ok {
		return stored.meta, nil
	} else {
		return FileMeta{}, errors.New("Not Found")
	}
}

func (this *DummyStorage) Delete(id int64) (int64, error) {
	if _, ok := this.store[id]; ok {
		delete(this.store, id)
		return id, nil
	} else {
		return -1, errors.New("Not Found")
	}
}

func (this *DummyStorage) List() ([]FileList, error) {
	list := make([]FileList, 0)
	for id, item := range this.store {
		list = append(list, FileList{id: id, createdAt: item.meta.createdAt})
	}
	return list, nil
}

func NewDummyStorage() *DummyStorage {
	dummyStorage := new(DummyStorage)
	dummyStorage.store = make(map[int64]DummyStorageColumn)
	dummyStorage.id = 0
	return dummyStorage
}

func TestDownload(t *testing.T) {
	dummyStorage := NewDummyStorage()
	now, _ := time.Parse(time.RFC3339, "2014-01-01T00:00:00Z09:00")
	id, _ := dummyStorage.Store([]byte("hoge"), FileMeta{createdAt: now})

	files := Files{storage: dummyStorage, expire: 1 * time.Minute}
	content, _, err := files.Download(id, now)
	assert.Nil(t, err)
	assert.Equal(t, string(content), "hoge")
}

func TestDownload2(t *testing.T) {
	dummyStorage := NewDummyStorage()
	now, _ := time.Parse(time.RFC3339, "2014-01-01T00:00:00Z09:00")

	id, _ := dummyStorage.Store([]byte("hoge"), FileMeta{createdAt: now})

	files := Files{storage: dummyStorage, expire: 1 * time.Minute}
	content, _, err := files.Download(id+1, now)
	assert.NotNil(t, err)
	assert.Nil(t, content)
}

func TestDownload3(t *testing.T) {
	dummyStorage := NewDummyStorage()
	now, _ := time.Parse(time.RFC3339, "2014-01-01T00:00:00Z09:00")

	id, _ := dummyStorage.Store([]byte("hoge"), FileMeta{createdAt: now})

	files := Files{storage: dummyStorage, expire: 1 * time.Minute}
	content, _, err := files.Download(id, now.Add(1*time.Minute))
	assert.NotNil(t, err)
	assert.Nil(t, content)
}

func TestUpload(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2014-01-01T00:00:00Z09:00")
	dummyStorage := NewDummyStorage()
	files := Files{storage: dummyStorage, expire: 1 * time.Minute}
	id, err := files.Upload([]byte("hoge"), false, "", "image/png", now)
	assert.Nil(t, err)
	content, _, _ := dummyStorage.Fetch(id)
	assert.Equal(t, content, []byte("hoge"))
}

func TestUpload2(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2014-01-01T00:00:00Z09:00")
	dummyStorage := NewDummyStorage()
	files := Files{storage: dummyStorage, expire: 1 * time.Minute}
	id1, err1 := files.Upload([]byte("hoge"), false, "", "image/png", now)
	id2, err2 := files.Upload([]byte("fuga"), false, "", "image/png", now)
	assert.NotEqual(t, id1, id2)
	assert.Nil(t, err1)
	assert.Nil(t, err2)
}

func TestDelete(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2014-01-01T00:00:00Z09:00")
	dummyStorage := NewDummyStorage()
	id, _ := dummyStorage.Store([]byte("hoge"), NewMeta(false, "image/png", now, "pass"))
	files := Files{storage: dummyStorage, expire: 1 * time.Minute}
	id2, err := files.Delete(id, "pass", now)
	assert.Equal(t, id, id2)
	content, _, err := files.Download(id, now)
	assert.NotNil(t, err)
	assert.Nil(t, content)
}

func TestDelete2(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2014-01-01T00:00:00Z09:00")
	dummyStorage := NewDummyStorage()
	id, _ := dummyStorage.Store([]byte("hoge"), NewMeta(false, "image/png", now, "pass"))
	files := Files{storage: dummyStorage, expire: 1 * time.Minute}
	_, err := files.Delete(id, "pa", now)
	assert.NotNil(t, err)
	content, _, err := files.Download(id, now)
	assert.Nil(t, err)
	assert.Equal(t, content, []byte("hoge"))
}

func TestDelete3(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2014-01-01T00:00:00Z09:00")
	dummyStorage := NewDummyStorage()
	id, _ := dummyStorage.Store([]byte("hoge"), NewMeta(false, "image/png", now, "pass"))
	files := Files{storage: dummyStorage, expire: 1 * time.Minute}
	_, err := files.Delete(id+1, "pass", now)
	assert.NotNil(t, err)
	content, _, err := files.Download(id, now)
	assert.Nil(t, err)
	assert.Equal(t, content, []byte("hoge"))
}

func TestDelete4(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2014-01-01T00:00:00Z09:00")
	dummyStorage := NewDummyStorage()
	id, _ := dummyStorage.Store([]byte("hoge"), NewMeta(false, "image/png", now, "pass"))
	files := Files{storage: dummyStorage, expire: 1 * time.Minute}
	_, err := files.Delete(id, "pass", now.Add(1*time.Minute))
	assert.NotNil(t, err)
}

func TestList(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2014-01-01T00:00:00Z09:00")
	dummyStorage := NewDummyStorage()
	files := Files{storage: dummyStorage, expire: 1 * time.Minute}
	list, _ := files.List(now)
	assert.Equal(t, len(list), 0)
}

func TestList2(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2014-01-01T00:00:00Z09:00")
	dummyStorage := NewDummyStorage()
	id, _ := dummyStorage.Store([]byte("hoge"), NewMeta(false, "image/png", now, "pass"))
	files := Files{storage: dummyStorage, expire: 1 * time.Minute}
	list, _ := files.List(now)
	assert.Equal(t, len(list), 1)
	assert.Equal(t, list[0].id, id)
}

func TestList3(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2014-01-01T00:00:00Z09:00")
	dummyStorage := NewDummyStorage()
	_, _ = dummyStorage.Store([]byte("hoge"), NewMeta(false, "image/png", now, "pass"))
	files := Files{storage: dummyStorage, expire: 1 * time.Minute}
	list, _ := files.List(now.Add(1 * time.Minute))
	assert.Equal(t, len(list), 0)
}
