// controllers/auth.go
package controllers

import (
	"context"
	"log"
	"strings"
	"time"

	"maliaki-backend/models"
	"maliaki-backend/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Register handles user registration
func Register(db *mongo.Collection, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("[WARNING] Invalid registration request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Input validation
		req.Username = strings.TrimSpace(req.Username)
		if req.Username == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
			return
		}

		// Check if user exists
		count, err := db.CountDocuments(context.Background(), bson.M{"username": req.Username})
		if err != nil {
			log.Printf("[ERROR] Database error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}

		// Hash password
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			log.Printf("[ERROR] Error hashing password: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing password"})
			return
		}

		// Create user
		user := models.User{
			Username:  req.Username,
			Password:  hashedPassword,
			CreatedAt: time.Now().Format(time.RFC3339),
		}

		result, err := db.InsertOne(context.Background(), user)
		if err != nil {
			log.Printf("[ERROR] Error creating user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
			return
		}

		// Generate JWT
		userID := result.InsertedID.(primitive.ObjectID).Hex()
		token, err := utils.GenerateJWT(userID, jwtSecret)
		if err != nil {
			log.Printf("[ERROR] Error generating token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

// Login handles user authentication
func Login(db *mongo.Collection, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("[WARNING] Invalid login request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Input validation
		req.Username = strings.TrimSpace(req.Username)
		if req.Username == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
			return
		}

		// Find user
		var user models.User
		err := db.FindOne(context.Background(), bson.M{"username": req.Username}).Decode(&user)
		if err != nil {
			log.Printf("[WARNING] User not found: %s", req.Username)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Check password
		if !utils.CheckPasswordHash(req.Password, user.Password) {
			log.Printf("[WARNING] Invalid password for user: %s", req.Username)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Generate JWT
		token, err := utils.GenerateJWT(user.ID.Hex(), jwtSecret)
		if err != nil {
			log.Printf("[ERROR] Error generating token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
