package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

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

func (handler *KeyPointHandler) GetAllByTourSortedByCreatedAt(writer http.ResponseWriter, req *http.Request) {
	tourIdStr := mux.Vars(req)["tourId"]

	tourId, err := uuid.Parse(tourIdStr)
	if err != nil {
		http.Error(writer, "Invalid UUID", http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")

	keyPoints, err := handler.KeyPointService.GetAllByTourSortedByCreatedAt(tourId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(keyPoints)
}

func (handler *KeyPointHandler) Create(writer http.ResponseWriter, req *http.Request) {
	// Parse multipart form data (10MB limit)
	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(writer, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// Extract form fields
	tourIdStr := req.FormValue("tourId")
	title := req.FormValue("title")
	description := req.FormValue("description")
	latitudeStr := req.FormValue("latitude")
	longitudeStr := req.FormValue("longitude")

	// Parse required fields
	tourId, err := uuid.Parse(tourIdStr)
	if err != nil {
		http.Error(writer, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		http.Error(writer, "Invalid latitude", http.StatusBadRequest)
		return
	}

	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		http.Error(writer, "Invalid longitude", http.StatusBadRequest)
		return
	}

	coordinates := model.Coordinates{
		Latitude:  latitude,
		Longitude: longitude,
	}

	// Handle image upload
	var image model.Image
	file, header, err := req.FormFile("image")
	if err == nil {
		defer file.Close()
		
		// Read image data
		imageData, err := io.ReadAll(file)
		if err != nil {
			http.Error(writer, "Failed to read image file", http.StatusInternalServerError)
			return
		}

		// Validate file type
		contentType := header.Header.Get("Content-Type")
		if !isValidImageType(contentType) {
			http.Error(writer, "Invalid image type. Only JPEG, PNG, GIF, and WebP are allowed", http.StatusBadRequest)
			return
		}

		image = model.Image{
			Data:     imageData,
			MimeType: contentType,
			Filename: header.Filename,
		}
	}
	// If no image uploaded, image will be empty struct

	// Create keypoint
	keyPoint := model.BeforeCreateKeyPoint(tourId, coordinates, title, description, image)
	
	err = handler.KeyPointService.Create(keyPoint)
	if err != nil {
		http.Error(writer, "Error while creating keypoint", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(keyPoint)
	log.Println("KeyPoint successfully created")
}

// Helper function to validate image types
func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png", 
		"image/gif",
		"image/webp",
	}
	
	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
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

	// Check content type to determine how to parse request
	contentType := req.Header.Get("Content-Type")
	
	var updatedKeyPoint model.KeyPoint
	
	if contentType == "application/json" {
		// Handle JSON request (existing behavior)
		var input struct {
			Title       string            `json:"title"`
			Description string            `json:"description"`
			Coordinates model.Coordinates `json:"coordinates"`
			Image       model.Image       `json:"image"`
		}

		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		updatedKeyPoint = model.KeyPoint{
			Coordinates: input.Coordinates,
			Title:       input.Title,
			Description: input.Description,
			Image:       input.Image,
		}
	} else {
		// Handle multipart form data (new behavior for file uploads)
		err := req.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(writer, "Failed to parse multipart form", http.StatusBadRequest)
			return
		}

		// Extract form fields
		title := req.FormValue("title")
		description := req.FormValue("description")
		latitudeStr := req.FormValue("latitude")
		longitudeStr := req.FormValue("longitude")

		// Parse coordinates
		latitude, err := strconv.ParseFloat(latitudeStr, 64)
		if err != nil {
			http.Error(writer, "Invalid latitude", http.StatusBadRequest)
			return
		}

		longitude, err := strconv.ParseFloat(longitudeStr, 64)
		if err != nil {
			http.Error(writer, "Invalid longitude", http.StatusBadRequest)
			return
		}

		coordinates := model.Coordinates{
			Latitude:  latitude,
			Longitude: longitude,
		}

		// Get existing keypoint to preserve image if no new image uploaded
		existingKeyPoint, err := handler.KeyPointService.GetById(id)
		if err != nil {
			http.Error(writer, "KeyPoint not found", http.StatusNotFound)
			return
		}

		// Handle image upload
		var image model.Image = existingKeyPoint.Image // Default to existing image
		file, header, err := req.FormFile("image")
		if err == nil {
			defer file.Close()
			
			// Read image data
			imageData, err := io.ReadAll(file)
			if err != nil {
				http.Error(writer, "Failed to read image file", http.StatusInternalServerError)
				return
			}

			// Validate image type
			mimeType := header.Header.Get("Content-Type")
			if mimeType != "image/jpeg" && mimeType != "image/png" && mimeType != "image/gif" && mimeType != "image/webp" {
				http.Error(writer, "Invalid image format. Only JPEG, PNG, GIF, and WebP are supported", http.StatusBadRequest)
				return
			}

			image = model.Image{
				Data:     imageData,
				MimeType: mimeType,
				Filename: header.Filename,
			}
		}

		updatedKeyPoint = model.KeyPoint{
			Coordinates: coordinates,
			Title:       title,
			Description: description,
			Image:       image,
		}
	}

	err = handler.KeyPointService.Update(id, updatedKeyPoint)
	if err != nil {
		http.Error(writer, "KeyPoint not found", http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]string{"message": "KeyPoint updated successfully"})
}

// GetImage serves the image data for a keypoint
func (handler *KeyPointHandler) GetImage(writer http.ResponseWriter, req *http.Request) {
	idStr := mux.Vars(req)["id"]
	
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(writer, "Invalid UUID", http.StatusBadRequest)
		return
	}

	keyPoint, err := handler.KeyPointService.GetById(id)
	if err != nil {
		http.Error(writer, "KeyPoint not found", http.StatusNotFound)
		return
	}

	// Check if image exists
	if len(keyPoint.Image.Data) == 0 {
		http.Error(writer, "No image found for this keypoint", http.StatusNotFound)
		return
	}

	// Set appropriate headers
	writer.Header().Set("Content-Type", keyPoint.Image.MimeType)
	writer.Header().Set("Content-Length", strconv.Itoa(len(keyPoint.Image.Data)))
	
	// Write image data
	writer.Write(keyPoint.Image.Data)
}
