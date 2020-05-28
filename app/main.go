package main

import (
	"flag"
	"io/ioutil"
	"log"
)

var (
	conf = flag.String("conf", "config.yaml", "path to a configuration file")
)

func main() {

	flag.Parse()

	configData, _ := ioutil.ReadFile(*conf)

	log.Println("Parsing config ...")
	config := ParseConfig(string(configData))

	log.Println("Connecting to ldap ...")
	ldap, err := NewLdap(config.ldapConfig)
	if err != nil {
		log.Println(err)
	}

	log.Println("Connecting to cache...")
	cache, _ := NewCache(config.cacheConfig)
	if err != nil {
		log.Println(err)
	}

	srv, _ := NewApiServer(config.apiConfig, ldap, cache)
	if err != nil {
		log.Println(err)
	}

	srv.Routes()

	log.Fatalln(srv.Run())

}
