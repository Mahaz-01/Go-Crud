package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"gin-crud/ent"
	"gin-crud/ent/user"
	"gin-crud/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {
	// Define a struct to bind the incoming request body to
	var user struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Bind the incoming JSON body into the user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Ensure the required fields are provided (binding:"required" should handle this, but keeping for consistency)
	if user.Username == "" || user.Email == "" || user.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Username, email, and password are required"})
		return
	}

	// Hash the user's password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		log.Printf("Error hashing password: %v", err)
		return
	}

	// Create the new user in the database
	ctx := context.Background()
	createdUser, err := models.Client.User.
		Create().
		SetUsername(user.Username).
		SetEmail(user.Email).
		SetPassword(string(hashedPassword)).
		Save(ctx)

	if err != nil {
		if ent.IsConstraintError(err) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		log.Printf("Error registering user: %v", err)
		return
	}

	// Generate a JWT token for the newly registered user
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
		return
	}

	// Create a JWT token with user details and expiration time
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      createdUser.ID,
		"username": createdUser.Username,
		"email":    createdUser.Email,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign the token with the secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		log.Printf("Error generating token: %v", err)
		return
	}

	// Return the response with a success message and the JWT token
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   tokenString,
	})
}

func LoginUser(c *gin.Context) {
	var credentials struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Bind the credentials from the request body
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Query the user by username
	ctx := context.Background()
	dbUser, err := models.Client.User.
		Query().
		Where(user.Username(credentials.Username)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		log.Printf("Error fetching user: %v", err)
		return
	}

	// Check if the provided password matches the stored password hash
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(credentials.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate a JWT token if the login is successful
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
		return
	}

	// Create a new JWT token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      dbUser.ID,
		"username": dbUser.Username,
		"email":    dbUser.Email,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign the JWT token with the secret key
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		log.Printf("Error generating token: %v", err)
		return
	}

	// Return the generated token in the response
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
