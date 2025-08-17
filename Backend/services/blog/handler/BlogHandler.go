package handler

import (
	"encoding/json"
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
	var blog model.Blog
	err := json.NewDecoder(req.Body).Decode(&blog)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.BlogService.Create(&blog)
	if err != nil {
		println("Error while creating a new blog")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
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
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Images      []string `json:"images"`
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
