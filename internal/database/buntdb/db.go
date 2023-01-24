package buntdb

import (
	"os"
	"path/filepath"
	"time"

	"github.com/tidwall/buntdb"
)

type Database struct {
	db *buntdb.DB
}

func NewDatabase(address string) *Database {
	if _, err := os.Stat(filepath.Dir(address)); os.IsNotExist(err) {
		if err = os.MkdirAll(filepath.Dir(address), 0755); err != nil {
			panic(err)
		}
	}
	db, err := buntdb.Open(address)
	if err != nil {
		panic(err)
	}
	return &Database{db}
}

func (d *Database) Get(server, key string) (string, bool) {
	var value string
	if err := d.db.View(func(tx *buntdb.Tx) error {
		var err error
		value, err = tx.Get(server + ":" + key)
		return err
	}); err != nil {
		return "", false
	}
	return value, true
}

func (d *Database) Set(server, key, value string, expire int64) error {
	return d.db.Update(func(tx *buntdb.Tx) error {
		opt := &buntdb.SetOptions{
			Expires: true,
			TTL:     time.Unix(expire, 0).Sub(time.Now()),
		}
		if expire < 0 {
			opt.Expires = false
			opt.TTL = 0
		}
		_, _, err := tx.Set(server+":"+key, value, opt)
		return err
	})
}

func (d *Database) Close() error {
	return d.db.Close()
}
