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

	var newUserData models.UserCreate
	newUserData.Nickname = nickname

	if err := c.BindJSON(&newUserData); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	users, err := h.services.User.CreateUser(newUserData)
	logrus.Println("HANDLER RES: ", users, err)
	if err.Code == http.StatusConflict {
		c.JSON(err.Code, users)
		return
	}

	if err.Code == http.StatusCreated {
		c.JSON(err.Code, newUserData)
		return
	}

	helpers.NewErrorResponse(c, err.Code, err.Message)
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

	var updatedData models.UserUpdate
	updatedData.Nickname = nickname

	if err := c.BindJSON(&updatedData); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	userData, err := h.services.User.UpdateUserProfile(updatedData)

	if err.Code != http.StatusOK {
		helpers.NewErrorResponse(c, err.Code, err.Message)
		return
	}

	c.JSON(err.Code, userData)
}
