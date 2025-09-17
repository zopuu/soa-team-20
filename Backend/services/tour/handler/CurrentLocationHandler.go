package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"tour.xws.com/model"
	"tour.xws.com/service"
)

type CurrentLocationHandler struct {
	Svc *service.CurrentLocationService
}

func (h *CurrentLocationHandler) Get(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["userId"]
	loc, err := h.Svc.Get(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError); return
	}
	if loc == nil {
		w.WriteHeader(http.StatusNoContent); return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loc)
}

type setReq struct {
	UserId     string           `json:"userId"`
	Coordinates model.Coordinates `json:"coordinates"`
}

func (h *CurrentLocationHandler) Set(w http.ResponseWriter, r *http.Request) {
	var req setReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest); return
	}
	if req.UserId == "" {
		http.Error(w, "userId required", http.StatusBadRequest); return
	}
	if err := h.Svc.Set(req.UserId, req.Coordinates); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError); return
	}
	w.WriteHeader(http.StatusNoContent)
}
