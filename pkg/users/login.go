package users

import (
	"fmt"
	"net/http"

	"github.com/MSEarn/go-neo4j/pkg/auth"
	"github.com/MSEarn/go-neo4j/pkg/http_helper"
	"go.uber.org/zap"
)

type UserLogin struct {
	User User `json:"user"`
}

type LoggedInUser struct {
	Username string
	Email    string
	Token    string
}

type Response struct {
	Token string `json:"token"`
}

func Login(user UserRepository, sign auth.SignFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req UserLogin
		err := http_helper.Decode(r, &req)
		if err != nil {
			panic(err)
		}

		err = user.LoginUser(&req.User)
		if err != nil {
			panic(err)
		}

		user := req.User
		token, err := sign(r.Context(), user.Username, map[string]interface{}{
			"username": user.Username,
			"token":    "",
		})

		if err != nil {
			zap.L().Error(fmt.Sprintf("LoginHandlerFunc: sign jwt token failed %+v", err))
			http_helper.RespondOK(w, uint64(5000), err.Error(), nil)
			return
		}

		http_helper.RespondOK(w, uint64(0), "success", &Response{
			Token: token,
		})
	})
}
