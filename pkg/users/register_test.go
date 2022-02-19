package users_test

import (
	"bytes"
	"io"
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

func (FakeUserRepository) LoginUser(user *users.User) error {
	return nil
}

var _ = Describe("Users", func() {
	It("should be register", func() {
		registerHandler := users.Register(&FakeUserRepository{}, nil)
		reqBody := bytes.NewReader([]byte(`{"user": {"email": "user@example.com","password": "12345678","username": "user"}}
		`))

		req := httptest.NewRequest("POST", "https://example/api/v1/auth/login", reqBody)
		wTest := httptest.NewRecorder()
		registerHandler.ServeHTTP(wTest, req)
		res := wTest.Result()

		Expect(res.StatusCode).To(Equal(201))
		resp, _ := io.ReadAll(res.Body)
		Expect(string(resp)).To(Equal("{\"user\":{\"username\":\"user\",\"email\":\"user@example.com\"}}"))
	})
})
