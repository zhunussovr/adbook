package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/ldap.v3"
)

const (
	ldapServer = "192.168.15.5:389"
	//ldapPort     = 389
	ldapBind     = "vagrant@uplift.local"
	ldapPassword = "vagrant"

	filterDN = "(&(objectClass=person)(memberOf:1.2.840.113556.1.4.1941:=CN=staff,OU=uGroups,OU=DEMO01,DC=uplift,DC=local)(|(sAMAccountName={username})(mail={username})))"
	baseDN   = "OU=DEMO01,DC=uplift,DC=local"

	loginUsername = "vagrant"
	loginPassword = "vagrant"
)

type uData struct {
	FullName, Title, Email, Phone string
}

// Employee to define structure of json
type Employee struct {
	AccountName string
	UserData    []uData
}

func main() {
	conn, err := connect()

	if err != nil {
		fmt.Printf("Failed to connect. %s", err)
		return
	}

	defer conn.Close()

	if err := list(conn); err != nil {
		fmt.Printf("%v", err)
		return
	}

	if err := auth(conn); err != nil {
		fmt.Printf("%v", err)
		return
	}
}

func connect() (*ldap.Conn, error) {
	conn, err := ldap.Dial("tcp", ldapServer)

	if err != nil {
		return nil, fmt.Errorf("Failed to connect. %s", err)
	}

	if err := conn.Bind(ldapBind, ldapPassword); err != nil {
		return nil, fmt.Errorf("Failed to bind. %s", err)
	}

	return conn, nil
}

func list(conn *ldap.Conn) error {
	result, err := conn.Search(ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter("*"),
		[]string{"dn", "sAMAccountName", "mail", "sn", "displayName", "telephoneNumber", "title"},
		nil,
	))

	if err != nil {
		return fmt.Errorf("Failed to search users. %s", err)
	}

	for _, entry := range result.Entries {

		userlist := Employee{
			AccountName: entry.GetAttributeValue("sAMAccountName"),
			UserData: []uData{
				{
					FullName: entry.GetAttributeValue("displayName"),
					Title:    entry.GetAttributeValue("title"),
					Email:    entry.GetAttributeValue("mail"),
					Phone:    entry.GetAttributeValue("telephoneNumber"),
				},
			},
		}

		file, _ := json.MarshalIndent(userlist, "", " ")
		_ = ioutil.WriteFile("users.json", file, 0644)

		fmt.Printf("%+v\n", userlist)

		// For using SCRAM-SHA-1 auth.mechanism
		dbcreds := options.Credential{
			Username: "adbookadm",
			Password: "adbookadm",
		}
		clientOptions := options.Client().ApplyURI("mongodb://localhost/adbook").SetAuth(dbcreds)
		client, err := mongo.Connect(context.TODO(), clientOptions)

		if err != nil {
			log.Fatal(err)
		}

		err = client.Ping(context.TODO(), nil)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Connected to MongoDB!")

		//database block - willb be moved later
		collection := client.Database("adbook").Collection("userdata")
		insertResult, err := collection.InsertOne(context.TODO(), userlist)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted multiple documents: ", insertResult.InsertedID)
	}

	return nil
}

func auth(conn *ldap.Conn) error {
	result, err := conn.Search(ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter(loginUsername),
		[]string{"dn"},
		nil,
	))

	if err != nil {
		return fmt.Errorf("Failed to find user. %s", err)
	}

	if len(result.Entries) < 1 {
		return fmt.Errorf("User does not exist")
	}

	if len(result.Entries) > 1 {
		return fmt.Errorf("Too many entries returned")
	}

	if err := conn.Bind(result.Entries[0].DN, loginPassword); err != nil {
		fmt.Printf("Failed to auth. %s", err)
	} else {
		fmt.Printf("Authenticated successfuly!")
	}

	return nil
}

func filter(needle string) string {
	res := strings.Replace(
		filterDN,
		"{username}",
		needle,
		-1,
	)

	return res
}
