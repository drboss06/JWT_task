package handler

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getToken(c *gin.Context) {

	var guIdQuery = c.Query("guid")

	if guIdQuery == "" {
		newErrorResponse(c, http.StatusBadRequest, "guid is empty")
		return
	}

	clientIp := c.ClientIP()
	if clientIp == "" {
		newErrorResponse(c, http.StatusBadRequest, "client ip is empty")
		return

	}

	token, refreshToken, err := h.services.Authorization.GenerateToken(guIdQuery, clientIp)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token":         token,
		"refresh_token": refreshToken,
	})
}

type refreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *Handler) refresh(c *gin.Context) {
	guIdQuery := c.Query("guid")

	var input refreshInput

	clientIp := c.ClientIP()
	if clientIp == "" {
		newErrorResponse(c, http.StatusBadRequest, "client ip is empty")
		return
	}

	if guIdQuery == "" {
		newErrorResponse(c, http.StatusBadRequest, "guid is empty")
		return
	}

	c.BindJSON(&input)

	if input.RefreshToken == "" {
		newErrorResponse(c, http.StatusBadRequest, "refresh token is empty")
		return
	}

	refreshTokenBase, err := base64.StdEncoding.DecodeString(input.RefreshToken)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	accessToken, refreshToken, err := h.services.RefreshToken(refreshTokenBase, guIdQuery, clientIp)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token":         accessToken,
		"refresh_token": refreshToken,
	})
}
