package adbook

import (
	"github.com/zhunussovr/adbook/backend"
	"github.com/zhunussovr/adbook/model"
)

type BookService struct {
	Backends map[string]backend.Interface
}

func NewBookService() *BookService {
	return &BookService{make(map[string]backend.Interface)}
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

func (b *BookService) Search(search string) ([]model.Person, error) {
	// TODO: Implement
	return nil
}
