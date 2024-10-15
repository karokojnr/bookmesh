package main

type storage struct {
	// mongo db
}

func NewStorage() *storage {
	return &storage{}
}

func (s *storage) Create() error {
	return nil
}
