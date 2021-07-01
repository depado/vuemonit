package storage

import (
	bolt "go.etcd.io/bbolt"
)

type iterator struct {
	initial func() ([]byte, []byte)
	next    func() ([]byte, []byte)
	count   int
}

func (i *iterator) add() {
	i.count++
}

func NewIterator(c *bolt.Cursor, reverse bool) *iterator {
	if reverse {
		return reverseIterator(c)
	}
	return stdIterator(c)
}

func stdIterator(c *bolt.Cursor) *iterator {
	return &iterator{
		initial: c.First,
		next:    c.Next,
	}
}

func reverseIterator(c *bolt.Cursor) *iterator {
	return &iterator{
		initial: c.Last,
		next:    c.Prev,
	}
}
