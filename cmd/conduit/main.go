package main

import (
	"net/http"
	"os"

	"github.com/MSEarn/go-neo4j/pkg/users"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func main() {
	neo4jURI, found := os.LookupEnv("NEO4J_URI")
	if !found {
		panic("NEO4J_URI not set")
	}
	neo4jUsername, found := os.LookupEnv("NEO4J_USERNAME")
	if !found {
		panic("NEO4J_USERNAME not set")
	}
	neo4jPassword, found := os.LookupEnv("NEO4J_PASSWORD")
	if !found {
		panic("NEO4J_PASSWORD not set")
	}

	handler := &users.UserHandler{
		Path: "/users",
		UserRepository: &users.UserNeo4jRepository{
			Driver: driver(neo4jURI, neo4j.BasicAuth(neo4jUsername, neo4jPassword, "")),
		},
	}

	server := http.NewServeMux()
	server.HandleFunc("/users", handler.Register)

	if err := http.ListenAndServe("3000", server); err != nil {
		panic(err)
	}
}

func driver(target string, token neo4j.AuthToken) neo4j.Driver {
	driver, err := neo4j.NewDriver(target, token)
	if err != nil {
		panic(err)
	}

	return driver
}
