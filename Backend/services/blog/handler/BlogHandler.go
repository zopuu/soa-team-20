package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"blog.xws.com/model"
	"blog.xws.com/service"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type BlogHandler struct {
	BlogService *service.BlogService
}

func (handler *BlogHandler) GetAll(writer http.ResponseWriter, req *http.Request) {
	blogs, err := handler.BlogService.GetAllBlogs()
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(blogs)
}

func (handler *BlogHandler) GetById(writer http.ResponseWriter, req *http.Request) {
	idStr := mux.Vars(req)["id"]
	log.Printf("Blog sa id-em %s", idStr)

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(writer, "Invalid UUID", http.StatusBadRequest)
		return
	}

	blog, err := handler.BlogService.GetById(id) // make sure your service takes uuid.UUID
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(blog)
}

func (handler *BlogHandler) Create(writer http.ResponseWriter, req *http.Request) {
	// Check content type to determine how to parse request
	contentType := req.Header.Get("Content-Type")
	
	var blog *model.Blog
	var err error

	if contentType == "application/json" {
		// Handle JSON request (existing behavior)
		var blogData model.Blog
		err = json.NewDecoder(req.Body).Decode(&blogData)
		if err != nil {
			println("Error while parsing json")
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		blog = &blogData
	} else {
		// Handle multipart form data (new behavior for file uploads)
		err = req.ParseMultipartForm(50 << 20) // 50MB limit for multiple images
		if err != nil {
			http.Error(writer, "Failed to parse multipart form", http.StatusBadRequest)
			return
		}

		// Extract form fields
		userId := req.FormValue("userId")
		title := req.FormValue("title")
		description := req.FormValue("description")

		// Handle multiple image uploads
		var images []model.Image
		files := req.MultipartForm.File["images"]
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				continue // Skip files that can't be opened
			}
			defer file.Close()

			// Read image data
			imageData, err := io.ReadAll(file)
			if err != nil {
				continue // Skip files that can't be read
			}

			// Validate image type
			mimeType := fileHeader.Header.Get("Content-Type")
			if mimeType != "image/jpeg" && mimeType != "image/png" && mimeType != "image/gif" && mimeType != "image/webp" {
				continue // Skip invalid image formats
			}

			images = append(images, model.Image{
				Data:     imageData,
				MimeType: mimeType,
				Filename: fileHeader.Filename,
			})
		}

		blog = model.BeforeCreateTour(userId, title, description, images)
	}

	err = handler.BlogService.Create(blog)
	if err != nil {
		println("Error while creating a new blog")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
	println("Blog successfully created")
}

func (handler *BlogHandler) Delete(writer http.ResponseWriter, req *http.Request) {
	idStr := mux.Vars(req)["id"]
	log.Printf("Deleting blog with ID: %s", idStr)

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(writer, "Invalid UUID", http.StatusBadRequest)
		return
	}

	err = handler.BlogService.Delete(id)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]string{"message": "Blog deleted successfully"})
}

func (handler *BlogHandler) Update(writer http.ResponseWriter, req *http.Request) {
	idStr := mux.Vars(req)["id"]
	log.Printf("Updating blog with ID: %s", idStr)

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(writer, "Invalid UUID", http.StatusBadRequest)
		return
	}

	var input struct {
		Title       string        `json:"title"`
		Description string        `json:"description"`
		Images      []model.Image `json:"images"`
	}

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	updatedBlog := model.Blog{
		Title:       input.Title,
		Description: input.Description,
		Images:      input.Images,
	}

	err = handler.BlogService.Update(id, updatedBlog)
	if err != nil {
		http.Error(writer, "Blog not found", http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]string{"message": "Blog updated successfully"})
}

func (handler *BlogHandler) GetAllByUser(writer http.ResponseWriter, req *http.Request) {
	userID := mux.Vars(req)["userId"]

	writer.Header().Set("Content-Type", "application/json")

	blogs, err := handler.BlogService.GetAllByUser(userID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(blogs)
}
