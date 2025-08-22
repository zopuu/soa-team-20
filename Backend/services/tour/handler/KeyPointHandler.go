package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"tour.xws.com/model"
	"tour.xws.com/service"
)

type KeyPointHandler struct {
	KeyPointService *service.KeyPointService
}

func (handler *KeyPointHandler) GetAll(writer http.ResponseWriter, req *http.Request) {
	keyPoints, err := handler.KeyPointService.GetAllKeyPoints()
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(keyPoints)
}

func (handler *KeyPointHandler) GetAllByTour(writer http.ResponseWriter, req *http.Request) {
	tourIdStr := mux.Vars(req)["tourId"]

	tourId, err := uuid.Parse(tourIdStr)
	if err != nil {
		http.Error(writer, "Invalid UUID", http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")

	keyPoints, err := handler.KeyPointService.GetAllByTour(tourId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(keyPoints)
}

func (handler *KeyPointHandler) Create(writer http.ResponseWriter, req *http.Request) {
	var keyPoint model.KeyPoint
	err := json.NewDecoder(req.Body).Decode(&keyPoint)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.KeyPointService.Create(&keyPoint)
	if err != nil {
		println("Error while creating a new keyPoint")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
	println("KeyPoint successfully created")
}

func (handler *KeyPointHandler) Delete(writer http.ResponseWriter, req *http.Request) {
	idStr := mux.Vars(req)["id"]
	log.Printf("Deleting keyPoint with ID: %s", idStr)

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(writer, "Invalid UUID", http.StatusBadRequest)
		return
	}

	err = handler.KeyPointService.Delete(id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]string{"message": "KeyPoint deleted successfully"})
}

func (handler *KeyPointHandler) Update(writer http.ResponseWriter, req *http.Request) {
	idStr := mux.Vars(req)["id"]
	log.Printf("Updating keyPoint with ID: %s", idStr)

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(writer, "Invalid UUID", http.StatusBadRequest)
		return
	}

	var input struct {
		Title       string            `json:"title"`
		Description string            `json:"description"`
		Coordinates model.Coordinates `json:"coordinates"`
		Image       string            `json:"image"`
	}

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	updatedKeyPoint := model.KeyPoint{
		Coordinates: input.Coordinates,
		Title:       input.Title,
		Description: input.Description,
		Image:       input.Image,
	}

	err = handler.KeyPointService.Update(id, updatedKeyPoint)
	if err != nil {
		http.Error(writer, "KeyPoint not found", http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]string{"message": "KeyPoint updated successfully"})
}
