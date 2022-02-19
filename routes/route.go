package routes

import (
	"net/http"

	"github.com/MSEarn/go-neo4j/config"
	"github.com/MSEarn/go-neo4j/pkg/auth"
	"github.com/MSEarn/go-neo4j/pkg/neo4j_driver"
	"github.com/MSEarn/go-neo4j/pkg/users"
	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func Setup(neo4jCfg config.Neo4j, neo4jDriver neo4j_driver.Driver, jwt *auth.JWT) *mux.Router {
	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()

	userRepository := &users.UserNeo4jRepository{
		Driver: neo4jDriver(neo4jCfg.URI, neo4j.BasicAuth(neo4jCfg.Username, neo4jCfg.Password, "")),
	}
	api.Handle("/users/register", users.Register(userRepository, jwt)).Methods(http.MethodPost)
	api.Handle("/auth/login", users.Login(userRepository, auth.NewSignFunc(jwt))).Methods(http.MethodPost)

	return router
}
