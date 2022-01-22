package users_test

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"

	"github.com/MSEarn/go-neo4j/pkg/users"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type FakeUserRepository struct {
}

func (FakeUserRepository) RegisterUser(user *users.User) error {
	return nil
}

var _ = Describe("Users", func() {
	It("should be register", func() {
		handler := users.UserHandler{
			Path:           "/users",
			UserRepository: &FakeUserRepository{},
		}
		reqBody := bytes.NewReader([]byte(`{"user": {"email": "user@example.com","password": "12345678","username": "user"}}
		`))

		wTest := httptest.NewRecorder()
		handler.Register(wTest, httptest.NewRequest("POST", "/users", reqBody))

		Expect(wTest.Code).To(Equal(201))
		resp, _ := ioutil.ReadAll(wTest.Body)
		Expect(string(resp)).To(Equal("{\"user\":{\"username\":\"user\",\"email\":\"user@example.com\"}}"))
	})
})
