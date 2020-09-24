package adbook

import (
	"fmt"

	"github.com/zhunussovr/adbook/backend"
	"github.com/zhunussovr/adbook/model"
)

type BookService struct {
	Backends map[string]backend.Interface
}

func NewBookService(backends map[string]backend.Interface) *BookService {
	return &BookService{backends}
}

func (b *BookService) Get(search string) ([]model.Person, error) {
	var persons []model.Person
	// TODO: Implement
	return persons, nil
}

func (b *BookService) Set(model.Person) error {
	// TODO: Implement
	return nil
}

func (b *BookService) Search(search, backendName string) ([]model.Person, error) {
	var persons []model.Person

	backend, ok := b.Backends[backendName]
	if !ok {
		return nil, fmt.Errorf("no backend found with the name: %s", backendName)
	}

	persons, err := backend.Search(search)
	if err != nil {
		return nil, err
	}

	return persons, nil
}
