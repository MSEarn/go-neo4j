package config

type Config struct {
	Neo4j  Neo4j
	Server Server
}

type Neo4j struct {
	URI      string
	Port     string
	Username string
	Password string
}

type Server struct {
	Port         int
	WriteTimeout int
	ReadTimeout  int
	IdleTimeout  int
}
