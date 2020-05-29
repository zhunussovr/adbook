package main

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber"
)

type ApiServer struct {
	config *ApiConfig
	fibapp *fiber.App
	ldap   *Ldap
	cache  *Cache
}

type ApiConfig struct {
	ListenAddress string
}

func NewApiServer(config *ApiConfig, ldap *Ldap, cache *Cache) (*ApiServer, error) {
	app := fiber.New()
	return &ApiServer{config: config, fibapp: app, ldap: ldap, cache: cache}, nil
}

func (a *ApiServer) Routes() {
	a.fibapp.Get("/person/:id?", a.getPerson)
}

func (a *ApiServer) getPerson(c *fiber.Ctx) {

	if c.Params("id") == "" {
		c.SendStatus(500)
		return
	}

	id := c.Params("id")
	employees, err := a.ldap.GetLDAPUsers(id)
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

func (a *ApiServer) Run() error {
	log.Println("Starting server on ", a.config.ListenAddress)
	err := a.fibapp.Listen(a.config.ListenAddress)

	return err
}
