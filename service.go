package adbook

import (
	"github.com/zhunussovr/adbook/backend"
	"github.com/zhunussovr/adbook/model"
	"golang.org/x/sync/errgroup"
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
	var persons []model.Person
	var personsCh chan []model.Person
	g := new(errgroup.Group)

	for _, backend := range b.Backends {
		g.Go(func() error {
			persons, err := backend.Search(search)
			if err != nil {
				return err
			}
			personsCh <- persons
			return nil
		})
	}

	for range b.Backends {
		persons = append(persons, <-personsCh...)
	}

	return persons, nil
}
