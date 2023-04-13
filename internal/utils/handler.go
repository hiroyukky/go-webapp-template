package utils

import(
        "fmt"
        "net/http"
        "encoding/json"
)

func ResponseJson(w http.ResponseWriter, res interface{}) {
     json, e := json.Marshal(res)
     if e != nil {
        http.Error(w, e.Error(), http.StatusInternalServerError)
        return
     }
     w.Header().Set("Content-Type", "application/json")
     fmt.Fprint(w, string(json))
}