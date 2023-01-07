package handler

import (
	"net/http"
	"strconv"

	"github.com/george007361/db-course-proj/app/helpers"
	"github.com/george007361/db-course-proj/app/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) forumCreate(c *gin.Context) {
	logrus.Println("Handle create forum")

	var newForumData models.Forum

	if err := c.BindJSON(&newForumData); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	logrus.Println(newForumData)

	forumData, err := h.services.Forum.CreateForum(newForumData)
	if err.Code != http.StatusCreated {
		helpers.NewErrorResponse(c, err.Code, err.Message)
		return
	}
	logrus.Println(forumData)

	c.JSON(http.StatusCreated, forumData)
}

func (h *Handler) forumGetData(c *gin.Context) {
	logrus.Println("Handle get forum")

	slug, isExist := c.Params.Get("slug")
	if !isExist {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "No slug in URL")
		return
	}

	forumData, err := h.services.Forum.GetForumData(slug)
	if err.Code != http.StatusOK {
		helpers.NewErrorResponse(c, err.Code, err.Message)
		return
	}

	c.JSON(err.Code, forumData)
}

func (h *Handler) forumGetUsers(c *gin.Context) {

	logrus.Println("Handle get forum users")

	// limit=100&since=george&desc=false

	slug, isExist := c.Params.Get("slug")
	if !isExist {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "No slug in URL")
		return
	}

	limitStr, isExist := c.GetQuery("limit")
	if !isExist {
		limitStr = "100"
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "Invalig query param limit")
		return
	}

	since, isExist := c.GetQuery("since")
	if !isExist {
		since = ""
	}

	descStr, isExist := c.GetQuery("desc")
	if !isExist {
		descStr = "false"
	}

	var desc bool
	switch descStr {
	case "false":
		desc = false
		break
	case "true":
		desc = true
		break
	default:
		helpers.NewErrorResponse(c, http.StatusBadRequest, "Invalig query param desc")
		return
	}

	queryParams := models.ForumUsersQueryParams{
		Slug:  slug,
		Limit: limit,
		Desc:  desc,
		Since: since,
	}

	users, errr := h.services.Forum.GetForumUsers(queryParams)
	if errr.Code != http.StatusOK {
		helpers.NewErrorResponse(c, errr.Code, errr.Message)
		return
	}

	c.JSON(errr.Code, users)
}

func (h *Handler) forumGetThreads(c *gin.Context) {
	logrus.Println("Handle get forum threads")

	slug, isExist := c.Params.Get("slug")
	if !isExist {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "No slug in URL")
		return
	}

	limitStr, isExist := c.GetQuery("limit")
	if !isExist {
		limitStr = "100"
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "Invalig query param limit")
		return
	}

	sinceStr, isExist := c.GetQuery("since")
	if !isExist {
		sinceStr = "-1"
	}

	since, err := strconv.Atoi(sinceStr)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "Invalig query param since")
		return
	}

	descStr, isExist := c.GetQuery("desc")
	if !isExist {
		descStr = "false"
	}

	var desc bool
	switch descStr {
	case "false":
		desc = false
		break
	case "true":
		desc = true
		break
	default:
		helpers.NewErrorResponse(c, http.StatusBadRequest, "Invalig query param desc")
		return
	}

	queryParams := models.ForumThreadsQueryParams{
		Slug:  slug,
		Limit: limit,
		Desc:  desc,
		Since: since,
	}

	users, errr := h.services.Forum.GetForumThreads(queryParams)
	if errr.Code != http.StatusOK {
		helpers.NewErrorResponse(c, errr.Code, errr.Message)
		return
	}

	c.JSON(errr.Code, users)
}

func (h *Handler) forumCreateThread(c *gin.Context) {
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

	logrus.Println(newThreadData)

	threadData, err := h.services.Forum.CreateThreadInForum(newThreadData)
	if err.Code != http.StatusCreated {
		helpers.NewErrorResponse(c, err.Code, err.Message)
		return
	}
	logrus.Println(threadData)

	c.JSON(http.StatusCreated, threadData)
}
