package handler

import(
        "net/http"

        "go-template/internal/utils"
)

var ServerStatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    type ResponseResult struct {
         Result bool `json:"result"`
    }
    res := ResponseResult{
        Result: true,
    }
    utils.ResponseJson(w, res)
})
