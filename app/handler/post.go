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
			userData, err := h.services.User.GetUserProfile(postData.AuthorNickname)
			if err.Code != http.StatusOK {
				helpers.NewErrorResponse(c, err.Code, err.Message)
				return
			}
			postAllData.Author = &userData
		case "forum":
			forumData, err := h.services.Forum.GetForumData(postData.ForumSlug)
			if err.Code != http.StatusOK {
				helpers.NewErrorResponse(c, err.Code, err.Message)
				return
			}
			postAllData.Forum = &forumData
		case "thread":
			threadData, err := h.services.Thread.GetThreadDataById(postData.ThreadId)
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

func (h *Handler) getPostsInThread(c *gin.Context) {
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
	var threadPosts []*models.Post

	id, err := strconv.Atoi(slugOrId)

	if err != nil {
		// slug
		threadPosts, errr = h.services.Post.GetPostsByThreadSlug(queryParams, slugOrId)
	} else {
		// id
		threadPosts, errr = h.services.Post.GetPostsByThreadId(queryParams, id)
	}

	if errr.Code != http.StatusOK {
		helpers.NewErrorResponse(c, errr.Code, errr.Message)
		return
	}

	// TODO FIX THIS SHIT
	if len(threadPosts) == 0 {
		c.JSON(errr.Code, []string{})
		return
	}

	c.JSON(errr.Code, threadPosts)
}

func (h *Handler) createPostsInThread(c *gin.Context) {
	slugOrId, isExists := c.Params.Get("slug_or_id")
	if !isExists {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "No slug or id")
		return
	}

	newPostsData := make([]*models.Post, 0)
	if err := c.BindJSON(&newPostsData); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	id, err := strconv.Atoi(slugOrId)

	// var postsData []*models.Post
	var errr models.Error

	if err != nil {
		// slug
		newPostsData, errr = h.services.Post.CreatePostsByThreadSlug(newPostsData, slugOrId)
	} else {
		// id
		newPostsData, errr = h.services.Post.CreatePostsByThreadId(newPostsData, id)
	}

	if errr.Code != http.StatusCreated {
		helpers.NewErrorResponse(c, errr.Code, errr.Message)
		return
	}

	c.JSON(errr.Code, newPostsData)
}
