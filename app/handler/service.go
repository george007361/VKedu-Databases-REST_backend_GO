package handler

import (
	"net/http"

	"github.com/george007361/db-course-proj/app/helpers"
	"github.com/gin-gonic/gin"
)

func (h *Handler) serviceClear(c *gin.Context) {
	err := h.services.Managment.Clear()
	if err.Code != http.StatusOK {
		helpers.NewErrorResponse(c, err.Code, err.Message)
		return
	}

	c.JSON(err.Code, err.Message)
}

func (h *Handler) serviceGetStatus(c *gin.Context) {
	data, err := h.services.Managment.GetStatus()
	if err.Code != http.StatusOK {
		helpers.NewErrorResponse(c, err.Code, err.Message)
		return
	}

	c.JSON(err.Code, data)
}
