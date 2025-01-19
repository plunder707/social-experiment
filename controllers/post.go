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
    "go.mongodb.org/mongo-driver/mongo/options"
)

// CreatePost handles creating a new post
func CreatePost(db *mongo.Collection, hub *websocket.Hub) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("userID")
        if !exists {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
            return
        }

        var req struct {
            Content string `json:"content"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
            log.Printf("[WARNING] Invalid post request: %v", err)
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
            return
        }

        req.Content = strings.TrimSpace(req.Content)
        if req.Content == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Post content cannot be empty"})
            return
        }

        safeContent := utils.SanitizeInput(req.Content)

        objectID, err := primitive.ObjectIDFromHex(userID.(string))
        if err != nil {
            log.Printf("[ERROR] Invalid user ID format: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing request"})
            return
        }

        var user models.User
        err = db.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
        if err != nil {
            log.Printf("[ERROR] Error fetching user: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing request"})
            return
        }

        post := models.Post{
            UserID:    userID.(string),  // Use userID as string
            Username:  user.Username,
            Content:   safeContent,
            CreatedAt: time.Now().Format(time.RFC3339),
        }

        result, err := db.InsertOne(context.Background(), post)
        if err != nil {
            log.Printf("[ERROR] Error creating post: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating post"})
            return
        }

        insertedID := result.InsertedID.(primitive.ObjectID).Hex()  // Convert ObjectID to string
        post.ID = insertedID  // Set the post ID as a string

        hub.BroadcastPost(post)

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
