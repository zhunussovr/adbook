package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/zhunussovr/adbook"
	"github.com/zhunussovr/adbook/backend"
	"github.com/zhunussovr/adbook/backend/ldap"
)

var (
	conf = flag.String("conf", "config.toml", "path to a configuration file")
)

func main() {

	flag.Parse()

	configData, err := ioutil.ReadFile(*conf)
	if err != nil {
		log.Fatal("error reading config:", err)
	}

	log.Println("Parsing config ...")
	config := adbook.ParseConfig(string(configData))

	backends := make(map[string]backend.Interface)

	log.Println("Initializing ldaps ...")
	for _, ldapConfig := range config.LdapConfigs {
		ldapBackend, err := ldap.New(ldapConfig)
		if err != nil {
			log.Println(err)
		}

		backends[ldapConfig.Name] = ldapBackend

	}

	bookService := adbook.NewBookService(backends)

	srv, _ := adbook.NewApiServer(config.APIConfig, bookService)
	if err != nil {
		log.Println(err)
	}

	srv.Routes()

	log.Fatalln(srv.Run())

}
