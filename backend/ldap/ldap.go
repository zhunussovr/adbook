package ldap

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/zhunussovr/adbook/model"
)

type Ldap struct {
	conn   *ldap.Conn
	config Config
}

type Config struct {
	Name     string
	Server   string
	Bind     string
	Password string
	FilterDN string
	BaseDN   string
}

func New(config Config) (*Ldap, error) {
	con, err := ConnectLdap(config)
	if err != nil {
		return nil, err
	}

	return &Ldap{con, config}, nil
}

func ConnectLdap(config Config) (*ldap.Conn, error) {
	conn, err := ldap.Dial("tcp", config.Server)

	if err != nil {
		return nil, fmt.Errorf("failed to connect. %v", err)
	}

	if err := conn.Bind(config.Bind, config.Password); err != nil {
		return nil, fmt.Errorf("failed to bind. %v", err)
	}

	return conn, nil
}

func (l *Ldap) Search(search string) ([]model.Person, error) {
	var persons []model.Person

	result, err := l.conn.Search(ldap.NewSearchRequest(
		l.config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		l.Filter(search),
		[]string{"uid", "cn", "mail", "sn"},
		nil,
	))

	if err != nil {
		return nil, fmt.Errorf("Failed to search users. %v", err)
	}

	for _, entry := range result.Entries {

		user := model.Person{
			ID:       entry.GetAttributeValue("uid"),
			FullName: entry.GetAttributeValue("cn"),
			LastName: entry.GetAttributeValue("sn"),
			Email:    entry.GetAttributeValue("mail"),
			// Phone:     entry.GetAttributeValue("telephoneNumber"),
		}
		persons = append(persons, user)

	}

	// debug
	fmt.Printf("%+v\n", persons)

	return persons, nil
}

func (l *Ldap) Auth(login, pass string) error {
	result, err := l.conn.Search(ldap.NewSearchRequest(
		l.config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		l.Filter(login),
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

	if err := l.conn.Bind(result.Entries[0].DN, pass); err != nil {
		fmt.Printf("Failed to auth. %s", err)
	} else {
		fmt.Printf("Authenticated successfuly!")
	}

	return nil
}

func (l *Ldap) Filter(needle string) string {
	res := strings.Replace(
		l.config.FilterDN,
		"{username}",
		needle,
		-1,
	)

	return res
}
