package adbook

import "github.com/zhunussovr/adbook/backend/ldap"

type Config struct {
	LdapConfigs []ldap.Config
	APIConfig   ApiConfig
}

func ParseConfig(s string) Config {
	// TODO: make real parsing of input toml string
	return Config{
		LdapConfigs: []ldap.Config{
			ldap.Config{
				Name:     "OpenLDAP01",
				Server:   "localhost:389",
				Bind:     "cn=admin,dc=example,dc=org",
				Password: "admin",
				FilterDN: "(&(objectClass=inetOrgPerson)(|(uid=*{username}*)(mail=*{username}*)))",
				BaseDN:   "dc=example,dc=org",
			},
		},
		APIConfig: ApiConfig{
			ListenAddress: ":8080",
		},
	}
}
