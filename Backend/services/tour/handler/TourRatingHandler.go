package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"tour.xws.com/model"
	"tour.xws.com/service"
)

type TourRatingHandler struct {
	RatingService *service.TourRatingService
}

type createRatingReq struct {
	Rating       int    `json:"rating"`
	Comment      string `json:"comment"`
	TouristName  string `json:"touristName"`
	TouristEmail string `json:"touristEmail"`
	VisitedAt    string `json:"visitedAt"`   // yyyy-mm-dd
	CommentedAt  string `json:"commentedAt"` // yyyy-mm-dd
}

func parseDate(d string) (time.Time, error) {
	if d == "" { return time.Time{}, nil }
	return time.Parse("2006-01-02", d)
}

func (h *TourRatingHandler) Create(w http.ResponseWriter, r *http.Request) {
	tourId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil { http.Error(w, "Invalid tour id", http.StatusBadRequest); return }

	var body createRatingReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Bad JSON", http.StatusBadRequest); return
	}

	visitedAt, err := parseDate(body.VisitedAt)
	if err != nil { http.Error(w, "Invalid visitedAt", http.StatusBadRequest); return }
	commentedAt, err := parseDate(body.CommentedAt)
	if err != nil { http.Error(w, "Invalid commentedAt", http.StatusBadRequest); return }

	item := &model.TourRating{
		Id:           uuid.New(),
		TourId:       tourId,
		Rating:       body.Rating,
		Comment:      body.Comment,
		TouristName:  body.TouristName,
		TouristEmail: body.TouristEmail,
		VisitedAt:    visitedAt,
		CommentedAt:  commentedAt,
		CreatedAt:    time.Now().UTC(),
	}

	if err := h.RatingService.Create(item); err != nil {
		http.Error(w, "Failed to save review", http.StatusInternalServerError); return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"id": item.Id})
}

func (h *TourRatingHandler) GetByTour(w http.ResponseWriter, r *http.Request) {
	tourId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil { http.Error(w, "Invalid tour id", http.StatusBadRequest); return }

	items, err := h.RatingService.GetByTour(tourId)
	if err != nil { http.Error(w, "Failed to load reviews", http.StatusInternalServerError); return }

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(items)
}
