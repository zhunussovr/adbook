package main

import (
	"fmt"
	"strings"

	"gopkg.in/ldap.v3"
)

type Ldap struct {
	conn   *ldap.Conn
	config *LdapConfig
}

type LdapConfig struct {
	LdapServer   string
	LdapBind     string
	LdapPassword string
	FilterDN     string
	BaseDN       string
}

func NewLdap(config *LdapConfig) (*Ldap, error) {
	con, err := ConnectLdap(config)
	if err != nil {
		return nil, err
	}

	return &Ldap{con, config}, nil
}

func ConnectLdap(config *LdapConfig) (*ldap.Conn, error) {
	conn, err := ldap.Dial("tcp", config.LdapServer)

	if err != nil {
		return nil, fmt.Errorf("Failed to connect. %v", err)
	}

	if err := conn.Bind(config.LdapBind, config.LdapPassword); err != nil {
		return nil, fmt.Errorf("Failed to bind. %v", err)
	}

	return conn, nil
}

func (l *Ldap) GetLDAPUsers(username string) ([]Employee, error) {
	var userlist []Employee

	result, err := l.conn.Search(ldap.NewSearchRequest(
		l.config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		l.Filter(username),
		[]string{"dn", "sAMAccountName", "mail", "sn", "displayName", "telephoneNumber", "title"},
		nil,
	))

	if err != nil {
		return nil, fmt.Errorf("Failed to search users. %v", err)
	}

	for _, entry := range result.Entries {

		user := Employee{
			AccountName: entry.GetAttributeValue("sAMAccountName"),
			FullName:    entry.GetAttributeValue("displayName"),
			Title:       entry.GetAttributeValue("title"),
			Email:       entry.GetAttributeValue("mail"),
			Phone:       entry.GetAttributeValue("telephoneNumber"),
		}
		userlist = append(userlist, user)

	}

	// debug
	fmt.Printf("%+v\n", userlist)

	return userlist, nil
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
