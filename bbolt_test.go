package bbrowse

import (
	"bytes"
	"encoding/gob"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
)

type a struct {
	B string
	C int
}

func TestCopyBoltDB_Success(t *testing.T) {
	boltDB, err := createTestDB(t)
	assert.NoError(t, err)

	_, err = copyBoltDB(boltDB)

	assert.NoError(t, err)
}

func createTestDB(t *testing.T) (*bolt.DB, error) {
	t.Helper()

	f, err := os.CreateTemp("", "*.db")
	if err != nil {
		return nil, err
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	f.Close()

	db, err := bolt.Open(f.Name(), 0600, nil)
	if err != nil {
		return nil, err
	}
	t.Cleanup(func() { db.Close() })

	if err := db.Update(func(tx *bolt.Tx) error {
		// Put a bucket in the root
		bucket, err := tx.CreateBucket([]byte("bucket"))
		if err != nil {
			return err
		}

		// Add a key/value in the bucket
		err = bucket.Put([]byte("key"), []byte("value"))
		if err != nil {
			return err
		}

		err = bucket.Put([]byte("gob"), createGOB())
		if err != nil {
			return err
		}

		// Add a nested bucket
		nestedBucket, err := bucket.CreateBucket([]byte("nested bucket"))
		if err != nil {
			return err
		}

		// Add a key/value in the nested bucket
		err = nestedBucket.Put([]byte("key"), []byte("value"))
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}

func createGOB() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(a{B: "hello", C: 42})
	return buf.Bytes()
}
