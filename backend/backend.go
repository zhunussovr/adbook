package backend

import "github.com/zhunussovr/adbook/model"

type Interface interface {
	// Get(string) ([]model.Person, error)
	// Set(model.Person) error
	Search(search string) ([]model.Person, error)
}
