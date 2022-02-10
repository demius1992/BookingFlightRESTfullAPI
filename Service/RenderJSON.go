package Service

import (
	"encoding/json"
	"net/http"
)

func RenderJSON(w http.ResponseWriter, v interface{}) {
	js, err1 := json.Marshal(v)
	if err1 != nil {
		http.Error(w, "problems with json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
