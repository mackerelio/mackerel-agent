package config

var DefaultConfig = &Config{
	Apibase: "https://mackerel.io",
	Root:    ".",
	Pidfile: "mackerel-agent.pid",
	Conffile: "mackerel-agent.conf",
	Roles:   []string{},
	Verbose: false,
}
