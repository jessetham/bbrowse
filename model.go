package bbrowse

type Bucket struct {
	Name    []byte
	Buckets []*Bucket
	Pairs   []*Pair
}

type Pair struct {
	Key   []byte
	Value []byte
}
