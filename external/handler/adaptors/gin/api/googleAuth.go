package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"project-management-service/external/auth/adaptors/google"
	"project-management-service/internal/core/entities"
	"project-management-service/internal/core/logic"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	userSrv logic.UserService
}

func NewAuthHandler(userSrv logic.UserService) *AuthHandler {
	return &AuthHandler{userSrv: userSrv}
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	googleConfig := google.Config()
	url := googleConfig.AuthCodeURL("randomstate")
	http.Redirect(c.Writer, c.Request, url, http.StatusSeeOther)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	state := c.DefaultQuery("state", "")
	if state != "randomstate" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "State mismatch"})
		return
	}

	code := c.DefaultQuery("code", "")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code is missing"})
		return
	}

	googleConfig := google.Config()
	tokeng, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token exchange failed"})
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tokeng.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user data"})
		return
	}
	defer resp.Body.Close()

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user data"})
		return
	}

	var userInfo struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.Unmarshal(userData, &userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal user data"})
		return
	}
	userId := uuid.New().String()

	h.userSrv.CreateUser(entities.User{
		ID:      userId,
		Email:   userInfo.Email,
		Name:    userInfo.Name,
		Picture: userInfo.Picture,
	})

	token, err := h.userSrv.GenerateJWT(userId, userInfo.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}
	c.SetCookie("Authorization", token, 3600*24, "/", "localhost", false, true)
	c.Data(http.StatusOK, "application/json", userData)
}
