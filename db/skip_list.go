package db

type SkipList struct {
	node Node
}

type Node struct {
	key string
	next []Node
}