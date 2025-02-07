package routes

import (
	"alphalabz/pkg/casbin"
	"alphalabz/pkg/pocketbase"
	"encoding/json"
	"net/http"
)

type testData struct {
	Jwt string `json:"jwt"`
}

func TestHandler(w http.ResponseWriter, r *http.Request, p *pocketbase.PocketBaseClient, ce *casbin.CasbinEnforcer) {
	var testData testData
	// decode body
	if err := json.NewDecoder(r.Body).Decode(&testData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// send response
	// json.NewEncoder(w).Encode(result)
}
