package main

import (
    "context"
    "log"

    "social-experiment/controllers"
    "social-experiment/middleware"
    "social-experiment/utils"
    "social-experiment/websocket"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    // Load configuration
    config := utils.LoadConfig()

    // Initialize MongoDB
    clientOptions := options.Client().ApplyURI(config.MongoURI)
    mongoClient, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatalf("[ERROR] Failed to connect to MongoDB: %v", err)
    }
    defer func() {
        if err := mongoClient.Disconnect(context.Background()); err != nil {
            log.Printf("[ERROR] Failed to disconnect MongoDB: %v", err)
        }
    }()

    userCollection := mongoClient.Database("social-experiment").Collection("users")
    postCollection := mongoClient.Database("social-experiment").Collection("posts")

    // Initialize WebSocket Hub
    hub := websocket.NewHub()
    go hub.Run()

    // Initialize Gin Router
    router := gin.Default()

    // Apply Security Headers Middleware
    if config.SecurityHeaders {
        router.Use(middleware.SecurityHeadersMiddleware())
    }

    // Apply CORS Middleware
    router.Use(func(c *gin.Context) {
        origin := c.GetHeader("Origin")
        if isAllowedOrigin(origin, config.CORSOrigins) {
            c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
            c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
            c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
        }
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    // Initialize Rate Limiter
    rl := middleware.NewRateLimiter(config.RateLimit, config.RateBurst)
    router.Use(middleware.RateLimitMiddleware(rl))

    // Define Routes
    router.POST("/register", controllers.Register(userCollection, config.JWTSecret))
    router.POST("/login", controllers.Login(userCollection, config.JWTSecret))
    router.POST("/posts", middleware.AuthMiddleware(config.JWTSecret), controllers.CreatePost(postCollection, hub))
    router.GET("/posts", middleware.AuthMiddleware(config.JWTSecret), controllers.GetPosts(postCollection))
    router.GET("/ws", func(c *gin.Context) {
        hub.HandleWebSocket(c.Writer, c.Request)
    })

    // Start Server
    address := ":" + config.ServerPort
    log.Printf("[INFO] Starting server on %s", address)
    if err := router.Run(address); err != nil {
        log.Fatalf("[ERROR] Failed to run server: %v", err)
    }
}

// isAllowedOrigin checks if the origin is allowed based on the CORS configuration
func isAllowedOrigin(origin string, allowedOrigins []string) bool {
    if len(allowedOrigins) == 0 {
        return false
    }
    for _, ao := range allowedOrigins {
        if ao == "*" || ao == origin {
            return true
        }
    }
    return false
}
