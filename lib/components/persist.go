package components

import (
	"strings"

	"github.com/x-hgg-x/arkanoid-go/lib/binconv"
	bolt "go.etcd.io/bbolt"
)

// Persist keeps the persistence overtime
type Persist struct {
	bucket []byte
	db     *bolt.DB
}

// NewPersist gives new instance of persistence
func NewPersist(dbName, bucket string) (*Persist, error) {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}

	b := []byte(bucket)

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(b)
		return err
	})

	return &Persist{b, db}, nil
}

// Close implements closer
func (p *Persist) Close() {
	p.db.Close()
}

// Update adds data(value) using key
func (p *Persist) Update(key, value []byte) error {
	return p.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(p.bucket)
		err := b.Put(key, value)
		return err
	})
}

// View gets data(value) using key
func (p *Persist) View(key []byte) (value []byte, err error) {
	err = p.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(p.bucket)
		value = b.Get(key)
		return nil
	})
	return
}

// UpdateList adds/update array(value) using key
func (p *Persist) UpdateList(key []byte, value [][]byte) error {
	return p.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(p.bucket)

		for i, v := range value {
			key = append(key, binconv.Itob(i)...)
			err := b.Put(key, v)
			if err != nil {
				return err // must return slice even if almost finish
			}
		}
		return nil
	})
}

// ViewList gets array(value) using key
func (p *Persist) ViewList(key []byte) (value [][]byte, err error) {
	err = p.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(p.bucket)
		b.ForEach(func(k, v []byte) error {
			if strings.HasPrefix(string(k), string(key)) {
				value = append(value, v)
			}
			return nil
		})
		return nil
	})
	return
}