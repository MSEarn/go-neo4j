package neo4j_driver

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Driver func(target string, token neo4j.AuthToken) neo4j.Driver

func NewDriver() Driver {
	return func(target string, token neo4j.AuthToken) neo4j.Driver {
		driver, err := neo4j.NewDriver(target, token)
		fmt.Printf("%s\n", target)
		fmt.Printf("%s\n", token)
		if err != nil {
			panic(err)
		}

		return driver
	}
}
