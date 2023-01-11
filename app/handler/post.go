package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/george007361/db-course-proj/app/helpers"
	"github.com/george007361/db-course-proj/app/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) postGetData(c *gin.Context) {
	logrus.Println("Handle get post data")

	idStr, _ := c.Params.Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "Invalid post id") //"No nickname in URL"
	}
	var postAllData models.PostAllData

	postData, errr := h.services.Post.GetPostData(id)
	if errr.Code != http.StatusOK {
		helpers.NewErrorResponse(c, errr.Code, errr.Message)
		return
	}
	postAllData.Post = postData

	related := strings.Split(c.Query("related"), ",")

	for _, item := range related {
		switch item {
		case "user":
			userData, err := h.services.User.GetUserProfile(postData.Author)
			if err.Code != http.StatusOK {
				helpers.NewErrorResponse(c, err.Code, err.Message)
				return
			}
			postAllData.Author = &userData
		case "forum":
			forumData, err := h.services.Forum.GetForumData(postData.Forum)
			if err.Code != http.StatusOK {
				helpers.NewErrorResponse(c, err.Code, err.Message)
				return
			}
			postAllData.Forum = &forumData
		case "thread":
			threadData, err := h.services.Thread.GetThreadDataById(postData.Thread)
			if err.Code != http.StatusOK {
				helpers.NewErrorResponse(c, err.Code, err.Message)
				return
			}
			postAllData.Thread = &threadData
		}
	}

	c.JSON(http.StatusOK, postAllData)
}

func (h *Handler) postUpdate(c *gin.Context) {
	logrus.Println("Handle post update data")

	idStr, _ := c.Params.Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "Invalid post id") //"No nickname in URL"
	}

	var newData models.PostUpdate
	if err := c.BindJSON(&newData); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedPostData, errr := h.services.Post.UpdatePostData(newData, id)
	if errr.Code != http.StatusOK {
		helpers.NewErrorResponse(c, errr.Code, errr.Message)
		return
	}

	c.JSON(http.StatusOK, updatedPostData)
}
