package handler

import (
	"net/http"
	"strconv"

	"github.com/george007361/db-course-proj/app/helpers"
	"github.com/george007361/db-course-proj/app/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

	logrus.Error(errr)
	if errr.Code != http.StatusCreated {
		helpers.NewErrorResponse(c, errr.Code, errr.Message)
		return
	}

	c.JSON(errr.Code, postsData)
}

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

func (h *Handler) threadGetPosts(c *gin.Context) {
	// Slug or id
	slugOrId, isExists := c.Params.Get("slug_or_id")
	if !isExists {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "No slug or id")
		return
	}

	var err error

	// Limit
	limitStr, isExist := c.GetQuery("limit")
	limit := 100
	if isExist {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			helpers.NewErrorResponse(c, http.StatusBadRequest, "Invalig query param limit")
			return
		}
	}

	// Since
	sinceStr, isExist := c.GetQuery("since")
	since := 0
	if isExist {
		since, err = strconv.Atoi(sinceStr)
		if err != nil {
			helpers.NewErrorResponse(c, http.StatusBadRequest, "Invalig query param since")
			return
		}
	}

	// Desc
	descStr, isExist := c.GetQuery("desc")
	desc := false
	if isExist {
		switch descStr {
		case "true":
			desc = true
		case "false":
			desc = false
		default:
			helpers.NewErrorResponse(c, http.StatusBadRequest, "Invalig query param desc")
			return
		}
	}

	// Sort
	sortStr, isExist := c.GetQuery("sort")
	sort := "flat"
	logrus.Println(sortStr, isExist)
	if isExist {
		sort = sortStr
		// switch sortStr {
		// case "flat":
		// case "tree":
		// case "parent_tree":
		// 	sort = sortStr
		// default:
		// 	helpers.NewErrorResponse(c, http.StatusBadRequest, "Invalig query param desc")
		// 	return
		// }
	}

	queryParams := models.ThreadGetPostsParams{
		Limit: limit,
		Desc:  desc,
		Since: since,
		Sort:  sort,
	}

	var errr models.Error
	var threadPosts []models.Post

	id, err := strconv.Atoi(slugOrId)

	if err != nil {
		// slug
		threadPosts, errr = h.services.Thread.GetThreadPostsBySlug(queryParams, slugOrId)
	} else {
		// id
		threadPosts, errr = h.services.Thread.GetThreadPostsById(queryParams, id)
	}

	if errr.Code != http.StatusOK {
		helpers.NewErrorResponse(c, errr.Code, errr.Message)
		return
	}

	logrus.Println(len(threadPosts))

	// TODO FIX THIS SHIT
	if len(threadPosts) == 0 {
		c.JSON(errr.Code, []string{})
		return
	}

	c.JSON(errr.Code, threadPosts)
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
	newThreadData.Forum = slug

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
