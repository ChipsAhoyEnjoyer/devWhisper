package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg Config) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	//  TODO: Validate incoming UUID
	//  TODO: Authenticate user before allowing deletion
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	params := struct {
		ID uuid.UUID `json:"id"`
	}{}
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err = cfg.DB.DeleteUser(r.Context(), params.ID)
	if err != nil {
		http.Error(w, "Error deleting user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
