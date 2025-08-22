package handler

import (
	"encoding/json"

	"net/http"

	"blog.xws.com/model"
	"blog.xws.com/service"
	"github.com/gorilla/mux"
)

type LikeHandler struct {
	LikeService *service.LikeService
}

func (handler *LikeHandler) Create(writer http.ResponseWriter, req *http.Request) {
	var like model.Like
	err := json.NewDecoder(req.Body).Decode(&like)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.LikeService.CreateLike(&like)
	if err != nil {
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(like)
}

func (handler *LikeHandler) Delete(writer http.ResponseWriter, req *http.Request) {
	userId := mux.Vars(req)["userId"]
	blogId := mux.Vars(req)["blogId"]

	/*userIdd, err := uuid.Parse(userId)
	blogIdd, err := uuid.Parse(blogId)

	if err != nil {
		http.Error(writer, "Invalid like user ID", http.StatusBadRequest)
		print(err)
		return
	}*/
	err := handler.LikeService.Delete(userId, blogId)
	if err != nil {
		println("Error while deleting a like")
		println(err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusAccepted)
}

func (handler *LikeHandler) GetByBlogId(writer http.ResponseWriter, req *http.Request) {
	blogId := mux.Vars(req)["blogId"]

	/*id, err := uuid.Parse(blogId)
	if err != nil {
		http.Error(writer, "Invalid blog status", http.StatusBadRequest)
		return
	}*/

	comments, err := handler.LikeService.GetByBlogId(blogId)
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
