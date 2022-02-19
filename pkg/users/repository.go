package users

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	RegisterUser(user *User) error
	LoginUser(user *User) error
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

func (u *UserNeo4jRepository) LoginUser(user *User) error {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer func() {
		session.Close()
	}()

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return u.authenticate(tx, user)
	}); err != nil {
		return err
	}

	return nil
}

func (u *UserNeo4jRepository) authenticate(tx neo4j.Transaction, user *User) (interface{}, error) {
	query := "MATCH (u:User {username: $username, email: $email}) RETURN u.username AS username, u.email AS email, u.password AS password"
	parameters := map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
	}
	res, err := tx.Run(query, parameters)
	if err != nil {
		return nil, err
	}
	rec, err := res.Single()
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", rec)
	if !isPasswordMatched(user.Password, rec.Values[2].(string)) {
		return nil, errors.New("Unauthenticated")
	}

	return nil, err
}

func isPasswordMatched(initPassword, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(initPassword)) == nil
}
