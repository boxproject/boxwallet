package util

type Pair struct {
	Key []byte
	Val []byte
}

type Database interface {
	Put(key, value []byte) error
	BatchPut(pairs []*Pair) error
	Get(key []byte) (value []byte, err error)
	Remove(key []byte) error
	Iterator(prefix []byte) (<-chan *Pair, error)
	Close() error
	SaveAndDelete(saveKey, deleteKey []byte, value []byte) error
}
