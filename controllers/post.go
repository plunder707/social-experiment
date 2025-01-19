// controllers/post.go
package controllers

import (
    "context"
    "log"
    "net/http"
    "strings"
    "time"

    "social-experiment/models"
    "social-experiment/utils"
    "social-experiment/websocket"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

// CreatePost handles creating a new post
func CreatePost(db *mongo.Collection, hub *websocket.Hub) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Retrieve userID from context
        userID, exists := c.Get("userID")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
            return
        }

        // Bind JSON input to request struct
        var req struct {
            Content string `json:"content" binding:"required"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
            log.Printf("[WARNING] Invalid post request: %v", err)
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
            return
        }

        // Input validation and sanitization
        req.Content = strings.TrimSpace(req.Content)
        if req.Content == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Post content cannot be empty"})
            return
        }

        safeContent := utils.SanitizeInput(req.Content)

        // Convert userID from string to primitive.ObjectID
        objectID, err := primitive.ObjectIDFromHex(userID.(string))
        if err != nil {
            log.Printf("[ERROR] Invalid user ID format: %v", err)
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
            return
        }

        // Retrieve user from database
        var user models.User
        err = db.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
        if err != nil {
            if err == mongo.ErrNoDocuments {
                c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            } else {
                log.Printf("[ERROR] Error fetching user: %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing request"})
            }
            return
        }

        // Create a new Post instance with a new ObjectID
        post := models.Post{
            ID:        primitive.NewObjectID(),
            UserID:    user.ID,
            Username:  user.Username,
            Content:   safeContent,
            CreatedAt: time.Now(),
        }

        // Insert the post into the database
        _, err = db.InsertOne(context.Background(), post)
        if err != nil {
            log.Printf("[ERROR] Error creating post: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating post"})
            return
        }

        // Broadcast the new post to WebSocket clients
        hub.BroadcastPost(post)

        // Respond with the created post
        c.JSON(http.StatusOK, post)
    }
}


// GetPosts handles retrieving all posts
func GetPosts(db *mongo.Collection) gin.HandlerFunc {
    return func(c *gin.Context) {
        cursor, err := db.Find(context.Background(), bson.M{}, options.Find().SetSort(bson.D{{"created_at", -1}}))
        if err != nil {
            log.Printf("[ERROR] Error fetching posts: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching posts"})
            return
        }
        defer cursor.Close(context.Background())

        var posts []models.Post
        if err = cursor.All(context.Background(), &posts); err != nil {
            log.Printf("[ERROR] Error decoding posts: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding posts"})
            return
        }

        c.JSON(http.StatusOK, posts)
    }
}
