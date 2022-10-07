package bbrowse

type DB struct {
	Buckets []*Bucket
}

type Bucket struct {
	Name    []byte
	Buckets []*Bucket
	Pairs   []*Pair
}

type Pair struct {
	Key   []byte
	Value []byte
}
