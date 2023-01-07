package handler

import (
	"net/http"

	"github.com/george007361/db-course-proj/app/helpers"
	"github.com/george007361/db-course-proj/app/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) userCreate(c *gin.Context) {
	logrus.Println("Handle create user")

	nickname, isExist := c.Params.Get("nickname")
	if !isExist {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "No nickname in URL") //"No nickname in URL"
		return
	}

	var newUserData models.User
	newUserData.Nickname = nickname

	if err := c.BindJSON(&newUserData); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	err := h.services.User.CreateUser(newUserData)
	if err.Code != http.StatusCreated {
		helpers.NewErrorResponse(c, err.Code, err.Message)
		return
	}

	c.JSON(http.StatusCreated, newUserData)
}

func (h *Handler) userGetProfile(c *gin.Context) {
	logrus.Println("Handle get user profile")

	nickname, isExist := c.Params.Get("nickname")
	if !isExist {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "No nickname in URL") //"No nickname in URL"
		return
	}

	userData, err := h.services.User.GetUserProfile(nickname)
	if err.Code != http.StatusOK {
		helpers.NewErrorResponse(c, err.Code, err.Message)
		return
	}

	c.JSON(err.Code, userData)

}

func (h *Handler) userUpdateProfile(c *gin.Context) {

	logrus.Println("Handle update user profile")

	nickname, isExist := c.Params.Get("nickname")
	if !isExist {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "No nickname in URL") //"No nickname in URL"
		return
	}

	var userData models.User
	userData.Nickname = nickname

	if err := c.BindJSON(&userData); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	userUpdatedData, err := h.services.User.UpdateUserProfile(userData)

	if err.Code != http.StatusOK {
		helpers.NewErrorResponse(c, err.Code, err.Message)
		return
	}

	c.JSON(err.Code, userUpdatedData)
}
