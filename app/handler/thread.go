package handler

import (
	"net/http"
	"strconv"

	"github.com/george007361/db-course-proj/app/helpers"
	"github.com/george007361/db-course-proj/app/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) threadCreatePosts(c *gin.Context) {
	slugOrId, isExists := c.Params.Get("slug_or_id")
	if !isExists {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "No slug or id")
		return
	}

	newPostsData := make([]models.Post, 0)
	if err := c.BindJSON(&newPostsData); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	id, err := strconv.Atoi(slugOrId)

	var postsData []models.Post
	var errr models.Error

	if err != nil {
		// slug
		postsData, errr = h.services.Thread.CreatePostsByThreadSlug(newPostsData, slugOrId)
	} else {
		// id
		postsData, errr = h.services.Thread.CreatePostsByThreadId(newPostsData, id)
	}

	if errr.Code != http.StatusCreated {
		helpers.NewErrorResponse(c, errr.Code, errr.Message)
		return
	}

	c.JSON(errr.Code, postsData)
}

func (h *Handler) threadGetData(c *gin.Context) {

}

func (h *Handler) threadUpdateData(c *gin.Context) {

}

func (h *Handler) threadGetPosts(c *gin.Context) {

}

func (h *Handler) threadVote(c *gin.Context) {

}
