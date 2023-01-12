package handler

import (
	"github.com/george007361/db-course-proj/app/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		forum := api.Group("/forum")
		{
			forum.POST("/create", h.forumCreate)
			forum.GET("/:slug/details", h.forumGetData)
			forum.POST("/:slug/create", h.createThread)
			forum.GET("/:slug/users", h.forumGetUsers)
			forum.GET("/:slug/threads", h.forumGetThreads)
		}

		post := api.Group("/post")
		{
			post.GET("/:id/details", h.postGetData)
			post.POST("/:id/details", h.postUpdate)
		}

		service := api.Group("/service")
		{
			service.POST("/clear", h.serviceClear)
			service.GET("/status", h.serviceGetStatus)
		}

		thread := api.Group("/thread/:slug_or_id")
		{
			thread.POST("/create", h.threadCreatePosts)
			thread.GET("/details", h.threadGetData)
			thread.POST("/details", h.threadUpdateData)
			thread.GET("/posts", h.threadGetPosts)
			thread.POST("/vote", h.threadVote)
		}

		user := api.Group("/user/:nickname")
		{
			user.POST("/create", h.userCreate)
			user.GET("/profile", h.userGetProfile)
			user.POST("/profile", h.userUpdateProfile)
		}
	}

	return router
}
