package main

import (
	"fmt"
	"strings"

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

type Employee struct {
	AccountName string `json:"accountname,omitempty" bson:"accountname,omitempty"`
	FullName    string `json:"fullname,omitempty" bson:"fullname,omitempty"`
	Title       string `json:"title,omitempty" bson:"title,omitempty"`
	Email       string `json:"email,omitempty" bson:"email,omitempty"`
	Phone       string `json:"phone,omitempty" bson:"phone,omitempty"`
}

func Connect() (*ldap.Conn, error) {
	conn, err := ldap.Dial("tcp", ldapServer)

	if err != nil {
		return nil, fmt.Errorf("Failed to connect. %s", err)
	}

	if err := conn.Bind(ldapBind, ldapPassword); err != nil {
		return nil, fmt.Errorf("Failed to bind. %s", err)
	}

	return conn, nil
}

func GetLDAPUsers(conn *ldap.Conn) error {
	result, err := conn.Search(ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		Filter("*"),
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
	}
	return nil
}

func Auth(conn *ldap.Conn) error {
	result, err := conn.Search(ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		Filter(loginUsername),
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

func Filter(needle string) string {
	res := strings.Replace(
		filterDN,
		"{username}",
		needle,
		-1,
	)

	return res
}
