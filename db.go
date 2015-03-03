package tatslack

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

// DB stores the slack data.
type DB struct {
	db *bolt.DB
}

// Open opens the underlying database.
func Open(path string) (*DB, error) {
	// Open the underlying database.
	d, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	// Initialize buckets.
	db := &DB{db: d}
	if err := db.db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("messages"))
		return nil
	}); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// Close closes the underlying database.
func (db *DB) Close() error {
	return db.db.Close()
}

// Message returns a list of messages for a channel.
func (db *DB) Messages(channel string) ([]*Message, error) {
	var a []*Message
	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("messages")).Bucket([]byte(channel))
		if b == nil {
			return nil
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			m := &Message{}
			if err := json.Unmarshal(v, &m); err != nil {
				return err
			}
			a = append(a, m)
		}

		return nil
	})
	return a, err
}

// SaveMessages persists a list of messages for a channel.
func (db *DB) SaveMessages(channel string, a []*Message) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		// Create bucket for channel.
		b, err := tx.Bucket([]byte("messages")).CreateBucketIfNotExists([]byte(channel))
		if err != nil {
			return err
		}

		// Iterate over messages and insert.
		for _, m := range a {
			// Encode message into bytes.
			buf, err := json.Marshal(m)
			if err != nil {
				return err
			}

			// Save message by timestamp.
			if err := b.Put([]byte(m.TS), buf); err != nil {
				return err
			}
		}

		return nil
	})
}
