package main

// Storage defines how todos get persisted, regardless of backend.
type Storage[T any] interface {
	Save(data T) error
	Load(data *T) error
}
