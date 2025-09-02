package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"tour.xws.com/model"
	"tour.xws.com/service"
)

type TourHandler struct {
	TourService *service.TourService
}

func (handler *TourHandler) GetAll(writer http.ResponseWriter, req *http.Request) {
	tours, err := handler.TourService.GetAllTours()
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(tours)
}

func (handler *TourHandler) Create(writer http.ResponseWriter, req *http.Request) {
	var tour model.Tour
	err := json.NewDecoder(req.Body).Decode(&tour)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.TourService.Create(&tour)
	if err != nil {
		println("Error while creating a new tour")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
	println("Tour successfully created")
}

func (handler *TourHandler) Delete(writer http.ResponseWriter, req *http.Request) {
	idStr := mux.Vars(req)["id"]
	log.Printf("Deleting tour with ID: %s", idStr)

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(writer, "Invalid UUID", http.StatusBadRequest)
		return
	}

	err = handler.TourService.Delete(id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]string{"message": "Tour deleted successfully"})
}

func (handler *TourHandler) Update(writer http.ResponseWriter, req *http.Request) {
	idStr := mux.Vars(req)["id"]
	log.Printf("Updating tour with ID: %s", idStr)

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(writer, "Invalid UUID", http.StatusBadRequest)
		return
	}

	var input struct {
		AuthorId      string               `json:"authorId"`
		Title         string               `json:"title"`
		Description   string               `json:"description"`
		Difficulty    model.TourDifficulty `json:"difficulty"`
		Tags          []string             `json:"tags"`
		Status        model.TourStatus     `json:"status"`
		Price         float64              `json:"price"`
		Distance      float64              `json:"distance"`
		PublishedAt   time.Time            `json:"publishedAt"`
		ArchivedAt    time.Time            `json:"archivedAt"`
		Duration      float64              `json:"duration"`
		TransportType model.TransportType  `json:"transportType"`
	}

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	updatedTour := model.Tour{
		AuthorId:     input.AuthorId,
		Title:       input.Title,
		Description: input.Description,
		Difficulty:  input.Difficulty,
		Tags:       input.Tags,
		Status:     input.Status,
		Price:     input.Price,
		Distance:  input.Distance,
		PublishedAt: input.PublishedAt,
		ArchivedAt:  input.ArchivedAt,
		Duration:    input.Duration,
		TransportType: input.TransportType,
	}

	err = handler.TourService.Update(id, updatedTour)
	if err != nil {
		http.Error(writer, "Tour not found", http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]string{"message": "Tour updated successfully"})
}

func (handler *TourHandler) GetAllByAuthor(writer http.ResponseWriter, req *http.Request) {
	userID := mux.Vars(req)["userId"]

	writer.Header().Set("Content-Type", "application/json")

	tours, err := handler.TourService.GetAllByAuthor(userID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(tours)
}

func (handler *TourHandler) GetById(writer http.ResponseWriter, req *http.Request) {
	idStr := mux.Vars(req)["id"]
	log.Printf("Fetching tour with ID: %s", idStr)

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(writer, "Invalid UUID", http.StatusBadRequest)
		return
	}

	tour, err := handler.TourService.GetById(id)
	if err != nil {
		// if repository returned not found, respond 404
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(tour)
}
