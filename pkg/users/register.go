package users

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type UserRegistration struct {
	User *User `json:"user"`
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type UserHandler struct {
	Path           string
	UserRepository UserRepository
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var req UserRegistration
	_ = json.Unmarshal(body, &req)
	reqUser := req.User
	_ = u.UserRepository.RegisterUser(reqUser)

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
}
