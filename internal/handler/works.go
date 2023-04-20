package handler

import (
	"fmt"
	"net/http"

	"go-template/internal/middleware"
	"go-template/internal/utils"
)

var WorksListHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(middleware.JWTContextKey{}).(middleware.AuthClaims)

	fmt.Printf("jwt token %+v\n", claims)
	if claims.AccountId == 0 {
		http.Error(w, "token is missing", http.StatusInternalServerError)
		return
	}

	//claims := token.Claims.(middleware.AuthClaims)
	//fmt.Printf("claims %+v\n", claims)
	type ResponseResult struct {
		Result bool `json:"result"`
	}
	res := ResponseResult{
		Result: true,
	}
	utils.ResponseJson(w, res)
})
