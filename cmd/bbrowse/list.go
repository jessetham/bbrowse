package main

import (
	"bbrowse"
	"fmt"
)

type Bucket bbrowse.Bucket

func (b Bucket) Title() string { return string(b.Name) }
func (b Bucket) Description() string {
	return fmt.Sprintf("# of nested buckets: %d | # of pairs: %d", len(b.Buckets), len(b.Pairs))
}
func (b Bucket) FilterValue() string { return string(b.Name) }

type Pair bbrowse.Pair

func (p Pair) Title() string       { return string(p.Key) }
// TODO: Need to convert value to string for structs
func (p Pair) Description() string { return string(p.Value) }
func (p Pair) FilterValue() string { return string(p.Key) }
