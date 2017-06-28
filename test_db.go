package main

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
)

func main() {
	db, err := bolt.Open("/Users/mylxsw/Downloads/test.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	db.Update(func (tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("test"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		b.Put([]byte("answer"), []byte("what are you doing"))

		return nil
	})

	db.View(func (tx *bolt.Tx) error {
		b := tx.Bucket([]byte("test"))
		fmt.Printf("%s", b.Get([]byte("answer")))

		return nil
	})

}
