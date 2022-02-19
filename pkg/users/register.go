package users

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/MSEarn/go-neo4j/pkg/auth"
)

type UserRegistration struct {
	User *User `json:"user"`
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

func Register(
	user UserRepository,
	authJWT *auth.JWT,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var req UserRegistration
		_ = json.Unmarshal(body, &req)
		reqUser := req.User
		_ = user.RegisterUser(reqUser)

		w.WriteHeader(201)
		w.Header().Add("Content-Type", "application/json")
		resp := &UserRegistration{
			User: &User{
				Email:    reqUser.Email,
				Username: reqUser.Username,
			},
		}
		bytes, _ := json.Marshal(&resp)
		_, _ = w.Write(bytes)
	})
}
