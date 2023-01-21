package memory

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var db sync.Map

type Database struct{}

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) Get(server, key string) (string, bool) {
	v, ok := db.Load(server + ":" + key)
	if !ok {
		return "", false
	}
	_v := v.(string)
	if _v[len(_v)-2:] == "-1" {
		return _v[:len(_v)-3], true
	}
	expire, err := strconv.ParseInt(_v[len(_v)-10:], 10, 64)
	if err != nil {
		return "", false
	}
	if expire < time.Now().Unix() {
		db.Delete(server + ":" + key)
		return "", false
	}
	return _v[:len(_v)-11], true
}

func (d *Database) Set(server, key, value string, expire int64) error {
	db.Store(server+":"+key, fmt.Sprintf("%s:%d", value, expire))
	return nil
}

func (d *Database) Close() error {
	return nil
}
