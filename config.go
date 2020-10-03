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
				Server:   "192.168.15.5:389",
				Bind:     "vagrant@uplift.local",
				Password: "vagrant",
				FilterDN: "(&(objectClass=person)(memberOf:1.2.840.113556.1.4.1941:=CN=staff,OU=uGroups,OU=DEMO01,DC=uplift,DC=local)(|(sAMAccountName={username})(mail={username})))",
				BaseDN:   "OU=DEMO01,DC=uplift,DC=local",
			},
		},
		APIConfig: ApiConfig{
			ListenAddress: ":8080",
		},
	}
}
