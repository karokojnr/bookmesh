package main

import "context"

type storage struct {
	// mongo db
}

func NewStorage() *storage {
	return &storage{}
}

func (s *storage) Create(ctx context.Context) error {
	return nil
}
