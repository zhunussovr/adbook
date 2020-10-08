package adbook

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber"
)

type ApiServer struct {
	book   *BookService
	config ApiConfig
	fibapp *fiber.App
}

type ApiConfig struct {
	ListenAddress string
}

func NewApiServer(config ApiConfig, b *BookService) (*ApiServer, error) {
	app := fiber.New()
	return &ApiServer{config: config, fibapp: app, book: b}, nil
}

func (a *ApiServer) Routes() {
	a.fibapp.Static("/", "./web")
	a.fibapp.Get("/api/search", a.search)
	a.fibapp.Get("/api/person/:id?", a.getPerson)
}

func (a *ApiServer) getPerson(c *fiber.Ctx) {

	if c.Params("id") == "" {
		c.SendStatus(500)
		return
	}

	id := c.Params("id")
	employees, err := a.book.Get(id)
	if err != nil {
		c.SendStatus(500)
		return
	}

	if employees == nil {
		c.SendStatus(404)
		return
	}

	json, _ := json.Marshal(employees)
	c.Send(json)
}

func (a *ApiServer) search(c *fiber.Ctx) {
	q := c.Query("query")
	b := c.Query("backend")
	log.Printf("parsing query and backend: %v %v", q, b)
	persons, err := a.book.Search(q, b)
	if err != nil {
		log.Printf("error parsing: %v", err)
		c.SendStatus(500)
		return
	}

	if persons == nil {
		c.SendStatus(404)
		return
	}

	json, _ := json.Marshal(persons)
	c.Send(json)
}

func (a *ApiServer) Run() error {
	log.Println("Starting server on ", a.config.ListenAddress)
	err := a.fibapp.Listen(a.config.ListenAddress)

	return err
}
