package db

import (
	"go_crypo_coin/utils"

	"github.com/boltdb/bolt"
)

const (
	dbName      = "blockchain.db"
	dataBucket  = "data"
	blockBucket = "blocks"
	checkpoint  = "checkpoint"
)

var db *bolt.DB

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open("blockchain.db", 0600, nil)
		utils.HandleErr(err)
		db = dbPointer
		err = db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleErr(err)
			_, err = t.CreateBucketIfNotExists([]byte(blockBucket))
			return err
		})
		utils.HandleErr(err)
	}
	return db
}

func Close() {
	DB().Close()
}

func SaveBlock(hash string, data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blockBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

func SaveCheckPoint(data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandleErr(err)
}

func Checkpoint() []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}

func Block(hash string) []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blockBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}
