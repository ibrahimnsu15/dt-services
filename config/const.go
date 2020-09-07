package config

const (
	WORKSTATION = "workstation"
	DEV         = "dev"
	STAGE       = "stage"
	PROD        = "prod"

	CONSUL_KV           = "dt"
	CONSUL_HOST_DEV     = "localhost"
	CONSUL_PORT_DEV     = "8500"
	CONSUL_HOST_CLUSTER = "consul-server"
	CONSUL_PORT_CLUSTER = "8500"
)

var (
	MONGOHOSTS_WORKSTATION = []string{"localhost:27017"}
)
