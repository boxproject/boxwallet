package db

import (
	"time"

	"fmt"

	"github.com/boxproject/boxwallet/util"
	"github.com/dgraph-io/badger"
	"github.com/sirupsen/logrus"
)

var (
	defaultStore util.Database
)

type internalDb struct {
	db *badger.DB
}

func GetStore() util.Database {
	return defaultStore
}

func (idb *internalDb) SaveAndDelete(saveKey, deleteKey []byte, value []byte) (err error) {
	txn := idb.db.NewTransaction(true)
	if err = txn.Set(saveKey, value); err != nil {
		txn.Discard()
		return
	}
	if err = txn.Delete(deleteKey); err != nil {
		txn.Discard()
		return
	}

	return txn.Commit()
}

func (idb *internalDb) Put(key, value []byte) error {
	return idb.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (idb *internalDb) BatchPut(pairs []*util.Pair) error {
	return idb.db.Update(func(txn *badger.Txn) error {
		txn = idb.db.NewTransaction(true)
		for _, pair := range pairs {
			if err := txn.Set(pair.Key, pair.Val); err == badger.ErrTxnTooBig {
				if err = txn.Commit(); err != nil {
					return err
				}
				txn = idb.db.NewTransaction(true)
				if err = txn.Set(pair.Key, pair.Val); err != nil {
					return err
				}
			}
		}
		return txn.Commit()
	})
}

func (idb *internalDb) Remove(key []byte) error {
	return idb.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func (idb *internalDb) Get(key []byte) (value []byte, err error) {
	err = idb.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		value, err = item.ValueCopy(value)
		return err
	})

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (idb *internalDb) Iterator(prefix []byte) (<-chan *util.Pair, error) {
	data := make(chan *util.Pair, 127)
	err := idb.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			val := new([]byte)
			err := item.Value(func(v []byte) error {
				fmt.Printf("key=%s, value=%s\n", k, v)
				val = &v
				return nil
			})

			if err != nil {
				close(data)
				return err
			}

			data <- &util.Pair{k, *val}
		}
		close(data)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (idb *internalDb) Close() error {
	if idb.db != nil {
		return idb.db.Close()
	}
	return nil
}

func Open(filePath string) error {
	opts := badger.DefaultOptions
	opts.Dir = filePath
	opts.ValueDir = filePath
	db, err := badger.Open(opts)
	if err != nil {
		logrus.WithError(err).Errorf("打开数据库失败")
		return err
	}

	defaultStore = &internalDb{
		db: db,
	}

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
		again:
			err := db.RunValueLogGC(0.7)
			if err == nil {
				goto again
			}
		}
	}()

	return nil
}

func Close() error {
	if defaultStore != nil {
		return defaultStore.Close()
	}

	return nil
}

func Put(key, value []byte) error {
	return defaultStore.Put(key, value)
}

func BatchPut(pairs []*util.Pair) error {
	return defaultStore.BatchPut(pairs)
}

func Get(key []byte) ([]byte, error) {
	return defaultStore.Get(key)
}

func Iterator(prefix []byte) (<-chan *util.Pair, error) {
	return defaultStore.Iterator(prefix)
}

func Remove(key []byte) error {
	return defaultStore.Remove(key)
}
