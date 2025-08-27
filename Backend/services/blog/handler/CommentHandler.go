package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"blog.xws.com/model"
	"blog.xws.com/service"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CommentHandler struct {
	CommentService *service.CommentService
}

func (handler *CommentHandler) Create(writer http.ResponseWriter, req *http.Request) {
	var comment model.Comment
	err := json.NewDecoder(req.Body).Decode(&comment)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.CommentService.CreateComment(&comment)
	if err != nil {
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(comment)
}

func (handler *CommentHandler) Update(writer http.ResponseWriter, req *http.Request) {
	var comment model.Comment
	err := json.NewDecoder(req.Body).Decode(&comment)
	if err != nil {
		fmt.Println("Error while parsing json:", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = handler.CommentService.UpdateComment(&comment)
	if err != nil {
		fmt.Println("Error while updating the blog:", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(comment)
}
func (handler *CommentHandler) GetById(writer http.ResponseWriter, req *http.Request) {
	idStr := mux.Vars(req)["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(writer, "Invalid UUID", http.StatusBadRequest)
		return
	}

	blog, err := handler.CommentService.GetById(id) // make sure your service takes uuid.UUID
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(blog)
}
func (handler *CommentHandler) Delete(writer http.ResponseWriter, req *http.Request) {
	idV := mux.Vars(req)["id"]

	id, err := uuid.Parse(idV)

	if err != nil {
		http.Error(writer, "Invalid comment ID", http.StatusBadRequest)
		return
	}
	err = handler.CommentService.Delete(id)
	if err != nil {
		println("Error while deleting a facility")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusAccepted)
}

func (handler *CommentHandler) GetByBlogId(writer http.ResponseWriter, req *http.Request) {
	idV := mux.Vars(req)["id"]

	id, err := uuid.Parse(idV)
	if err != nil {
		http.Error(writer, "Invalid blog status", http.StatusBadRequest)
		return
	}

	comments, err := handler.CommentService.GetByBlogId(id)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(writer).Encode(comments); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
