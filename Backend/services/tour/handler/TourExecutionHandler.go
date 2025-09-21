package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"tour.xws.com/model"
	"tour.xws.com/service"
)

type TourExecutionHandler struct {
	Svc *service.TourExecutionService
}

type startReq struct {
	UserId string `json:"userId"`
	TourId string `json:"tourId"`
}

func (h *TourExecutionHandler) Start(w http.ResponseWriter, r *http.Request) {
	var req startReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if req.UserId == "" || req.TourId == "" {
		http.Error(w, "userId and tourId are required", http.StatusBadRequest)
		return
	}

	tourUUID, err := uuid.Parse(req.TourId)
	if err != nil {
		http.Error(w, "invalid tourId", http.StatusBadRequest)
		return
	}

	te, err := h.Svc.Start(req.UserId, tourUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(te)
}

type checkReq struct {
    UserId string            `json:"userId"`
    Coords model.Coordinates `json:"coords"`
}

func (h *TourExecutionHandler) CheckProximity(w http.ResponseWriter, r *http.Request) {
	var req checkReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if req.UserId == "" {
		http.Error(w, "userId required", http.StatusBadRequest)
		return
	}

	res, err := h.Svc.CheckProximity(req.UserId, req.Coords)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

type abandonReq struct {
    UserId string `json:"userId"`
}

func (h *TourExecutionHandler) Abandon(w http.ResponseWriter, r *http.Request) {
    var req abandonReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "bad request", http.StatusBadRequest); return
    }
    if req.UserId == "" {
        http.Error(w, "userId required", http.StatusBadRequest); return
    }
    if err := h.Svc.Abandon(req.UserId); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest); return
    }
    w.WriteHeader(http.StatusNoContent)
}
type activeReq struct {
    UserId string `json:"userId"`
}
func (h *TourExecutionHandler) GetActive(w http.ResponseWriter, r *http.Request) {
    var req activeReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "bad request", http.StatusBadRequest); return
    }
    if req.UserId == "" {
        http.Error(w, "userId required", http.StatusBadRequest); return
    }
    te, err := h.Svc.GetActive(req.UserId)
    if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
    if te == nil { w.WriteHeader(http.StatusNoContent); return }
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(te)
}
