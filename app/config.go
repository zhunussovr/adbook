package main

type Config struct {
	ldapConfig  *LdapConfig
	apiConfig   *ApiConfig
	cacheConfig *CacheConfig
}

func ParseConfig(s string) *Config {
	// TODO: make real parsing of incoming string
	return &Config{
		ldapConfig: &LdapConfig{
			LdapServer:   "192.168.15.5:389",
			LdapBind:     "vagrant@uplift.local",
			LdapPassword: "vagrant",
			FilterDN:     "(&(objectClass=person)(memberOf:1.2.840.113556.1.4.1941:=CN=staff,OU=uGroups,OU=DEMO01,DC=uplift,DC=local)(|(sAMAccountName={username})(mail={username})))",
			BaseDN:       "OU=DEMO01,DC=uplift,DC=local",
		},
		apiConfig: &ApiConfig{
			ListenAddress: ":8080",
		},
		cacheConfig: &CacheConfig{
			Username:       "adbookadm",
			Password:       "adbookadm",
			DbName:         "adbook",
			CollectionName: "userdata",
			Port:           8081,
		},
	}
}
