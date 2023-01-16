package handler

import (
	"net/http"
	"strconv"

	"github.com/george007361/db-course-proj/app/helpers"
	"github.com/george007361/db-course-proj/app/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) threadGetData(c *gin.Context) {
	slugOrId, isExists := c.Params.Get("slug_or_id")
	if !isExists {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "No slug or id")
		return
	}

	id, err := strconv.Atoi(slugOrId)
	var errr models.Error
	var threadData models.Thread

	if err != nil {
		// slug
		threadData, errr = h.services.Thread.GetThreadData(slugOrId)
	} else {
		// id
		threadData, errr = h.services.Thread.GetThreadDataById(id)
	}

	if errr.Code != http.StatusOK {
		helpers.NewErrorResponse(c, errr.Code, errr.Message)
		return
	}

	c.JSON(errr.Code, threadData)
}

func (h *Handler) threadUpdateData(c *gin.Context) {
	slugOrId, isExists := c.Params.Get("slug_or_id")
	if !isExists {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "No slug or id")
		return
	}

	var newThreadData models.UpdateThread
	if err := c.BindJSON(&newThreadData); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var errr models.Error
	var threadData models.Thread

	id, err := strconv.Atoi(slugOrId)

	if err != nil {
		// slug
		threadData, errr = h.services.Thread.UpdateThreadBySlug(newThreadData, slugOrId)
	} else {
		// id
		threadData, errr = h.services.Thread.UpdateThreadById(newThreadData, id)
	}

	if errr.Code != http.StatusOK {
		helpers.NewErrorResponse(c, errr.Code, errr.Message)
		return
	}

	c.JSON(errr.Code, threadData)
}

func (h *Handler) threadVote(c *gin.Context) {
	slugOrId, isExists := c.Params.Get("slug_or_id")
	if !isExists {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "No slug or id")
		return
	}

	var vote models.Vote
	if err := c.BindJSON(&vote); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var errr models.Error
	var threadData models.Thread

	id, err := strconv.Atoi(slugOrId)

	if err != nil {
		// slug
		threadData, errr = h.services.Thread.VoteThreadBySlug(vote, slugOrId)
	} else {
		// id
		threadData, errr = h.services.Thread.VoteThreadById(vote, id)
	}

	if errr.Code != http.StatusOK {
		helpers.NewErrorResponse(c, errr.Code, errr.Message)
		return
	}

	c.JSON(errr.Code, threadData)
}

func (h *Handler) createThread(c *gin.Context) {
	logrus.Println("Handle forum create thread")

	slug, isExist := c.Params.Get("slug")
	if !isExist {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "No slug in URL")
		return
	}

	var newThreadData models.Thread
	newThreadData.ForumSlug = slug

	if err := c.BindJSON(&newThreadData); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	threadData, err := h.services.Thread.CreateThread(newThreadData)
	if err.Code == http.StatusCreated || err.Code == http.StatusConflict {
		c.JSON(err.Code, threadData)
		return
	}

	helpers.NewErrorResponse(c, err.Code, err.Message)
}
