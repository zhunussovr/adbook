package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

const dbName = "adbook"
const collectionName = "userdata"
const port = 8081

// Employee to define structure of json
type Employee struct {
	AccountName string //`json:"accountname,omitempty" bson:"accountname,omitempty"`
	FullName    string //`json:"fullname,omitempty" bson:"fullname,omitempty"`
	Title       string //`json:"title,omitempty" bson:"title,omitempty"`
	Email       string //`json:"email,omitempty" bson:"email,omitempty"`
	Phone       string //`json:"phone,omitempty" bson:"phone,omitempty"`
}

var client *mongo.Client

func getPerson(c *fiber.Ctx) {
	collection, err := GetCollections(dbName, collectionName)
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	var filter bson.M = bson.M{}

	if c.Params("id") != "" {
		id := c.Params("id")
		objID, _ := primitive.ObjectIDFromHex(id)
		filter = bson.M{"_id": objID}
	}

	var results []bson.M
	cur, err := collection.Find(context.Background(), filter)
	defer cur.Close(context.Background())

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	cur.All(context.Background(), &results)

	if results == nil {
		c.SendStatus(404)
		return
	}

	json, _ := json.Marshal(results)
	c.Send(json)
}

func main() {
	app := fiber.New()

	app.Get("/person/:id?", getPerson)
	app.Listen(port)

	// LDAP connection
	conn, err := connect()
	if err != nil {
		fmt.Printf("Failed to connect. %s", err)
		return
	}

	defer conn.Close()

	if err := getldapusers(conn); err != nil {
		fmt.Printf("%v", err)
		return
	}

	if err := auth(conn); err != nil {
		fmt.Printf("%v", err)
		return
	}
	// API block
	//http.Handle("/users", usersHandler)
	//log.Fatal(http.ListenAndServe(":8080", nil))
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

func getldapusers(conn *ldap.Conn) error {
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
			FullName:    entry.GetAttributeValue("displayName"),
			Title:       entry.GetAttributeValue("title"),
			Email:       entry.GetAttributeValue("mail"),
			Phone:       entry.GetAttributeValue("telephoneNumber"),
		}
		fmt.Printf("%+v\n", userlist)

		//collection := client.Database("adbook").Collection("userdata")
		//err := collection.InsertOne(context.TODO(), userlist)
		//if err != nil {
		//	log.Fatal(err)
		//}
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
