package bbrowse

import (
	"errors"
	"os"

	bolt "go.etcd.io/bbolt"
)

func OpenAndCopyBoltDB(filename string) (*Bucket, error) {
	if _, err := os.Stat(filename); err != nil {
		return nil, err
	}
	db, err := bolt.Open(filename, 0666, &bolt.Options{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return copyBoltDB(db)
}

func copyBoltDB(boltDB *bolt.DB) (*Bucket, error) {
	root := Bucket{Name: []byte("root")}
	err := boltDB.View(func(tx *bolt.Tx) error {
		if err := tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			bucket, err := copyBoltBucket(b)
			if err != nil {
				return err
			}
			bucket.Name = name
			root.Buckets = append(root.Buckets, bucket)
			return nil
		}); err != nil {
			return err
		}
		return nil
	})
	return &root, err
}

func copyBoltBucket(b *bolt.Bucket) (*Bucket, error) {
	root := Bucket{}
	queue := []*bolt.Bucket{b}
	copied := map[*bolt.Bucket]*Bucket{b: &root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if err := node.ForEach(func(k, v []byte) error {
			copiedNode := copied[node]
			if v != nil {
				// Key/Value pair
				pair := Pair{Key: copyBytes(k), Value: copyBytes(v)}
				copiedNode.Pairs = append(copiedNode.Pairs, &pair)
			} else {
				// Bucket
				copiedBucket := Bucket{Name: copyBytes(k)}
				copiedNode.Buckets = append(copiedNode.Buckets, &copiedBucket)

				bucket := node.Bucket(k)
				if bucket == nil {
					return errors.New("unable to retrieve bucket")
				}
				copied[bucket] = &copiedBucket
				queue = append(queue, bucket)
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}
	return &root, nil
}

func copyBytes(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}
