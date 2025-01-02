package google

// import (
// 	"context"
// 	"encoding/json"
// 	"io"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt/v5"
// )

// func GoogleCallback(c *gin.Context) {
// 	state := c.DefaultQuery("state", "")
// 	if state != "randomstate" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "State mismatch"})
// 		return
// 	}

// 	code := c.DefaultQuery("code", "")
// 	if code == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Code is missing"})
// 		return
// 	}

// 	googleConfig := config.SetupConfig()
// 	tokeng, err := googleConfig.Exchange(context.Background(), code)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token exchange failed"})
// 		return
// 	}

// 	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tokeng.AccessToken)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user data"})
// 		return
// 	}
// 	defer resp.Body.Close()

// 	userData, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user data"})
// 		return
// 	}

// 	// Parse user data
// 	var userInfo struct {
// 		ID    string `json:"id"`
// 		Email string `json:"email"`
// 		Name  string `json:"name"`
// 	}
// 	if err := json.Unmarshal(userData, &userInfo); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal user data"})
// 		return
// 	}

// 	user, err := ac.UserRepository.FindOrCreateUser(userInfo.Email)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database operation failed"})
// 		return
// 	}

// 	// Generate JWT
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"sub": user.Email,
// 		"exp": time.Now().Add(24 * time.Hour).Unix(),
// 	})

// 	secret := os.Getenv("SECRET")
// 	tokenString, err := token.SignedString([]byte(secret))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
// 		return
// 	}

// 	c.SetCookie("Authorization", tokenString, 3600*24, "/", "localhost", false, true)
// 	// Send the user data as a JSON response
// 	c.Data(http.StatusOK, "application/json", userData)
// }
