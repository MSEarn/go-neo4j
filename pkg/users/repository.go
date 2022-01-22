package users

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	RegisterUser(user *User) error
}

type UserNeo4jRepository struct {
	Driver neo4j.Driver
}

func (u *UserNeo4jRepository) RegisterUser(user *User) error {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		session.Close()
	}()

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return u.persistUser(tx, user)
	}); err != nil {
		return err
	}

	return nil
}

func (u *UserNeo4jRepository) persistUser(tx neo4j.Transaction, user *User) (interface{}, error) {
	query := "CREATE (:User {email: $email, username: $username, password: $password})"
	hashPwd, err := hash(user.Password)
	if err != nil {
		return nil, err
	}
	parameters := map[string]interface{}{
		"email":    user.Email,
		"username": user.Username,
		"password": string(hashPwd),
	}
	_, err = tx.Run(query, parameters)

	return nil, err
}

func hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
