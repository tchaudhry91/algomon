package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/tchaudhry91/algoprom/actions"
	"github.com/tchaudhry91/algoprom/algochecks"
)

type BoltStore struct {
	db     *bolt.DB
	logger *log.Logger
}

func NewBoltStore(f string, logger *log.Logger) (*BoltStore, error) {
	db, err := bolt.Open(f, 0600, nil)
	if err != nil {
		return nil, err
	}
	s := &BoltStore{
		db:     db,
		logger: logger,
	}
	return s, nil
}

func (s *BoltStore) PutCheck(ctx context.Context, check *algochecks.Check, output *algochecks.Output) (key string, err error) {
	key = ""
	name := check.Name
	err = s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("checks"))
		if err != nil {
			return err
		}
		bucket, err = bucket.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}
		key = strconv.FormatInt(output.Timestamp.Unix(), 10)
		val, err := json.Marshal(output)
		if err != nil {
			return fmt.Errorf("Error Marshalling Output to JSON: %v", err)
		}
		bucket.Put([]byte(key), val)
		return nil
	})
	return key, err
}

func (s *BoltStore) PutAction(ctx context.Context, checkName string, action *actions.ActionMeta, output *actions.Output) (key string, err error) {
	key = checkName + "_" + action.Name + "_"
	err = s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("actions"))
		if err != nil {
			return err
		}
		key += strconv.FormatInt(output.Timestamp.Unix(), 10)
		val, err := json.Marshal(output)
		if err != nil {
			return fmt.Errorf("Error Marshalling Output to JSON: %v", err)
		}
		bucket.Put([]byte(key), val)
		return nil
	})
	return key, err
}
