package handler

import (
	"authentication-service/external/auth/adaptors/google"
	"authentication-service/internal/core/entities"
	"authentication-service/internal/core/logic"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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

	redirectTo := c.DefaultQuery("redirect_to", "http://localhost:3000/project")
	state := "randomstate|" + redirectTo

	url := googleConfig.AuthCodeURL(state)
	http.Redirect(c.Writer, c.Request, url, http.StatusSeeOther)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	// Parse state to extract redirect URL
	state := c.DefaultQuery("state", "")
	parts := strings.SplitN(state, "|", 2)
	if len(parts) != 2 || parts[0] != "randomstate" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "State mismatch"})
		return
	}
	redirectTo := parts[1] // Extract the redirect URL

	code := c.DefaultQuery("code", "")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code is missing"})
		return
	}

	googleConfig := google.Config()
	tokeng, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Print(err.Error())
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
	err = h.userSrv.CreateUser(entities.User{
		ID:      userId,
		Email:   userInfo.Email,
		Name:    userInfo.Name,
		Picture: userInfo.Picture,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	token, err := h.userSrv.GenerateJWT(userId, userInfo.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "Authorization",
		Value:    token,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect to the original URL or a default one
	c.Redirect(http.StatusFound, redirectTo)
}
